<script setup lang="ts">
import { computed } from 'vue'

const modelValue = defineModel<boolean>()
const emit = defineEmits(['toggleRole'])
defineProps({
    options: {
        type: Array,
        required: true,
        validator(value: string[]) {
            if (value.length !== 2) {
                console.error(`Error: 'item' props must have 2 items, but received ${value.length}`)
                return false
            }
            return true
        }
    }
})
const handleClick = () => {
    emit('toggleRole')
}
const hoverClass = 'transition-all cursor-pointer shadow-sm hover:shadow-lg hover:shadow-lg'
const buttonClass = 'w-full text-center transition-all select-none p-2 font-medium'

const selectionClass = computed(() =>{
    return {
        first: modelValue.value ? 'bg-action text-tertiary scale-115 -translate-x-1' : 'bg-gray-100 text-gray-300',
        second: !modelValue.value ? 'bg-action text-tertiary scale-115 translate-x-1' : 'bg-gray-100 text-gray-300'
}
}) 
</script>


<template>
    <div @click="handleClick" :class="hoverClass" class="flex flex-row scale-80 w-30 place-self-end justify-center items-center rounded-xl">
        <div :class="[selectionClass.first, buttonClass, 'rounded-s-xl']">{{ options.at(0) }}</div>
        <div :class="[selectionClass.second, buttonClass, 'rounded-e-xl']">{{ options.at(1) }}</div>
    </div>
</template>