<script setup lang="ts">
import PriceSummaryCard from '@/components/PriceSummaryCard.vue';
import ProductListCard from '@/components/ProductListCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import CartService from '@/service/CartService';
import type { Cart } from '@/types';
import { onMounted, ref, watchEffect } from 'vue';

const carts = ref<Cart[] | null>(null)
onMounted(() => {
    watchEffect(() => {
            CartService.getCarts(1)
                .then(res => carts.value = res.data)
                .catch(err => console.log(err))
    })
})

const handleDelete = (userId: number, gameId: number) => {
    CartService.deleteCart(userId, gameId)
        .then(res => alert(res.data))
}
</script>

<template>
    <ProductMainContainer class="w-2xl place-self-center">
        <h1 class="w-2xl">Carts</h1>
        <div class="grid grid-cols-4 w-full gap-x-3">
            <div v-if="carts?.length" class="w-full flex flex-col col-span-3 gap-y-2.5">
                <ProductListCard @delete-item="handleDelete" v-for="cart in carts" :product="cart" :key="cart.game.id"/>
            </div>
            <div v-else class="col-span-3 flex flex-col items-center justify-center opacity-75 bg-tertiary rounded-xl">
                <h2>Looks like your cart is empty</h2>
                <h3>Find more games <RouterLink :to="{ name: 'search-result-view' }">here</RouterLink></h3>
            </div>
            <PriceSummaryCard/>
        </div>
    </ProductMainContainer>
</template>