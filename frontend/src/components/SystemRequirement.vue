<script setup lang="ts">
import type { Hardware } from '@/types';

defineProps<{
    hardware?: Hardware
}>()
const valueClasses = 'font-bold text-end'
const fields: { label: string; key: keyof Hardware }[] = [
    { label: 'OS', key: 'operatingSystem' },
    { label: 'Processor', key: 'processor' },
    { label: 'Memory', key: 'memory' },
    { label: 'Graphics', key: 'graphicsCard' },
    { label: 'DirectX', key: 'directXVersion' },
    { label: 'Network', key: 'network' },
    { label: 'Storage', key: 'storage' },
    { label: 'Soundcard', key: 'soundcard' },
    { label: 'Note', key: 'additionalNote' }
]
</script>

<template>
    <div class="grid grid-cols-2">
        <template v-for="field in fields" :key="field.key">
            <p>{{ field.label }}:</p><p :class="valueClasses">{{ hardware?.[field.key] ?? '-' }}</p>
        </template>
    </div>
</template>
