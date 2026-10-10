<script setup lang="ts">
import GameCard from '@/components/GameCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import GameService from '@/service/GameService';
import type { Game } from '@/types';
import { computed, onMounted, ref, watchEffect } from 'vue';
import { useRoute } from 'vue-router';
const route = useRoute()
const games = ref<Game[] | null>(null)
const keyword = computed(() => {
    return typeof route.query.keyword === 'string' ? route.query.keyword : ''
})
defineProps<{
    page: number
}>()
onMounted(() => {
    games.value = null
    watchEffect(() => {
        updateKeyword()
    })
})

const updateKeyword = () => {
    let queryFunction
    if (keyword.value === '') {
        queryFunction = GameService.getGames()
    } else {
        queryFunction = GameService.getGamesByKeyword(keyword.value)
    }
    queryFunction.then(res => {
        games.value = res.data
    })
}
</script>


<template>
    <ProductMainContainer class="flex flex-col w-4xl place-self-center">
        <h1 class="w-full">Search Results: {{ $route.query.keyword }}</h1>
        <div class="grid grid-cols-4 gap-5 w-auto">
            <GameCard v-for="game in games" :key="game.id" :game/>
        </div>
        <div class="bg-black">
            <RouterLink v-if="page >= 1" :to="{ query: { page: page - 1 }}"/>
            <RouterLink :to="{ query: { page: page + 1 }}"/>
        </div>
    </ProductMainContainer>
</template>