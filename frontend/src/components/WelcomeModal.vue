<script setup lang="ts">
import { Gamepad2, X } from '@lucide/vue';

defineProps<{
    show: boolean
}>()

const emit = defineEmits<{
    (e: 'close'): void
}>()
</script>

<template>
    <div
        v-if="show"
        @click.self="emit('close')"
        class="fixed inset-0 bg-black/80 backdrop-blur-md flex items-center justify-center z-50 p-4 transition-all duration-300">
        <div class="bg-linear-to-br from-secondary via-primary to-tertiary rounded-3xl p-8 max-w-lg w-full shadow-2xl border border-white/20 text-center relative flex flex-col items-center gap-y-4 animate-in fade-in zoom-in duration-300">
            <button
                @click="emit('close')"
                type="button"
                class="absolute top-4 right-4 text-white/60 hover:text-white transition-all cursor-pointer">
                <X class="size-6"/>
            </button>

            <div class="bg-action/20 text-action p-4 rounded-2xl mb-1">
                <Gamepad2 class="size-12"/>
            </div>

            <h2 class="text-white text-3xl font-bold tracking-wide">
                Welcome to <br/>
                <span class="text-action">Kai-Dao-5-Fong</span>
            </h2>

            <p class="text-white/80 text-base max-w-sm font-medium">
                Your game paradise awaits. Discover, collect, and play the finest titles today!
            </p>

            <button
                @click="emit('close')"
                type="button"
                class="mt-3 bg-action hover:bg-lime-400 text-tertiary font-bold px-8 py-3 rounded-xl shadow-lg transition-all cursor-pointer hover:scale-105">
                Start Exploring
            </button>

            <p class="text-white/40 text-xs mt-1">Click anywhere to close</p>
        </div>
    </div>
</template>

