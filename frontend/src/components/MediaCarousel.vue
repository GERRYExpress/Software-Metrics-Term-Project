<script setup lang="ts">
import type { Media } from '@/types';
import { ChevronLeft, ChevronRight, Image as ImageIcon, Play } from '@lucide/vue';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
    media: Media[] | null,
}>()

const activeIndex = ref(0)

const mediaList = computed(() => {
    return props.media && props.media.length > 0 ? props.media : []
})

watch(() => props.media, () => {
    activeIndex.value = 0
})

const prev = () => {
    if (mediaList.value.length === 0) return
    activeIndex.value = (activeIndex.value - 1 + mediaList.value.length) % mediaList.value.length
}

const next = () => {
    if (mediaList.value.length === 0) return
    activeIndex.value = (activeIndex.value + 1) % mediaList.value.length
}

const currentMedia = computed(() => {
    return mediaList.value[activeIndex.value] || null
})

</script>

<template>
    <div class="flex flex-col gap-y-3 w-full items-center">
        <!-- Main Media Preview -->
        <div class="relative w-full h-80 bg-tertiary rounded-xl overflow-hidden flex items-center justify-center group">
            <template v-if="currentMedia">
                <img
                    v-if="currentMedia.type === 'image'"
                    :src="currentMedia.mediaUrl"
                    alt="Gallery item"
                    class="w-full h-full object-contain transition-all duration-300"
                />
                <iframe
                    v-else
                    frameborder="0"
                    :src="currentMedia.mediaUrl"
                    title="Game trailer"
                    allow="accelerometer; autoplay; encrypted-media; clipboard-write; gyroscope; picture-in-picture"
                    allowfullscreen
                    class="w-full h-full object-contain"
                />
            </template>
            <div v-else class="flex flex-col items-center justify-center text-white/40 gap-y-2">
                <ImageIcon class="size-16"/>
                <p class="text-sm">No gallery media available</p>
            </div>

            <!-- Carousel Controls -->
            <button
                v-if="mediaList.length > 1"
                @click="prev"
                type="button"
                class="absolute left-2 top-1/2 -translate-y-1/2 bg-black/60 hover:bg-black/90 text-white p-2 rounded-full transition-all cursor-pointer">
                <ChevronLeft class="size-6"/>
            </button>
            <button
                v-if="mediaList.length > 1"
                @click="next"
                type="button"
                class="absolute right-2 top-1/2 -translate-y-1/2 bg-black/60 hover:bg-black/90 text-white p-2 rounded-full transition-all cursor-pointer">
                <ChevronRight class="size-6"/>
            </button>

            <!-- Counter Badge -->
            <div v-if="mediaList.length > 0" class="absolute bottom-2 right-2 bg-black/70 px-2.5 py-1 rounded-md text-xs text-white">
                {{ activeIndex + 1 }} / {{ mediaList.length }}
            </div>
        </div>

        <!-- Thumbnails -->
        <div v-if="mediaList.length > 1" class="flex flex-row gap-x-2 overflow-x-auto w-full pb-1">
            <button
                v-for="(item, idx) in mediaList"
                :key="item.id || idx"
                @click="activeIndex = idx"
                type="button"
                class="shrink-0 w-20 h-14 rounded-lg overflow-hidden border-2 transition-all cursor-pointer"
                :class="activeIndex === idx ? 'border-action scale-105' : 'border-transparent opacity-60 hover:opacity-100'">
                <img
                    v-if="item.type === 'image'"
                    :src="item.mediaUrl"
                    alt="Thumbnail"
                    class="w-full h-full object-cover"
                />
                <span v-else class="flex h-full items-center justify-center bg-black/60 text-white">
                    <Play class="size-5" />
                </span>
            </button>
        </div>
    </div>
</template>