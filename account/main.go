package main

import (
	"context"
	"log"
	"memorize/inject"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting server....")

	dataSources, err := inject.InitDataSources()
	if err != nil {
		log.Fatalf("Unable to initialze data sources: %v\n", err)
	}

	repositories := inject.InitRepositories(dataSources)

	services, err := inject.InitServices(repositories)
	if err != nil {
		log.Fatalf("Unable to initilaze services: %v\n", err)
	}

	router := inject.InitRouter(services)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server started.")

	gracefullyShutdown(dataSources, server)
}

// Shutdown server and close connection
func gracefullyShutdown(dataSources *inject.DataSources, server *http.Server) {
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initilize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", server.Addr)

	// wait for kill signal of channel
	quit := make(chan os.Signal, 10)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This bocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the requiest it is currntly hadling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := dataSources.Close(); err != nil {
		log.Fatalf("A problem occured gracefully shutting down data sources: %v\n", err)
	}

	// Shutdown server
	log.Println("Shutting down server....")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
