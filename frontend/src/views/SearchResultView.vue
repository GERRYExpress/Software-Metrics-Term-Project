<script setup lang="ts">
import GameCard from '@/components/GameCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
defineProps<{
    page: number
}>()

</script>


<template>
    <ProductMainContainer class="flex flex-col">
        <h1 class="w-full">Search Results</h1>
        <div class="grid grid-cols-5 gap-5 w-auto">
            <GameCard v-for="value in 10" :key="value"/>
        </div>
        <div class="bg-black">
            <RouterLink v-if="page >= 1" :to="{ query: { page: page - 1 }}"/>
            <RouterLink :to="{ query: { page: page + 1 }}"/>
        </div>
    </ProductMainContainer>
</template>