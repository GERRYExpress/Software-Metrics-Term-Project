<script setup lang="ts">
import { Plus, Trash2 } from '@lucide/vue'
import { computed } from 'vue'

const props = withDefaults(defineProps<{
    modelValue: string[]
    maxItems?: number
}>(), {
    maxItems: 4
})

const emit = defineEmits<{
    (event: 'update:modelValue', value: string[]): void
}>()

const slots = computed(() => {
    const values = [...props.modelValue]
    while (values.length < props.maxItems) values.push('')
    return values.slice(0, props.maxItems)
})

const updateSlot = (index: number, value: string) => {
    const values = [...props.modelValue]
    values[index] = value
    emit('update:modelValue', values.filter(Boolean))
}

const removeSlot = (index: number) => {
    const values = [...props.modelValue]
    values.splice(index, 1)
    emit('update:modelValue', values)
}
</script>

<template>
    <section class="flex flex-col gap-y-4">
        <div>
            <h2 class="text-xl">Gallery and media</h2>
            <p class="text-sm text-white/60">Add up to {{ maxItems }} image or video URLs.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div
                v-for="(value, index) in slots"
                :key="index"
                class="flex flex-col gap-y-2 rounded-xl border border-white/10 bg-tertiary/60 p-3">
                <div class="flex items-center justify-between">
                    <label :for="`media-${index}`" class="font-bold">Media {{ index + 1 }}</label>
                    <button
                        v-if="value"
                        type="button"
                        class="text-red hover:text-red-300 cursor-pointer"
                        :aria-label="`Remove media ${index + 1}`"
                        @click="removeSlot(index)">
                        <Trash2 class="size-4"/>
                    </button>
                </div>
                <input
                    :id="`media-${index}`"
                    :value="value"
                    type="url"
                    placeholder="https://..."
                    class="w-full rounded-lg bg-secondary p-2 text-white outline-none focus:ring-2 focus:ring-action"
                    @input="updateSlot(index, ($event.target as HTMLInputElement).value)"
                />
                <img
                    v-if="value"
                    :src="value"
                    :alt="`Media preview ${index + 1}`"
                    class="h-32 w-full rounded-lg object-cover"
                    @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
                <div v-else class="flex h-32 items-center justify-center rounded-lg border border-dashed border-white/20 text-white/40">
                    <Plus class="size-6"/>
                </div>
            </div>
        </div>
    </section>
</template>
