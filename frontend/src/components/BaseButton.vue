<script setup lang="ts">
interface ButtonProps {
    type?: 'danger' | 'danger-link' | 'neutral' | 'action'
}
withDefaults(defineProps<ButtonProps>(), {
    type: 'action'
})
</script>
<template>
    <button class="flex flex-row items-center justify-center transition-all w-full px-5 py-3 rounded-xl shadow-md hover:shadow-xl cursor-pointer font-bold"
        :class="{
            'bg-secondary border-white border-2 hover:outline-3 text-white': type === 'neutral',
            'bg-red hover:bg-amber-950 text-white': type === 'danger',
            'bg-action hover:bg-lime-400 text-tertiary' : type === 'action',
            'bg-none underline text-red hover:scale-105 hover:shadow-none! shadow-none' : type === 'danger-link'
            }"
        v-bind="$attrs">
        <slot/>
    </button>
</template>