<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import BaseButton from '@/components/BaseButton.vue';
import MediaCarousel from '@/components/MediaCarousel.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import SystemRequirement from '@/components/SystemRequirement.vue';
import { useGameStore } from '@/stores/game';
import { useUserActionStore } from '@/stores/userAction';
import { Check, Heart } from '@lucide/vue';
import { Image } from '@unpic/vue';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';

const gameStore = useGameStore()
const userActionStore = useUserActionStore()

const { id, name, description, releaseDate, categories, price, media, hardwareRequirement, posterUrl, account } = storeToRefs(gameStore)

const inCart = computed(() => (id.value ? userActionStore.isInCart(id.value) : false))
const inWishlist = computed(() => (id.value ? userActionStore.isInWishlist(id.value) : false))

const formattedCategories = computed(() => {
    if (categories.value) {
        return categories.value.map(category => category.name)
    }
    return []
})
const formattedDate = computed(() => {
    if (releaseDate.value) {
        const date = new Date(releaseDate.value)
        return date.getDate() + '/' + (date.getMonth() + 1) + '/' + date.getFullYear()
    }
    return null
})

const handleAddToCart = () => {
    if (id.value) {
        userActionStore.addToCart(id.value, name.value)
    }
}

const handleWishlistAction = () => {
    if (id.value) {
        userActionStore.toggleWishlist(id.value, name.value)
    }
}
</script>

<template>
    <ProductMainContainer>
        <AppCard class="grid grid-cols-2 pt-7.5 gap-x-5 relative">
            <RouterLink class="absolute top-1 left-5" @click="$router.go(-1)" to="">Back</RouterLink>
            <Image class="rounded-xl w-full" layout="constrained" :src="posterUrl"/>
            <div class="flex flex-col place-self-start h-full gap-y-2">
                <h2>{{ name }}</h2>
                <h3>Published on <b>{{ formattedDate ?? 'XXXXXXX' }}</b> by <b>{{ account?.username ?? 'USER' }}</b></h3>
                <h3><span v-for="category in formattedCategories" :key="category">{{ category + ' ' }}</span></h3>
                <h2 class="mt-auto text-end">${{ Number(price || 0).toFixed(2) }}</h2>
                <div class="flex gap-x-2.5 w-full">
                    <!-- Add to Cart Button -->
                    <BaseButton
                        @click.prevent="handleAddToCart"
                        :type="inCart ? 'neutral' : 'action'"
                        class="w-2/3 transition-all cursor-pointer font-bold"
                        :class="inCart ? 'border-action! bg-tertiary! text-action!' : ''">
                        <span v-if="inCart" class="flex items-center justify-center gap-x-2">
                            <Check class="size-5 text-action stroke-[3]"/>
                            In Cart
                        </span>
                        <span v-else>Add to Cart</span>
                    </BaseButton>

                    <!-- Wishlist Button -->
                    <BaseButton
                        @click.prevent="handleWishlistAction"
                        type="neutral"
                        class="w-1/3! transition-all cursor-pointer"
                        :class="inWishlist ? 'border-red-500! bg-red-950/40!' : ''"
                        :title="inWishlist ? 'Remove from wishlist' : 'Add to wishlist'">
                        <Heart
                            class="size-5 transition-all"
                            :class="inWishlist ? 'fill-red-500 text-red-500 scale-110' : 'text-white'"
                        />
                    </BaseButton>
                </div>
            </div>
        </AppCard>
        <h2 class="w-2xl">Gallery</h2>
        <AppCard class="text-justify">
            <MediaCarousel :media="media"/>
        </AppCard>
        <h2 class="w-2xl">Description</h2>
        <AppCard class="text-justify">
            {{ description }}
        </AppCard>
        <h2 class="w-2xl">System Requirements</h2>
        <AppCard class="grid grid-cols-2 gap-5">
            <SystemRequirement :hardware="hardwareRequirement?.at(0)"/>
            <SystemRequirement :hardware="hardwareRequirement?.at(1)"/>
        </AppCard>
    </ProductMainContainer>
</template>