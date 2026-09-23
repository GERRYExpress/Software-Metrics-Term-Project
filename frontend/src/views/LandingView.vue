<script setup lang="ts">
import GameCard from '@/components/GameCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';

</script>

<template>
    <ProductMainContainer>
        <h1 class="w-full">HIGHLIGHT</h1>
        <div class="bg-tertiary w-full h-64 rounded-xl"></div>
        <div class="grid grid-cols-5 gap-5 w-auto">
            <GameCard v-for="value in 10" :key="value"/>
        </div>
        <RouterLink :to="{ name: 'search-result-view' }" class="place-self-center">EXPLORE MORE</RouterLink>
    </ProductMainContainer>
</template>