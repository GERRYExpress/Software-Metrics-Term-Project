<script setup lang="ts">
</script>
<template>
    <button class="bg-action flex flex-row items-center justify-center hover:bg-lime-400 transition-all w-full px-5 py-3 rounded-xl shadow-md hover:shadow-xl cursor-pointer text-tertiary font-bold" v-bind="$attrs">
        <slot/>
    </button>
</template>