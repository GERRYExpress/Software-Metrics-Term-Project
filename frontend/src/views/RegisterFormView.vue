<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import AuthMainContainer from '@/components/AuthMainContainer.vue';
import BaseButton from '@/components/BaseButton.vue';
import BaseInput from '@/components/BaseInput.vue';
import type { RegisterForm } from '@/types';
import { UserRoundPlus } from '@lucide/vue';
import { computed, ref } from 'vue';
const credential = ref<RegisterForm>({
    username: '',
    password: '',
    confirmPassword: ''
})

const isPasswordValid = computed(() => {
    return credential.value.password.trim() === credential.value.confirmPassword.trim()
})

const handleRegister = () => {
    if (!isPasswordValid.value) {
        return
    }
    /* Create account query */
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
            <form @submit.prevent="handleRegister" class="flex flex-col w-full gap-y-1.5">
                <BaseInput label="Username" v-model="credential.username"/>
                <BaseInput label="Username" v-model="credential.password"/>
                <BaseInput label="Username" v-model="credential.confirmPassword"/>
                <span class="w-full text-center">
                    Already have an account? <RouterLink :to="{ name: 'login-view' }">Sign in here</RouterLink>
                </span>
                <BaseButton>Register</BaseButton>
            </form>
        </AppCard>
    </AuthMainContainer>
</template>