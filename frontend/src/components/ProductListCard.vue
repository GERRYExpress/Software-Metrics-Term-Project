<script setup lang="ts">
import AppCard from './AppCard.vue';
import BaseButton from './BaseButton.vue';

</script>



<template>
    <AppCard class="grid grid-cols-3 place-items-start! shadow-md hover:shadow-xl transition-all hover:scale-101 relative">
        <img class="w-24 rounded-xl col-span-1" src="https://picsum.photos/128/192" alt="">
        <div class="grid grid-cols-2 col-span-2 w-full">
            <div class=" w-full text-start">
                <p class="text-md font-bold">GAME 1</p>
                <p class="text-sm">CATEGORY</p>
                <p class="text-md font-thin ">Added:</p>
            </div>
            <div class="w-full text-end">
                <p>$PRICE</p>
            </div>
        </div>
        <div class="absolute bottom-0 right-0 w-24!">
            <BaseButton type="danger-link">Remove</BaseButton>
        </div>
    </AppCard>

</template>