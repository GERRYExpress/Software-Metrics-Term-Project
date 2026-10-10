<script setup lang="ts">
import { Check, Heart, ShoppingCartPlus } from '@lucide/vue';
import BaseButton from './BaseButton.vue';
import type { Game } from '@/types.ts';
import { Image } from '@unpic/vue';
import { useUserActionStore } from '@/stores/userAction';
import { computed } from 'vue';

const props = defineProps<{
    game: Game
}>()

const userActionStore = useUserActionStore()

const inCart = computed(() => userActionStore.isInCart(props.game.id))
const inWishlist = computed(() => userActionStore.isInWishlist(props.game.id))

const optimizedPoster = (url: string) =>
    `https://wsrv.nl/?url=${encodeURIComponent(url)}&w=288&h=192&fit=outside&output=webp&q=70`

const handleAddToCart = () => {
    userActionStore.addToCart(props.game.id, props.game.name)
}

const handleWishlist = () => {
    userActionStore.toggleWishlist(props.game.id, props.game.name)
}
</script>

<template>
    <RouterLink :to="{ name: 'game-detail-view', params: { id: game.id } }" class="w-auto w-max-72 rounded-xl bg-tertiary flex flex-col relative no-underline! transition-all hover:scale-102 hover:shadow-xl">
        <Image
            layout="constrained"
            loading="lazy"
            decoding="async"
            fetchpriority="low"
            class="w-full h-48 rounded-t-xl pointer-events-none object-cover object-top"
            :src="optimizedPoster(game.posterUrl)"
            alt=""
        />
        <div class="bg-tertiary w-32 px-2 text-center text-sm truncate line-clamp-1 absolute place-self-end font-normal rounded-bl-xl rounded-tr-xl border-b border-l border-white/10">
            {{ game.categories?.map(category => category.name).join(', ') ?? 'CATEGORY' }}
        </div>
        <div class="p-2 w-full h-full flex flex-col gap-y-2">
            <p class="line-clamp-2 h-full font-semibold">{{ game.name ?? 'NAME' }}</p>
            <p class="text-end text-xl font-bold">${{ Number(game.price || 0).toFixed(2) }}</p>
            <div class="flex flex-row gap-x-1.5">
                <!-- Wishlist Button -->
                <BaseButton
                    @click.prevent.stop="handleWishlist"
                    type="neutral"
                    class="transition-all cursor-pointer"
                    :class="inWishlist ? 'border-red-500! bg-red-950/40! shadow-red-500/20' : ''"
                    :title="inWishlist ? 'Remove from wishlist' : 'Add to wishlist'">
                    <Heart
                        class="size-4 transition-all"
                        :class="inWishlist ? 'fill-red-500 text-red-500 scale-110' : 'text-white'"
                    />
                </BaseButton>

                <!-- Add to Cart Button -->
                <BaseButton
                    @click.prevent.stop="handleAddToCart"
                    :type="inCart ? 'neutral' : 'action'"
                    class="transition-all cursor-pointer font-bold"
                    :class="inCart ? 'border-action! bg-tertiary! text-action!' : ''"
                    :title="inCart ? 'Already in cart' : 'Add to cart'">
                    <span v-if="inCart" class="flex items-center gap-x-1 text-xs">
                        <Check class="size-4 text-action stroke-3"/>
                        In Cart
                    </span>
                    <ShoppingCartPlus v-else class="size-4"/>
                </BaseButton>
            </div>
        </div>
    </RouterLink>
</template>