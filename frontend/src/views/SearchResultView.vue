<script setup lang="ts">
import GameCard from '@/components/GameCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import GameService from '@/service/GameService';
import type { Game } from '@/types';
import { onMounted, ref, watchEffect } from 'vue';
import { useRoute } from 'vue-router';
const route = useRoute()
const games = ref<Game[] | null>(null)
defineProps<{
    page: number
}>()
onMounted(() => {
    watchEffect(() => {
        /* If not searched or filter */
        if (route.query.keyword === '' && route.query.category === '') {
            GameService.getGames()
                .then(res => {
                    games.value = res.data
                })
        }
    })
})
</script>


<template>
    <ProductMainContainer class="flex flex-col">
        <h1 class="w-full">Search Results</h1>
        <div class="grid grid-cols-5 gap-5 w-auto">
            <GameCard v-for="game in games" :key="game.id" :game/>
        </div>
        <div class="bg-black">
            <RouterLink v-if="page >= 1" :to="{ query: { page: page - 1 }}"/>
            <RouterLink :to="{ query: { page: page + 1 }}"/>
        </div>
    </ProductMainContainer>
</template>