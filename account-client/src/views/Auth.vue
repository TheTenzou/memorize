<template>
  <div class="max-w-xl mx-auto px-4">
    <div class="rounded-lg shadow-lg p-4">
      <ul class="mx-8 mb-2 flex justify-center">
        <li
          @click="setIsLogin(true)"
          class="mx-2 coursor-pointer text-center hover:opacity-75 trasition-opacity"
          :class="{ 'border-b-2 border-blue-400': isLogin }"
        >
          Login
        </li>
        <li
          @click="setIsLogin(false)"
          class="mx-2 cursor-pointer text-center hover:opacity-75 translition-opacity"
          :class="{ 'border-b-2 border-blue-400': !isLogin }"
        >
          Sign Up
        </li>
      </ul>
      <LoginForm
        :isLogin="isLogin"
        :isSubmitting="isLoading"
        @submitAuth="authSubmitted"
      />
      <div v-if="error" class="text-center my-2">
        <p class="text-red-400">{{ error.message }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref } from 'vue'
import LoginForm from '../components/LoginForm.vue'
import { useAuth } from '../store/auth'

export default defineComponent({
  name: 'Auth',

  components: {
    LoginForm,
  },

  setup() {
    const isLogin = ref(true)
    const { currentUser, error, isLoading, signin, signup } = useAuth({
      onAuthRoute: '/',
    })

    const setIsLogin = (nextVal) => {
      isLogin.value = nextVal
    }

    const authSubmitted = ({ login, password }) => {
      isLogin.value ? signin(login, password) : signup(login, password)
    }

    return {
      isLogin,
      setIsLogin,
      authSubmitted,
      currentUser,
      error,
      isLoading,
    }
  },
})
</script>
