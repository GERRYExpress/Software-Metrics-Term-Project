<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import BaseButton from '@/components/BaseButton.vue';
import MediaCarousel from '@/components/MediaCarousel.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import SystemRequirement from '@/components/SystemRequirement.vue';
import { useGameStore } from '@/stores/game';
import { Heart } from '@lucide/vue';
import { Image } from '@unpic/vue';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';
const gameStore = useGameStore()
const { name, description, releaseDate, categories, price, media, hardwareRequirement, posterUrl, account } = storeToRefs(gameStore)

const formattedCategories = computed(() => {
    if (categories.value) {
        return categories.value.map(category => category.name)
    }
})
const formattedDate = computed(() => {
    if (releaseDate.value) {
        const date = new Date(releaseDate.value)
        return date.getDate() + '/' + date.getMonth() + '/' + date.getFullYear()
    }
})
</script>



<template>
    <ProductMainContainer>
        <AppCard class="grid grid-cols-2 gap-x-5">
            <Image class="rounded-xl w-full" layout="constrained" :src="posterUrl"/>
            <div class="flex flex-col place-self-start h-full gap-y-2">
                <h2>{{ name }}</h2>
                <h3>Published on <b>{{ formattedDate ?? 'XXXXXXX' }}</b> by <b>{{ account?.username ?? 'USER' }}</b></h3>
                <h3><span v-for="category in formattedCategories" :key="category">{{ category + ' ' }}</span></h3>
                <h2 class="mt-auto text-end">${{ price.toFixed(2) }}</h2>
                <div class="flex gap-x-2.5">
                    <BaseButton class="w-2/3">Add to Cart</BaseButton>
                    <BaseButton type="neutral" class="w-1/3!"><Heart/></BaseButton>
                </div>
            </div>
        </AppCard>
        <h2 class="w-2xl">Gallery</h2>
        <AppCard class="text-justify">
            <!-- CAROUSEL MODIFIACTION NEEDED --Krittameth to Boki Chan -->
            <MediaCarousel :media/>
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