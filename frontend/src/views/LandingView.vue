<script setup lang="ts">
import GameCard from '@/components/GameCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import WelcomeModal from '@/components/WelcomeModal.vue';
import GameService from '@/service/GameService';
import type { Game } from '@/types';
import { ChevronLeft, ChevronRight } from '@lucide/vue';
import { onMounted, onUnmounted, ref, watchEffect } from 'vue';

const games = ref<Game[] | null>(null)
const showWelcome = ref(false)

const banners = [
    '/Landing/Landing-Page-Carousel-1.png',
    '/Landing/Landing-Page-Carousel-2.png',
    '/Landing/Landing-Page-Carousel-3.png',
    '/Landing/Landing-Page-Carousel-4.png',
    '/Landing/Landing-Page-Carousel-5.png',
]

const currentBannerIndex = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const nextBanner = () => {
    currentBannerIndex.value = (currentBannerIndex.value + 1) % banners.length
}

const prevBanner = () => {
    currentBannerIndex.value = (currentBannerIndex.value - 1 + banners.length) % banners.length
}

onMounted(() => {
    watchEffect(() => {
        GameService.getGames(8, 1)
            .then(res => {
                games.value = res.data
            })
            .catch(err => {
                console.error('Failed to load games:', err)
            })
    })

    // Auto-advance banner every 5 seconds
    timer = setInterval(nextBanner, 5000)

    // Check if welcome popup was already displayed
    if (!sessionStorage.getItem('kd5f_welcome_seen')) {
        showWelcome.value = true
    }
})

onUnmounted(() => {
    if (timer) clearInterval(timer)
})

const closeWelcome = () => {
    showWelcome.value = false
    sessionStorage.setItem('kd5f_welcome_seen', 'true')
}
</script>

<template>
    <ProductMainContainer class="flex flex-col w-4xl place-self-center gap-y-6">
        <!-- Highlights Carousel -->
        <div class="flex flex-col w-full gap-y-2">
            <h1 class="w-full text-2xl">HIGHLIGHTS</h1>
            <div class="relative w-full h-100 bg-tertiary rounded-2xl overflow-hidden shadow-2xl group">
                <img
                    :src="banners[currentBannerIndex]"
                    alt="Highlight banner"
                    class="w-full h-full object-cover transition-all duration-500 ease-in-out"
                />

                <!-- Prev/Next Controls -->
                <button
                    @click="prevBanner"
                    type="button"
                    class="absolute left-3 top-1/2 -translate-y-1/2 bg-black/60 hover:bg-black/90 text-white p-2.5 rounded-full transition-all cursor-pointer opacity-80 hover:opacity-100 hover:scale-110">
                    <ChevronLeft class="size-6"/>
                </button>
                <button
                    @click="nextBanner"
                    type="button"
                    class="absolute right-3 top-1/2 -translate-y-1/2 bg-black/60 hover:bg-black/90 text-white p-2.5 rounded-full transition-all cursor-pointer opacity-80 hover:opacity-100 hover:scale-110">
                    <ChevronRight class="size-6"/>
                </button>

                <!-- Indicator Dots -->
                <div class="absolute bottom-4 left-1/2 -translate-x-1/2 flex flex-row gap-x-2 bg-black/50 px-3 py-1.5 rounded-full">
                    <button
                        v-for="(_, idx) in banners"
                        :key="idx"
                        @click="currentBannerIndex = idx"
                        type="button"
                        class="size-2.5 rounded-full transition-all cursor-pointer"
                        :class="idx === currentBannerIndex ? 'bg-action scale-125' : 'bg-white/50 hover:bg-white'"
                    />
                </div>
            </div>
        </div>

        <!-- Featured Games Grid -->
        <div class="flex flex-col w-full gap-y-3">
            <h1 class="w-full text-2xl">FEATURED GAMES</h1>
            <div v-if="games" class="grid grid-cols-4 gap-5 w-full">
                <GameCard v-for="game in games.slice(0, 8)" :game="game" :key="game.id"/>
            </div>
            <div v-else class="text-center py-10 text-white/50">
                Loading games...
            </div>
        </div>

        <div class="flex flex-col items-center gap-y-2 my-4 text-center">
            <RouterLink :to="{ name: 'search-result-view' }" class="text-xl">
                EXPLORE MORE
            </RouterLink>
            <p class="text-white/60 text-sm">
                Or get started:
                <RouterLink :to="{ name: 'login-view' }" class="text-action">Sign In</RouterLink> or
                <RouterLink :to="{ name: 'register-view' }" class="text-action">Register</RouterLink>
            </p>
        </div>

        <!-- Welcome Greeting Modal -->
        <WelcomeModal :show="showWelcome" @close="closeWelcome"/>
    </ProductMainContainer>
</template>