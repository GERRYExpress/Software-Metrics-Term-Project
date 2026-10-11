<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import AuthMainContainer from '@/components/AuthMainContainer.vue';
import BaseButton from '@/components/BaseButton.vue';
import BaseInput from '@/components/BaseInput.vue';
import BaseToggle from '@/components/BaseToggle.vue';
import AccountService from '@/service/AccountService';
import { useAuthStore } from '@/stores/auth';
import type { LoginForm } from '@/types';
import { UserRound } from '@lucide/vue';
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const loginOptions = ['Buyer', 'Seller']
const credential = ref<LoginForm>({
    username: '',
    password: '',
    isSeller: route.query.role === 'true'
})
const errorMessage = ref('')
const isLoading = ref(false)

const toggleRole = () => {
    credential.value.isSeller = !credential.value.isSeller
    router.push({ query: {
        role: String(credential.value.isSeller)
    }})
}

const handleLogin = async () => {
    errorMessage.value = ''
    isLoading.value = true
    try {
        const response = await AccountService.authenticate(credential.value)
        authStore.setSession(response.data.account, response.data.token)
        await router.push({
            name: credential.value.isSeller ? 'seller-dashboard-view' : 'landing-page-view'
        })
    } catch (error) {
        errorMessage.value = getErrorMessage(error)
    } finally {
        isLoading.value = false
    }
}

function getErrorMessage(error: unknown) {
    const responseMessage = (error as { response?: { data?: { error?: string } } })?.response?.data?.error
    return responseMessage || 'Unable to sign in. Check that the backend is running and try again.'
}
</script>


<template>
    <AuthMainContainer>
        <h1>Welcome Back!</h1>
        <AppCard class="w-100">
            <RouterLink class=" place-self-start" :to="{ name: 'landing-page-view' }">
                Back to Home
            </RouterLink>
            <UserRound class="size-16"/>
            <p v-if="errorMessage" class="w-full rounded-lg bg-red/20 p-3 text-center text-red-100" role="alert">
                {{ errorMessage }}
            </p>
            <form @submit.prevent="handleLogin" class="flex flex-col w-full gap-y-1.5">
                <BaseInput label="Username" v-model="credential.username" placeholder="Username"/>
                <BaseInput label="Password" type="password" v-model="credential.password" placeholder="Password"/>
                <span class="w-full text-center">
                    No account? <RouterLink :to="{name: 'register-view'}">Sign up here</RouterLink>
                </span>
                <div class="flex flex-row items-center justify-around w-full">
                    <p class="font-bold">Login as</p>
                    <BaseToggle @toggleRole="toggleRole" v-model="credential.isSeller" :options="loginOptions"/>
                </div>
                <BaseButton :disabled="isLoading">{{ isLoading ? 'Signing in...' : 'Login' }}</BaseButton>
            </form>
        </AppCard>
    </AuthMainContainer>
</template>