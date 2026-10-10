<script setup lang="ts">
import { useUserActionStore } from '@/stores/userAction';
import { CheckCircle2, Heart, Info, ShoppingCart, X } from '@lucide/vue';
import { storeToRefs } from 'pinia';

const userActionStore = useUserActionStore()
const { activeNotification } = storeToRefs(userActionStore)
</script>

<template>
    <Transition
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 translate-y-6 scale-95"
        enter-to-class="opacity-100 translate-y-0 scale-100"
        leave-active-class="transition-all duration-200 ease-in"
        leave-from-class="opacity-100 translate-y-0 scale-100"
        leave-to-class="opacity-0 translate-y-4 scale-95">
        <div
            v-if="activeNotification"
            class="fixed bottom-6 right-6 z-50 flex items-center gap-x-3 bg-secondary/95 backdrop-blur-md text-white px-5 py-3.5 rounded-2xl shadow-2xl border border-action/50 max-w-sm pointer-events-auto">
            <!-- Icon -->
            <div class="p-2 rounded-xl bg-tertiary flex items-center justify-center shrink-0">
                <ShoppingCart v-if="activeNotification.type === 'cart'" class="size-5 text-action"/>
                <Heart v-else-if="activeNotification.type === 'wishlist'" class="size-5 fill-red-500 text-red-500"/>
                <CheckCircle2 v-else-if="activeNotification.type === 'success'" class="size-5 text-action"/>
                <Info v-else class="size-5 text-cyan-400"/>
            </div>

            <!-- Message Content -->
            <div class="flex-1 text-sm font-medium pr-1">
                <p class="text-white leading-tight font-semibold">
                    {{ activeNotification.message }}
                </p>
            </div>

            <!-- Close button -->
            <button
                @click="userActionStore.clearNotification"
                type="button"
                class="text-white/60 hover:text-white transition-all p-1 rounded-lg hover:bg-tertiary cursor-pointer">
                <X class="size-4"/>
            </button>
        </div>
    </Transition>
</template>

