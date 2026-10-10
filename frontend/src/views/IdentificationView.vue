<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import AuthMainContainer from '@/components/AuthMainContainer.vue';
import { ShoppingBag, Store, ArrowLeft } from '@lucide/vue';
import { useRouter } from 'vue-router';

const router = useRouter()

const chooseRole = (isSeller: boolean) => {
    router.push({
        name: 'login-view',
        query: {
            role: String(isSeller)
        }
    })
}
</script>

<template>
    <AuthMainContainer>
        <h1>Sign In</h1>
        <AppCard class="w-100 flex flex-col items-center gap-y-6 p-8">
            <RouterLink class="place-self-start flex items-center gap-x-1 text-sm font-semibold" :to="{ name: 'landing-page-view' }">
                <ArrowLeft class="size-4"/> Back to Home
            </RouterLink>

            <h2 class="text-2xl font-bold text-center">Why are you here?</h2>

            <div class="flex flex-col gap-y-4 w-full">
                <button
                    @click="chooseRole(false)"
                    type="button"
                    class="bg-action hover:bg-lime-400 text-tertiary font-bold px-6 py-4 rounded-2xl shadow-lg transition-all hover:scale-102 flex items-center justify-center gap-x-3 text-lg cursor-pointer">
                    <ShoppingBag class="size-6"/>
                    I'm here to buy a game!
                </button>

                <button
                    @click="chooseRole(true)"
                    type="button"
                    class="bg-secondary border-2 border-white/30 hover:border-action text-white hover:text-action font-bold px-6 py-4 rounded-2xl shadow-lg transition-all hover:scale-102 flex items-center justify-center gap-x-3 text-lg cursor-pointer">
                    <Store class="size-6"/>
                    I'm here to sell a game!
                </button>
            </div>
        </AppCard>
    </AuthMainContainer>
</template>
