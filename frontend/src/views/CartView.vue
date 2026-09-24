<script setup lang="ts">
import PriceSummaryCard from '@/components/PriceSummaryCard.vue';
import ProductListCard from '@/components/ProductListCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';

</script>

<template>
    <ProductMainContainer class="w-2xl place-self-center">
        <h1 class="w-2xl">Cart</h1>
        <div class="grid grid-cols-4 w-full gap-x-3">
            <div class="w-full col-span-3 flex flex-col gap-y-2.5">
                <ProductListCard v-for="value in 5" :key="value"/>
            </div>
            <PriceSummaryCard/>
        </div>
    </ProductMainContainer>
</template>