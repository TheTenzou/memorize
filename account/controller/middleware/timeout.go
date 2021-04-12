package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"memorize/models/apperrors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(timeout time.Duration, errTimeout *apperrors.Error) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// wrap ginContext.Writer with timeoutWriter
		timeoutWr := &timeoutWriter{
			ResponseWriter: ginContext.Writer,
			header:         make(http.Header),
		}

		// update gin writer
		ginContext.Writer = timeoutWr

		// wrap the request context with a timeout
		timeoutContext, cancel := context.WithTimeout(ginContext.Request.Context(), timeout)
		defer cancel()

		// update gin request context
		ginContext.Request = ginContext.Request.WithContext(timeoutContext)

		finished := make(chan struct{})        // to indicate handler finished
		panicChan := make(chan interface{}, 1) // used to handle panics if we can't recover

		go func() {
			// recovering panic
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			ginContext.Next() // calls subsequent middleware(s) and handler
			finished <- struct{}{}
		}()

		select {
		case <-panicChan:
			handlePanic(timeoutWr)
		case <-finished:
			handleFinished(timeoutWr)
		case <-timeoutContext.Done():
			handleTimeout(timeoutWr, errTimeout, ginContext)
		}
	}
}

func handlePanic(timeoutWr *timeoutWriter) {
	// if we cannot recover from panic,
	// send internal server error
	err := apperrors.NewInternal()
	timeoutWr.ResponseWriter.WriteHeader(err.Status())
	eResp, _ := json.Marshal(gin.H{
		"error": err,
	})

	timeoutWr.ResponseWriter.Write(eResp)
}

func handleFinished(timeoutWr *timeoutWriter) {
	// if finished, set headers and write resp
	timeoutWr.mutex.Lock()
	defer timeoutWr.mutex.Unlock()

	// map Headers from tw.Header() (written to by gin)
	// to timeoutWr.ResponseWriter for response
	header := timeoutWr.ResponseWriter.Header()
	for i, value := range timeoutWr.Header() {
		header[i] = value
	}

	timeoutWr.ResponseWriter.WriteHeader(timeoutWr.code)
	// timeoutWr.wbuf will have been written to already when gin writes to timeoutWr.Write()
	timeoutWr.ResponseWriter.Write(timeoutWr.wbuf.Bytes())
}

func handleTimeout(timeoutWr *timeoutWriter, errTimeout *apperrors.Error, ginContext *gin.Context) {
	// timeout has occurred, send errTimeout and write headers
	timeoutWr.mutex.Lock()
	defer timeoutWr.mutex.Unlock()
	// ResponseWriter from gin
	timeoutWr.ResponseWriter.Header().Set("Content-Type", "application/json")
	timeoutWr.ResponseWriter.WriteHeader(errTimeout.Status())

	errorResponse, _ := json.Marshal(gin.H{
		"error": errTimeout,
	})

	timeoutWr.ResponseWriter.Write(errorResponse)
	ginContext.Abort()
	timeoutWr.SetTimedOut()
}

// implements http.Writer, but tracks if Writer has timed out
// or has already written its header to prevent
// header and body overwrites
// also locks access to this writer to prevent race conditions
// holds the gin.ResponseWriter which we'll manually call Write()
// on in the middleware function to send response
type timeoutWriter struct {
	gin.ResponseWriter

	header      http.Header
	wbuf        bytes.Buffer // The zero value for Buffer is an empty buffer ready to use.
	mutex       sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

// Writes the response, but first makes sure there hasn't already been a timeout
func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mutex.Lock()
	defer tw.mutex.Unlock()

	if tw.timedOut {
		return 0, nil
	}

	return tw.wbuf.Write(b)
}

// In http.ResponseWriter interface
func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	tw.mutex.Lock()
	defer tw.mutex.Unlock()

	// We do not write the header if we've timed out or written the header
	if tw.timedOut || tw.wroteHeader {
		return
	}

	tw.wroteHeader = true
	tw.code = code
}

// Header "relays" the header, h, set in struct
// In http.ResponseWriter interface
func (tw *timeoutWriter) Header() http.Header {
	return tw.header
}

// SetTimeOut sets timedOut field to true
func (tw *timeoutWriter) SetTimedOut() {
	tw.timedOut = true
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
