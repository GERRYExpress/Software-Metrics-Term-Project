<script setup lang="ts">
import { Heart, ShoppingCartPlus } from '@lucide/vue';
import BaseButton from './BaseButton.vue';
import type { Game } from '@/types.ts';
withDefaults(defineProps<{
    game?: Game | null
}>(), ({
    game: null
}))
</script>


<template>
    <RouterLink :to="{ name: '' }" class="w-auto w-max-72 rounded-xl bg-tertiary flex flex-col relative">
        <img class="w-full h-2/3 rounded-t-xl pointer-events-none" src="https://picsum.photos/128" alt="">
        <div class="bg-tertiary w-24 px-2 text-center truncate line-clamp-1 absolute place-self-end font-normal rounded-bl-xl rounded-tr-xl">
            {{ game?.categories.at(0) ?? 'CATEGORY' }}
        </div>
        <div class="p-2 w-full h-full flex flex-col gap-y-2">
            <p>{{ game?.name ?? 'NAME' }}</p>
            <p class="text-end text-xl font-bold">${{ game?.price ?? 'PRICE' }}</p>
            <div class="flex flex-row gap-x-1">
                <BaseButton>
                    <Heart/>
                </BaseButton>
                <BaseButton>
                    <ShoppingCartPlus/>
                </BaseButton>
            </div>
        </div>
    </RouterLink>
</template>