<script setup lang="ts">
const modelValue = defineModel<string>()
interface BaseInputProps {
    label?: string,
    placeholder?: string,
    type?: string
}
withDefaults(defineProps<BaseInputProps>(), {
    label: '',
    placeholder: '',
    type: 'text'
})
</script>

<template>
    <label v-if="label" :for="label" class="w-full font-bold">{{ label }}</label>
    <input
        class="w-full bg-tertiary rounded-lg text-white p-2"
        :id="label"
        v-model="modelValue"
        v-bind="$attrs"
        :placeholder="placeholder"
        :type="type"/>
</template>