<script setup lang="ts">
import GameCard from '@/components/GameCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import GameService from '@/service/GameService';
import type { Game } from '@/types';
import { onMounted, ref, watchEffect } from 'vue';

const games = ref<Game[] | null>()

onMounted(() => {
    watchEffect(() => {
        GameService.getGames()
            .then(res => {
                games.value = res.data
            })
    })
})

</script>

<template>
    <ProductMainContainer class="flex flex-col w-4xl place-self-center">
        <h1 class="w-full">HIGHLIGHT</h1>
        <img src="/Landing/Landing-Page-Carousel-5.png" class="bg-tertiary w-full h-full rounded-xl"/>
        <h1 class="w-full">GAMES</h1>
        <div v-if="games" class="grid grid-cols-4 gap-5 w-auto">
            <GameCard v-for="game in games.slice(0, 8)" :game :key="game.id"/>
        </div>
        <RouterLink :to="{ name: 'search-result-view' }" class="place-self-center">EXPLORE MORE</RouterLink>
    </ProductMainContainer>
</template>