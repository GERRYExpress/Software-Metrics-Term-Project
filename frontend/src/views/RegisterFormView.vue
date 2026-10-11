<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import AuthMainContainer from '@/components/AuthMainContainer.vue';
import BaseButton from '@/components/BaseButton.vue';
import BaseInput from '@/components/BaseInput.vue';
import AccountService from '@/service/AccountService';
import type { RegisterForm } from '@/types';
import { UserRoundPlus } from '@lucide/vue';
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
const router = useRouter()
const credential = ref<RegisterForm>({
    username: '',
    password: '',
    confirmPassword: ''
})
const errorMessage = ref('')
const isLoading = ref(false)

const isPasswordValid = computed(() => {
    return credential.value.password.trim() === credential.value.confirmPassword.trim()
})

const handleRegister = async () => {
    if (!isPasswordValid.value) {
        errorMessage.value = 'Passwords do not match.'
        return
    }
    errorMessage.value = ''
    isLoading.value = true
    try {
        await AccountService.register(credential.value)
        await router.push({ name: 'login-view' })
    } catch (error) {
        const responseMessage = (error as { response?: { data?: { error?: string } } })?.response?.data?.error
        errorMessage.value = responseMessage || 'Unable to create the account. Please try again.'
    } finally {
        isLoading.value = false
    }
}

</script>


<template>
    <AuthMainContainer>
        <h1>Create new account</h1>
        <AppCard class="w-100">
            <RouterLink class=" place-self-start" :to="{ name: 'landing-page-view' }">
                Back to Home
            </RouterLink>
            <UserRoundPlus class="size-16"/>
            <p v-if="errorMessage" class="w-full rounded-lg bg-red/20 p-3 text-center text-red-100" role="alert">
                {{ errorMessage }}
            </p>
            <form @submit.prevent="handleRegister" class="flex flex-col w-full gap-y-1.5">
                <BaseInput label="Username" v-model="credential.username" placeholder="Username"/>
                <BaseInput label="Password" type="password" v-model="credential.password" placeholder="Password"/>
                <BaseInput label="Confirm Password" type="password" v-model="credential.confirmPassword" placeholder="Confirm Password"/>
                <span class="w-full text-center">
                    Already have an account? <RouterLink :to="{ name: 'login-view' }">Sign in here</RouterLink>
                </span>
                <BaseButton :disabled="isLoading">{{ isLoading ? 'Creating account...' : 'Register' }}</BaseButton>
            </form>
        </AppCard>
    </AuthMainContainer>
</template>