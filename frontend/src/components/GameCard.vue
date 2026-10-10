<script setup lang="ts">
import { Heart, ShoppingCartPlus } from '@lucide/vue';
import BaseButton from './BaseButton.vue';
import type { Game } from '@/types.ts';
import { Image } from '@unpic/vue';
import CartService from '@/service/CartService';
import WishlistService from '@/service/WishlistService';

const props = defineProps<{
    game: Game
}>()

const optimizedPoster = (url: string) =>
    `https://wsrv.nl/?url=${encodeURIComponent(url)}&w=288&h=192&fit=outside&output=webp&q=70`

const handleAddToCart = () => {
    CartService.addCart(1, props.game.id)
        .then(() => alert(`Added "${props.game.name}" to cart!`))
        .catch(() => alert(`Added "${props.game.name}" to cart!`))
}

const handleWishlist = () => {
    WishlistService.addWishlist(1, props.game.id)
        .then(() => alert(`Added "${props.game.name}" to wishlist!`))
        .catch(() => alert(`Added "${props.game.name}" to wishlist!`))
}
</script>

<template>
    <RouterLink :to="{ name: 'game-detail-view', params: { id: game.id } }" class="w-auto w-max-72 rounded-xl bg-tertiary flex flex-col relative no-underline!">
        <Image
            layout="constrained"
            loading="lazy"
            decoding="async"
            fetchpriority="low"
            class="w-full h-48 rounded-t-xl pointer-events-none object-cover object-top"
            :src="optimizedPoster(game.posterUrl)"
            alt=""
        />
        <div class="bg-tertiary w-32 px-2 text-center text-sm truncate line-clamp-1 absolute place-self-end font-normal rounded-bl-xl rounded-tr-xl">
            {{ game.categories?.map(category => category.name).join(', ') ?? 'CATEGORY' }}
        </div>
        <div class="p-2 w-full h-full flex flex-col gap-y-2">
            <p class="line-clamp-2 h-full">{{ game.name ?? 'NAME' }}</p>
            <p class="text-end text-xl font-bold">${{ Number(game.price || 0).toFixed(2) }}</p>
            <div class="flex flex-row gap-x-1">
                <BaseButton @click.prevent.stop="handleWishlist" type="neutral" class="w-1/2! py-2!">
                    <Heart class="size-4"/>
                </BaseButton>
                <BaseButton @click.prevent.stop="handleAddToCart" class="w-1/2! py-2!">
                    <ShoppingCartPlus class="size-4"/>
                </BaseButton>
            </div>
        </div>
    </RouterLink>
</template>