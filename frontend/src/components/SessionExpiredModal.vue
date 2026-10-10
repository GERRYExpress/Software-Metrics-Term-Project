<script setup lang="ts">
import { Clock } from '@lucide/vue';
import BaseButton from './BaseButton.vue';

defineProps<{
    show: boolean
}>()

const emit = defineEmits<{
    (e: 'confirm'): void
}>()
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-black/80 backdrop-blur-xs flex items-center justify-center z-50 p-4">
        <div class="bg-secondary rounded-2xl p-6 max-w-sm w-full shadow-2xl flex flex-col items-center gap-y-4 text-center border border-white/10 animate-in fade-in zoom-in duration-200">
            <div class="text-action">
                <Clock class="size-16"/>
            </div>
            <h2 class="text-white text-2xl font-bold">Session Expired</h2>
            <p class="text-white/80 text-sm">
                You have been idle for too long. Please log in again to continue.
            </p>
            <BaseButton @click="emit('confirm')" class="w-full mt-2">
                OK
            </BaseButton>
        </div>
    </div>
</template>

