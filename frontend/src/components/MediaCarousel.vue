<script setup lang="ts">
import type { Media } from '@/types';
import { computed } from 'vue';
const props = defineProps<{
    media: Media[] | null,
    posterUrl: string
}>()

const posterUrl = computed(() => {
    return props.posterUrl === '' ? 'https://placehold.co/128x192?text=Poster+Unavailable' : props.posterUrl
})
</script>

<template>
    <img class="rounded-xl" :src="posterUrl"/>
</template>