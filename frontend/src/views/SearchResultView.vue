<script setup lang="ts">
import AppPagination from '@/components/AppPagination.vue';
import GameCard from '@/components/GameCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import BaseInput from '@/components/BaseInput.vue';
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
const props = defineProps<{
    page: number
}>()
onMounted(() => {
    fetchGames()
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
        <div class="w-full text-xl flex flex-row items-center gap-x-3">
            <RouterLink v-if="page > 1" :to="{ query: { ...route.query, page: page - 1}}">Prev</RouterLink>
            <BaseInput disabled class="w-16! m-0! text-center" v-model="page"/>
            <RouterLink v-if="hasNextPage" :to="{ query: { ...route.query,  page: page + 1 }}">Next</RouterLink>
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