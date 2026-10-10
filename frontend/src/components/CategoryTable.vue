<script setup lang="ts">
import type { Category } from '@/types'
import { Edit3, Eye, EyeOff, Trash2 } from '@lucide/vue'
import { ref } from 'vue'

defineProps<{
    categories: Category[]
}>()

const emit = defineEmits<{
    save: [category: Category, name: string]
    toggle: [category: Category]
    remove: [category: Category]
}>()

// Parent owns which row is being edited, so it can close the editor after a successful save.
const editingId = defineModel<number | null>('editingId', { default: null })
const editingName = ref('')

function beginEdit(category: Category) {
    editingId.value = category.id
    editingName.value = category.name
}

function cancelEdit() {
    editingId.value = null
    editingName.value = ''
}
</script>

<template>
    <div class="overflow-x-auto rounded-xl border border-white/20 bg-tertiary/60 shadow-lg">
        <table class="w-full min-w-[640px] text-left">
            <caption class="sr-only">Storefront categories and their visibility</caption>
            <thead class="bg-white/10 text-sm font-bold uppercase tracking-wider text-white">
                <tr class="border-b-2 border-action/50">
                    <th scope="col" class="px-4 py-4">Category name</th>
                    <th scope="col" class="px-4 py-4">Visibility</th>
                    <th scope="col" class="px-4 py-4 text-right">Actions</th>
                </tr>
            </thead>
            <tbody class="divide-y divide-white/10">
                <tr
                    v-for="category in categories"
                    :key="category.id"
                    class="odd:bg-white/[0.03] even:bg-white/[0.08] transition-colors hover:bg-white/15">
                    <td class="px-4 py-4">
                        <template v-if="editingId === category.id">
                            <label class="sr-only" :for="`category-name-${category.id}`">
                                Category name for {{ category.name }}
                            </label>
                            <input
                                :id="`category-name-${category.id}`"
                                v-model="editingName"
                                type="text"
                                class="w-full max-w-sm rounded-lg border border-action/60 bg-tertiary px-3 py-2 text-white outline-none focus:ring-2 focus:ring-action"
                                @keyup.enter="emit('save', category, editingName)"
                                @keyup.esc="cancelEdit"
                            />
                        </template>
                        <span v-else class="text-lg font-semibold text-white">{{ category.name }}</span>
                    </td>
                    <td class="px-4 py-4">
                        <span
                            class="inline-flex items-center gap-2 rounded-full border px-3 py-1 text-sm font-bold"
                            :class="category.isHidden
                                ? 'border-white/30 bg-white/10 text-white/80'
                                : 'border-action/60 bg-action/20 text-action'">
                            <span class="size-2 rounded-full" :class="category.isHidden ? 'bg-white/50' : 'bg-action'"></span>
                            {{ category.isHidden ? 'Hidden' : 'Visible' }}
                        </span>
                    </td>
                    <td class="px-4 py-4">
                        <div class="flex justify-end gap-2">
                            <template v-if="editingId === category.id">
                                <button type="button" class="cursor-pointer rounded-lg bg-action px-4 py-2 text-sm font-bold text-tertiary hover:brightness-110" @click="emit('save', category, editingName)">Save</button>
                                <button type="button" class="cursor-pointer rounded-lg bg-secondary px-4 py-2 text-sm text-white hover:brightness-125" @click="cancelEdit">Cancel</button>
                            </template>
                            <template v-else>
                                <button type="button" class="cursor-pointer rounded-lg bg-secondary p-2.5 text-white hover:brightness-125" :aria-label="`Edit ${category.name}`" :title="`Edit ${category.name}`" @click="beginEdit(category)">
                                    <Edit3 class="size-5"/>
                                </button>
                                <button type="button" class="cursor-pointer rounded-lg bg-secondary p-2.5 text-white hover:brightness-125" :aria-label="category.isHidden ? `Show ${category.name}` : `Hide ${category.name}`" :title="category.isHidden ? `Show ${category.name}` : `Hide ${category.name}`" @click="emit('toggle', category)">
                                    <EyeOff v-if="category.isHidden" class="size-5"/>
                                    <Eye v-else class="size-5"/>
                                </button>
                                <button type="button" class="cursor-pointer rounded-lg bg-red p-2.5 text-white hover:brightness-110" :aria-label="`Delete ${category.name}`" :title="`Delete ${category.name}`" @click="emit('remove', category)">
                                    <Trash2 class="size-5"/>
                                </button>
                            </template>
                        </div>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>