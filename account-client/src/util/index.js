import axios from 'axios'
import jwt_decode from 'jwt-decode'

import { ref, onMounted } from 'vue'

export const doRequest = async (requestOptions) => {
  let error
  let data

  try {
    const response = await axios.request(requestOptions)
    data = response.data
  } catch (err) {
    if (err.response) {
      error = err.response.data.error
    } else if (err.request) {
      error = err.request
    } else {
      error = err
    }
  }

  return {
    data,
    error,
  }
}

export const useRequest = (requestOptions, options) => {
  const { execOnMounted } = options || {}
  const error = ref(null)
  const data = ref(null)
  const loading = ref(false)

  const exec = async (requestData) => {
    data.value = null
    loading.value = true
    error.value = null

    if (requestData) {
      requestOptions = {
        ...requestOptions,
        data: requestData,
      }
    }

    const response = await doRequest(requestOptions)

    data.value = response.data
    error.value = response.error
    loading.value = false
  }

  onMounted(() => {
    if (execOnMounted) {
      exec()
    }
  })

  return {
    exec,
    error,
    data,
    loading,
  }
}

const accessTokenKey = '__memorizerAccess'
const refreshTokenKey = '__memorizerRefresh'

// store access and refresh tokens
export const storeTokens = (accessToken, refreshToken) => {
  localStorage.setItem(accessTokenKey, accessToken)
  localStorage.setItem(refreshTokenKey, refreshToken)
}

export const getTokens = () => {
  return [
    localStorage.getItem(accessTokenKey),
    localStorage.getItem(refreshTokenKey),
  ]
}

export const getTokenPayload = (token) => {
  if (!token) {
    return null
  }

  const tokenClaims = jwt_decode(token)

  if (Date.now() / 1000 >= tokenClaims.exp) {
    return null
  }

  return tokenClaims
}
