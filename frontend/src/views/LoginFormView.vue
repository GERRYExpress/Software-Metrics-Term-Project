<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import BaseButton from '@/components/BaseButton.vue';
import BaseInput from '@/components/BaseInput.vue';
import BaseToggle from '@/components/BaseToggle.vue';
import type { LoginForm } from '@/types';
import { UserRound } from '@lucide/vue';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
const router = useRouter()
const loginOptions = ['Buyer', 'Seller']
const credential = ref<LoginForm>({
    username: '',
    password: '',
    isSeller: false
})

const toggleRole = () => {
    credential.value.isSeller = !credential.value.isSeller
    router.push({ query: {
        role: String(credential.value.isSeller)
    }})
}

const handleLogin = () => {
    router.push({
        name: 'landing-page-view'
    })
}

</script>


<template>
    <main class="flex flex-col items-center justify-center my-16 gap-y-5">
        <h1>Welcome Back!</h1>
        <AppCard class="w-100">
            <RouterLink class=" place-self-start" :to="{ path: '' }">
                Back to Home
            </RouterLink>
            <UserRound class="size-16"/>
            <form @submit.prevent="handleLogin" class="flex flex-col w-full gap-y-1.5">
                <BaseInput label="Username:" v-model="credential.username"/>
                <BaseInput label="Password:" v-model="credential.password"/>
                <span class="w-full text-center">
                    No account? <RouterLink :to="{name: 'register-view'}">Sign up here</RouterLink>
                </span>
                <div class="flex flex-row items-center justify-around w-full">
                    <p class="font-bold">Login as</p>
                    <BaseToggle @toggleRole="toggleRole" v-model="credential.isSeller" :options="loginOptions"/>
                </div>
                <BaseButton text="Login">Sign In</BaseButton>
            </form>
        </AppCard>
    </main>
</template>