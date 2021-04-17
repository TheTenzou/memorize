import { reactive, provide, inject, toRefs, readonly } from 'vue'
import { storeTokens, doRequest, getTokenPayload } from '../util'

const state = reactive({
  currentUser: null,
  idToken: null,
  isLoading: false,
  error: null,
})

const signin = async (login, password) =>
  await authenticate(login, password, '/api/account/signin')

const signup = async (login, password) =>
  await authenticate(login, password, '/api/account/signup')

export const authStrore = {
  ...toRefs(readonly(state)),
  signin,
  signup,
}

const storeSymbol = Symbol()

export function provideAuth() {
  provide(storeSymbol, authStrore)
}

export function useAuth() {
  const store = inject(storeSymbol)

  if (!store) {
    throw new Error('Auth store has not been instantiated!')
  }

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

  console.log(tokens)
  console.log(tokens.accessToken)

  const tokenClaims = getTokenPayload(tokens.accessToken)

  state.accessToken = tokens.accessToken
  state.currentUser = tokenClaims.user
  state.isLoading = false
}
