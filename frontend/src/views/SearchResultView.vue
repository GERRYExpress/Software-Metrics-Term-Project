<script setup lang="ts">
import AppPagination from '@/components/AppPagination.vue';
import GameCard from '@/components/GameCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import GameService from '@/service/GameService';
import type { Game } from '@/types';
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute()
const router = useRouter()
const allGames = ref<Game[]>([])
const isLoading = ref(false)
const perPage = 8

const keyword = computed(() => {
    return typeof route.query.keyword === 'string' ? route.query.keyword : ''
})

const currentPage = computed(() => {
    const p = Number(route.query.page)
    return isNaN(p) || p < 1 ? 1 : p
})

const totalPages = computed(() => {
    return Math.max(1, Math.ceil(allGames.value.length / perPage))
})

const paginatedGames = computed(() => {
    const start = (currentPage.value - 1) * perPage
    return allGames.value.slice(start, start + perPage)
})

const fetchGames = () => {
    isLoading.value = true
    const query = keyword.value ? GameService.getGamesByKeyword(keyword.value) : GameService.getGames()
    query
        .then(res => {
            allGames.value = res.data || []
        })
        .catch(err => {
            console.error('Failed to fetch search games:', err)
            allGames.value = []
        })
        .finally(() => {
            isLoading.value = false
        })
}

onMounted(() => {
    fetchGames()
})

watch(() => route.query.keyword, () => {
    fetchGames()
})

const handlePageChange = (newPage: number) => {
    router.push({
        query: {
            ...route.query,
            page: newPage
        }
    })
}
</script>

<template>
    <ProductMainContainer class="flex flex-col w-5xl place-self-center">
        <h1 class="w-full">
            {{ keyword ? `Search Results: "${keyword}"` : 'All Games' }}
        </h1>

        <div v-if="paginatedGames.length > 0" class="grid grid-cols-4 gap-5 w-full my-2">
            <GameCard v-for="game in paginatedGames" :key="game.id" :game="game"/>
        </div>

        <div v-else-if="!isLoading" class="flex flex-col items-center justify-center py-16 bg-tertiary rounded-xl w-full text-center">
            <h2 class="text-xl">No games found</h2>
            <p class="text-white/60 text-sm mt-1">Try searching with a different title or keyword</p>
        </div>

        <!-- Pagination Component -->
        <AppPagination
            v-if="allGames.length > perPage"
            :current-page="currentPage"
            :total-pages="totalPages"
            :total-items="allGames.length"
            :item-count="paginatedGames.length"
            @change-page="handlePageChange"
        />
    </ProductMainContainer>
</template>