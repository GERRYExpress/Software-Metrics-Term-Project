<script setup lang="ts">
import type { Game } from '@/types.ts';
import AppCard from './AppCard.vue';
import BaseButton from './BaseButton.vue';
import { computed } from 'vue';

const props = defineProps<{
    games?: Game[]
}>()

const emit = defineEmits<{
    (e: 'checkout'): void
}>()

const subtotal = computed(() => {
    if (!props.games || props.games.length === 0) return 0
    return props.games.reduce((sum, g) => sum + (Number(g.price) || 0), 0)
})

const vat = computed(() => {
    return Math.ceil(subtotal.value * 7) / 100
})

const total = computed(() => {
    return subtotal.value + vat.value
})
</script>

<template>
    <AppCard class="items-start justify-start place-self-start w-auto! min-w-64 h-min gap-y-3">
        <h3 class="font-bold border-b border-white/20 pb-2 w-full">Order Summary</h3>
        <div v-if="games && games.length > 0" class="flex flex-col gap-y-1.5 w-full max-h-48 overflow-y-auto text-sm">
            <div v-for="game in games" :key="game.id" class="flex flex-row justify-between gap-x-2">
                <p class="truncate w-32">{{ game.name }}</p>
                <p class="font-semibold w-fit">${{ Number(game.price).toFixed(2) }}</p>
            </div>
        </div>
        <div class="border-t border-white/20 pt-2 w-full flex flex-col gap-y-1 text-sm">
            <div class="flex flex-row justify-between w-full">
                <p class="text-white/80">Subtotal</p>
                <p>${{ subtotal.toFixed(2) }}</p>
            </div>
            <div class="flex flex-row justify-between w-full">
                <p class="text-white/80">VAT (7%)</p>
                <p>${{ vat.toFixed(2) }}</p>
            </div>
            <div class="flex flex-row justify-between w-full font-bold text-base pt-1 border-t border-white/10">
                <p>Total</p>
                <p class="text-action">${{ total.toFixed(2) }}</p>
            </div>
        </div>
        <BaseButton :disabled="!games || games.length === 0" @click="emit('checkout')" class="mt-2">
            Checkout
        </BaseButton>
    </AppCard>
</template>