<template>
  <h1 class="text-3xl text-center">User Account</h1>
  <loader
    v-if="meLoading"
    :height="256"
    class="animate-spin stroke-current text-blue-500 mx-auto"
  />
  <UserForm
    v-if="user"
    :user="user"
    @detailsSubmitted="handleDetailsSubmitted"
  />
  <button
    type="button"
    class="btn btn-red w-32 block mx-auto my-2"
    @click="signout"
  >
    signout
  </button>
  <p v-if="meError" class="text-center text-red-500">Error fetching user</p>
  <p v-if="updateDetailsError" class="text-center text-red-500">
    Failed to update user details
  </p>
</template>

<script>
import { computed, defineComponent } from 'vue'
import { useAuth } from '../store/auth'
import { useRequest } from '../util/index'
import UserForm from '../components/UserForm.vue'
import Loader from '../components/ui/Loader.vue'

export default defineComponent({
  name: 'Details',
  components: {
    UserForm,
    Loader,
  },

  setup() {
    const { accessToken, signout } = useAuth()

    const { data: meData, error: meError, loading: meLoading } = useRequest(
      {
        url: '/api/account/me',
        method: 'get',
        headers: {
          Authorization: `Bearer ${accessToken.value}`,
        },
      },
      {
        execOnMounted: true,
      }
    )

    const {
      exec: updateDetails,
      data: updateDetailsData,
      error: updateDetailsError,
      loading: updateDetailsLoading,
    } = useRequest({
      url: '/api/account/details',
      method: 'put',
      headers: {
        Authorization: `Bearer ${accessToken.value}`,
      },
    })

    const handleDetailsSubmitted = (userDetails) => {
      updateDetails(userDetails)
    }

    const loading = computed(() => {
      return meLoading || updateDetailsLoading.value
    })

    const user = computed(() => {
      return updateDetailsData?.value?.user || meData?.value?.user
    })

    return {
      loading,
      meData,
      meError,
      user,
      handleDetailsSubmitted,
      updateDetailsError,
      signout,
    }
  },
})
</script>
