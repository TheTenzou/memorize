import { reactive, provide, inject, toRefs, readonly, watchEffect } from 'vue'
import { storeTokens, doRequest, getTokenPayload, getTokens, removeTokens } from '../util'
import { useRouter } from 'vue-router'

const state = reactive({
  currentUser: null,
  accessToken: null,
  isLoading: false,
  error: null,
})

const storeSymbol = Symbol()

const signin = async (login, password) =>
  await authenticate(login, password, '/api/account/signin')

const signup = async (login, password) =>
  await authenticate(login, password, '/api/account/signup')

const signout = async () => {
  state.isLoading = true
  state.error = null

  const { error } = await doRequest({
    url: '/api/account/signout',
    method: 'post',
    headers: {
      Authorization: `Bearer ${state.accessToken}`,
    },
  })

  if (error) {
    state.error = error
    state.isLoading = false
    return
  }

  state.currentUser = null
  state.accessToken = null

  removeTokens()

  state.isLoading = false
}                          

const initializeUser = async () => {
  state.isLoading = true
  state.error = null

  const [accessToken, refreshToken] = getTokens()

  const accessTokenClaims = getTokenPayload(accessToken)
  const refreshTokenClaims = getTokenPayload(refreshToken)

  if (accessTokenClaims) {
    state.accessToken = accessToken
    state.currentUser = accessTokenClaims.user
  }

  state.isLoading = false

  if (!refreshTokenClaims) {
    return
  }

  const { data, error } = await doRequest({
    url: '/api/account/tokens',
    method: 'post',
    data: {
      refreshToken,
    },
  })

  if (error) {
    return
  }

  const { tokens } = data
  storeTokens(tokens.accessToken, tokens.refreshToken)

  const updatedTokenClaims = getTokenPayload(tokens.accessToken)

  state.currentUser = updatedTokenClaims.user
  state.accessToken = tokens.accessToken
}

export const createAuthStore = (authStoreOption) => {
  const { onAuthRoute, requireAuthRoute } = authStoreOption || {}

  const authStore = {
    ...toRefs(readonly(state)),
    signin,
    signup,
    signout,
    initializeUser,
    onAuthRoute,
    requireAuthRoute,
  }

  return {
    authStore,
    install: (app) => {
      app.provide(storeSymbol, authStore)
    },
  }
}

export const authStore = {
  ...toRefs(readonly(state)),
  signin,
  signup,
}

export function useAuth() {
  const store = inject(storeSymbol)

  if (!store) {
    throw new Error('Auth store has not been instantiated!')
  }

  const router = useRouter()

  watchEffect(() => {
    if (store.currentUser.value && store.onAuthRoute) {
      router.push(store.onAuthRoute)
    }

    if (!store.currentUser.value && store.requireAuthRoute) {
      router.push(store.requireAuthRoute)
    }
  })

  return store
}

const authenticate = async (login, password, url) => {
  state.isLoading = true
  state.error = null

  const { data, error } = await doRequest({
    url,
    method: 'post',
    data: {
      login: login,
      password,
    },
  })

  if (error) {
    state.error = error
    state.isLoading = false
    return
  }

  const { tokens } = data

  storeTokens(tokens.accessToken, tokens.refreshToken)

  const tokenClaims = getTokenPayload(tokens.accessToken)

  state.accessToken = tokens.accessToken
  state.currentUser = tokenClaims.user
  state.isLoading = false
}
