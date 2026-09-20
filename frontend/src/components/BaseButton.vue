<script setup lang="ts">
interface BaseButtonProp {
    text: string
}
withDefaults(defineProps<BaseButtonProp>(), {
    text: ''
}
)
</script>
<template>
    <button class="bg-action hover:bg-lime-400 transition-all w-full px-5 py-3 rounded-xl shadow-md hover:shadow-xl cursor-pointer text-tertiary font-bold" v-bind="$attrs">
        {{ text }}
    </button>
</template>