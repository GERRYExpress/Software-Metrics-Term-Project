<script setup lang="ts">
import type { Game } from '@/types.ts';
import AppCard from './AppCard.vue';
import BaseButton from './BaseButton.vue';
defineProps<{
    games?: Game[]
}>()

</script>


<template>
    <AppCard :class="{ 'bg-action': !games}" class="items-start justify-start place-self-start w-auto! h-min max-h-96">
        <div v-for="game in games" :key="game.id" class="flex flex-row justify-between w-full">
            <p>{{ game.name }}</p>
            <p>${{ game.price }}</p>
        </div>
        <div class="flex flex-row justify-between w-full font-bold">
            <p>Total</p>
            <p>$0.00</p>
        </div>
        <BaseButton>Checkout</BaseButton>
    </AppCard>
</template>