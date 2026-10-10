<script setup lang="ts">
import BaseButton from './BaseButton.vue';
import { CheckCircle2 } from '@lucide/vue';

defineProps<{
    show: boolean
}>()

const emit = defineEmits<{
    (e: 'done'): void
}>()
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-black/75 backdrop-blur-xs flex items-center justify-center z-50 p-4">
        <div class="bg-secondary rounded-2xl p-6 max-w-sm w-full shadow-2xl flex flex-col items-center gap-y-4 text-center border border-white/10 animate-in fade-in zoom-in duration-200">
            <CheckCircle2 class="size-16 text-action"/>
            <h2 class="text-white text-2xl font-bold">Payment Successful!</h2>
            <p class="text-white/80 text-sm">
                Thank you for your purchase. Your games have been added to your library!
            </p>
            <BaseButton @click="emit('done')" class="w-full mt-2">
                Continue Shopping
            </BaseButton>
        </div>
    </div>
</template>

