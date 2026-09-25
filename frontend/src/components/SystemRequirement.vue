<script setup lang="ts">
import type { Hardware } from '@/types';

defineProps<{
    hardware?: Hardware
}>()
const valueClasses = 'font-bold text-end'
</script>

<template>
    <div class="grid grid-cols-2">
        <p>OS:</p><p :class="valueClasses">{{ hardware?.operatingSystem ?? '-' }}</p>
        <p>Processor:</p><p :class="valueClasses">{{ hardware?.processor ?? '-' }}</p>
        <p>Memory:</p><p :class="valueClasses">{{ hardware?.memory ?? '-' }}</p>
        <p>Graphics:</p><p :class="valueClasses">{{ hardware?.graphicsCard ?? '-' }}</p>
        <p>DirectX:</p><p :class="valueClasses">{{ hardware?.directXVersion ?? '-' }}</p>
        <p>Network:</p><p :class="valueClasses">{{ hardware?.network ?? '-' }}</p>
        <p>Storage:</p><p :class="valueClasses">{{ hardware?.storage ?? '-' }}</p>
        <p>Soundcard:</p><p :class="valueClasses">{{ hardware?.soundcard ?? '-' }}</p>
        <p>Note:</p><p :class="valueClasses">{{ hardware?.additionalNote ?? '-' }}</p>
    </div>
</template>