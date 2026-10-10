<script setup lang="ts">
import GameCard from '@/components/GameCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import BaseInput from '@/components/BaseInput.vue';
import GameService from '@/service/GameService';
import type { Game } from '@/types';
import { computed, onMounted, ref, watchEffect } from 'vue';
import { useRoute } from 'vue-router';
const route = useRoute()
const games = ref<Game[] | null>(null)
const keyword = computed(() => {
    return typeof route.query.keyword === 'string' ? route.query.keyword : ''
})
const props = defineProps<{
    page: number
}>()
onMounted(() => {
    games.value = null
    watchEffect(() => {
        updateKeyword()
    })
})

const page = computed(() => props.page)
const totalPages = ref<number>(0)
const hasNextPage = computed(() => {
    return page.value < Math.ceil(totalPages.value / 12)
})
const updateKeyword = () => {
    let queryFunction
    if (keyword.value === '') {
        queryFunction = GameService.getGames(12, page.value)
    } else {
        queryFunction = GameService.getGamesByKeyword(keyword.value)
    }
    queryFunction.then(res => {
        games.value = res.data
        totalPages.value = res.headers['x-total-count']
    })
}
</script>


<template>
    <ProductMainContainer class="flex flex-col w-4xl place-self-center">
        <h1 class="w-full">Search Results: {{ $route.query.keyword }}</h1>
        <div class="grid grid-cols-4 gap-5 w-auto">
            <GameCard v-for="game in games" :key="game.id" :game/>
        </div>
        <div class="w-full text-xl flex flex-row items-center gap-x-3">
            <RouterLink v-if="page > 1" :to="{ query: { ...route.query, page: page - 1}}">Prev</RouterLink>
            <BaseInput disabled class="w-16! m-0! text-center" v-model="page"/>
            <RouterLink v-if="hasNextPage" :to="{ query: { ...route.query,  page: page + 1 }}">Next</RouterLink>
        </div>
    </ProductMainContainer>
</template>