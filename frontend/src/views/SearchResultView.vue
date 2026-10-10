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

const selectedCategory = computed(() => {
    return typeof route.query.category === 'string' ? route.query.category : ''
})

const currentPage = computed(() => {
    const p = Number(route.query.page)
    return Number.isNaN(p) || p < 1 ? 1 : p
})

const filteredGames = computed(() => {
    if (!selectedCategory.value) return allGames.value
    return allGames.value.filter(g =>
        g.categories?.some(c => c.name.toLowerCase() === selectedCategory.value.toLowerCase())
    )
})

const totalPages = computed(() => {
    return Math.max(1, Math.ceil(filteredGames.value.length / perPage))
})

const paginatedGames = computed(() => {
    const start = (currentPage.value - 1) * perPage
    return filteredGames.value.slice(start, start + perPage)
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

watch(() => [route.query.keyword, route.query.category], () => {
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

const clearFilters = () => {
    router.push({ name: 'search-result-view' })
}
</script>

<template>
    <ProductMainContainer class="flex flex-col w-5xl place-self-center">
        <div class="flex flex-row justify-between items-center w-full">
            <h1 class="text-3xl">
                {{ keyword || selectedCategory ? `Results: ${[keyword ? `"${keyword}"` : '', selectedCategory ? `[${selectedCategory}]` : ''].filter(Boolean).join(' ')}` : 'All Games' }}
            </h1>
            <button
                v-if="keyword || selectedCategory"
                @click="clearFilters"
                type="button"
                class="text-xs bg-tertiary hover:bg-secondary text-action px-3 py-1.5 rounded-lg border border-action/30 transition-all cursor-pointer">
                Clear filters
            </button>
        </div>

        <div v-if="paginatedGames.length > 0" class="grid grid-cols-4 gap-5 w-full my-4">
            <GameCard v-for="game in paginatedGames" :key="game.id" :game="game"/>
        </div>

        <div v-else-if="!isLoading" class="flex flex-col items-center justify-center py-16 bg-tertiary rounded-xl w-full text-center">
            <h2 class="text-xl">No games found</h2>
            <p class="text-white/60 text-sm mt-1">Try searching with a different title or clearing category filters</p>
        </div>

        <!-- Pagination Component -->
        <AppPagination
            v-if="filteredGames.length > perPage"
            :current-page="currentPage"
            :total-pages="totalPages"
            :total-items="filteredGames.length"
            :item-count="paginatedGames.length"
            @change-page="handlePageChange"
        />
    </ProductMainContainer>
</template>