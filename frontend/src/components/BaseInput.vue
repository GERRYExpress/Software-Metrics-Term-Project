<script setup lang="ts">
const modelValue = defineModel()
interface BaseInputProps {
    label?: string,
    placeholder?: string
}
withDefaults(defineProps<BaseInputProps>(), {
    label: '',
    placeholder: ''
})
</script>

<template>
    <label :for="label" class="w-full font-bold">{{ label }}</label>
    <input
        class="w-full bg-tertiary rounded-lg text-white p-2"
        :id="label"
        v-model="modelValue"
        v-bind="$attrs"
        :placeholder="placeholder"
        type="text"/>
</template>