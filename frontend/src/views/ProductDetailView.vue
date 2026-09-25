<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import BaseButton from '@/components/BaseButton.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import SystemRequirement from '@/components/SystemRequirement.vue';
import { useGameStore } from '@/stores/game';
import { Heart } from '@lucide/vue';
import { storeToRefs } from 'pinia';
const gameStore = useGameStore()
const { name, description, releaseDate, categories, price, media, hardwareRequirement, posterUrl, account } = storeToRefs(gameStore)
</script>



<template>
    <ProductMainContainer>
        <AppCard class="grid grid-cols-2">
            <div>CAROUSEL {{ media }} {{ posterUrl }}</div>
            <div class="flex flex-col">
                <h2>{{ name }}</h2>
                <h3>Published on <b>{{ releaseDate ?? 'XX/XX/XXXX' }}</b> by <b>{{ account?.username ?? 'USER' }}</b></h3>
                <h3>{{ categories }}</h3>
                <h2>${{ price.toFixed(2) }}</h2>
                <div class="flex gap-x-2.5">
                    <BaseButton class="w-2/3">Add to Cart</BaseButton>
                    <BaseButton type="neutral" class="w-1/3!"><Heart/></BaseButton>
                </div>
            </div>
        </AppCard>
        <h2 class="w-2xl">Description</h2>
        <AppCard class="text-justify">
            {{ description }}
            Lorem ipsum, dolor sit amet consectetur adipisicing elit. Quisquam quas unde adipisci beatae placeat dolore at libero nisi hic doloribus magnam omnis ab, aliquid, corrupti voluptatum nostrum ut? Laudantium, quisquam.
        </AppCard>
        <h2 class="w-2xl">System Requirements</h2>
        <AppCard class="grid grid-cols-2 gap-5">
            <SystemRequirement :hardware="hardwareRequirement?.at(0)"/>
            <SystemRequirement :hardware="hardwareRequirement?.at(1)"/>
        </AppCard>
    </ProductMainContainer>
</template>