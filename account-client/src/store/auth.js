import { reactive, provide, inject, toRefs, readonly } from 'vue'

const state = reactive({
  currentUser: null,
  idToken: null,
  isLoading: false,
  error: null,
})

export const authStrore = {
  ...toRefs(readonly(state)),
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
