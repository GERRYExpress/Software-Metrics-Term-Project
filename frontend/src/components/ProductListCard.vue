<script setup lang="ts">
import type { Cart, Wishlist } from '@/types.ts';
import AppCard from './AppCard.vue';
import BaseButton from './BaseButton.vue';
import { computed } from 'vue';
const props = defineProps<{
    product: Wishlist | Cart
}>()
const emit = defineEmits(['delete-item'])

const categories = computed(() => {
    return props.product.game.categories.map(category => category.name).join(', ')
})

const date = computed(() => {
    if ('wishlistedDate' in props.product) {
        const wishlistDate = new Date(props.product.wishlistedDate)
        return wishlistDate.getDate() + '/' + wishlistDate.getMonth() + '/' + wishlistDate.getFullYear()
    }
    if ('cartDate' in props.product) {
        const cartDate = new Date(props.product.cartDate)
        return cartDate.getDate() + '/' + cartDate.getMonth() + '/' + cartDate.getFullYear()
    }
    return null
})

const handleDelete = () => {
    emit('delete-item', props.product.userId, props.product.productId)
}
</script>



<template>
    <AppCard class="flex flex-row place-items-start! w-lg shadow-md hover:shadow-xl transition-all hover:scale-101 relative">
        <img class="w-24 rounded-xl" :src="product.game.posterUrl" alt="">
        <div class="grid grid-cols-4 col-span-2 w-full h-full">
            <div class="*:w-full flex flex-col mx-5 col-span-3">
                <p class="text-md font-bold">{{ product.game.name }}</p>
                <p class="text-sm line-clamp-1">{{ categories }}</p>
                <p class="text-md font-thin ">Added: {{ date ?? '' }}</p>
            </div>
            <div class="w-full text-end font-medium">
                <p>${{ product.game.price }}</p>
            </div>
        </div>
        <div class="absolute bottom-0 right-0 w-24!">
            <BaseButton @click.prevent="handleDelete" type="danger-link">Remove</BaseButton>
        </div>
    </AppCard>
</template>