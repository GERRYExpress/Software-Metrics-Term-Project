<script setup lang="ts">
import ProductListCard from '@/components/ProductListCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import WishlistService from '@/service/WishlistService';
import type { Wishlist } from '@/types';
import { onMounted, ref, watchEffect } from 'vue';

const wishlists = ref<Wishlist[] | null>(null)
onMounted(() => {
    watchEffect(() => {
        WishlistService.getWishlists(1)
            .then(res => {
                wishlists.value = res.data
            })
            .catch(err => {
                console.error(err)
            })
    })
})

const handleDelete = (userId: number, gameId: number) => {
    WishlistService.deleteWishlist(userId, gameId)
        .then(() => {
            if (wishlists.value) {
                wishlists.value = wishlists.value.filter(w => w.productId !== gameId)
            }
        })
        .catch(err => console.error(err))
}
</script>


<template>
    <ProductMainContainer class="w-2xl place-self-center">
        <h1 class="w-2xl">Wishlists</h1>
        <div v-if="wishlists?.length" class="w-full *:w-full col-span-3 flex flex-col gap-y-2.5">
            <ProductListCard @delete-item="handleDelete" v-for="wishlist in wishlists" :product=wishlist :key="wishlist.game.id"/>
        </div>
        <div v-else class="flex flex-col items-center p-5 w-full justify-center opacity-75 bg-tertiary rounded-xl">
            <h2>Looks like your wishlist is empty</h2>
            <h3>Find more games <RouterLink :to="{ name: 'search-result-view' }">here</RouterLink></h3>
        </div>
    </ProductMainContainer>
</template>