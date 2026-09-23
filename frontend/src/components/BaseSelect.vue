<script setup lang="ts">
import type { Category } from '@/types';
import { ChevronsUpDown } from '@lucide/vue';

const modelValue = defineModel()
interface OptionProps {
    options: Category[],
    label?: string
}
defineProps<OptionProps>()
</script>

<template>
    <label v-if="label">{{ label }}</label>
    <div class="flex flex-row">
        <select 
            class="bg-white text-tertiary! transition-all hover:scale-101 h-full pl-2.5 pr-6 text-lg align-middle appearance-none border border-gray-400 focus:border-secondary focus:outline-none"
            v-bind="$attrs" v-if="options"
            v-model="modelValue">
            <option v-for="option in options" :value="option" :key="option.id"
                class="odd:bg-tertiary/10 transition-all"
                :selected="option.id === modelValue">{{ option.name }}</option>
        </select>
        <ChevronsUpDown class="absolute top-1.5 right-0.5 pointer-events-none fill-tertiary"/>
    </div>
</template>