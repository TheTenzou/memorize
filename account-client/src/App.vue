<template>
  <h1 class="text-4xl font-bol text-center my-12">Scaffolded App Works Well!</h1>
  <h3 class="text-sl text-center" v-if="errorCode">
    Error code: {{ errorCode }}
  </h3>
  <h3 class="text-center" v-if="errorMessage">{{ errorMessage }}</h3>
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue'

export default defineComponent({
  name: 'App',
  setup() {
    const errorCode = ref(null)
    const errorMessage = ref(null)

    onMounted(async () => {
      const response = await fetch('/api/account/me', {
        method: 'GET',
      })

      const body = await response.json()
      
      errorCode.value = response.status
      errorMessage.value = body.error.message
    })
    return {
      errorCode,
      errorMessage,
    }
  },
})
</script>
