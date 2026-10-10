<script setup lang="ts">
import { computed, ref } from 'vue';
import { ChevronLeft, ChevronRight } from '@lucide/vue';

const props = withDefaults(defineProps<{
    currentPage?: number,
    totalPages?: number,
    totalItems?: number,
    itemCount?: number
}>(), {
    currentPage: 1,
    totalPages: 1,
    totalItems: 0,
    itemCount: 0
})

const emit = defineEmits<{
    (e: 'changePage', page: number): void
}>()

const jumpPage = ref<number | ''>('')

const visiblePages = computed(() => {
    const pages: number[] = []
    const total = Math.max(1, props.totalPages)
    const current = Math.max(1, Math.min(props.currentPage, total))

    let start = Math.max(1, current - 2)
    let end = Math.min(total, current + 2)

    if (current <= 3) {
        end = Math.min(total, 5)
    }
    if (current >= total - 2) {
        start = Math.max(1, total - 4)
    }

    for (let i = start; i <= end; i++) {
        pages.push(i)
    }
    return pages
})

const goToPage = (page: number) => {
    if (page >= 1 && page <= props.totalPages && page !== props.currentPage) {
        emit('changePage', page)
    }
}

const handleJump = () => {
    if (typeof jumpPage.value === 'number') {
        const target = Math.max(1, Math.min(props.totalPages, jumpPage.value))
        goToPage(target)
        jumpPage.value = ''
    }
}
</script>

<template>
    <div class="flex flex-col items-center justify-center gap-y-4 w-full my-6">
        <!-- Pagination Buttons Row -->
        <div class="flex flex-row items-center justify-center gap-x-2">
            <!-- Prev Button -->
            <button
                @click="goToPage(currentPage - 1)"
                :disabled="currentPage <= 1"
                type="button"
                class="size-11 rounded-full flex items-center justify-center transition-all cursor-pointer font-bold border"
                :class="currentPage <= 1 ? 'border-white/10 text-white/30 cursor-not-allowed bg-tertiary/40' : 'border-action/40 bg-tertiary text-white hover:bg-secondary hover:border-action'">
                <ChevronLeft class="size-5"/>
            </button>

            <!-- Page Number Circles -->
            <button
                v-for="p in visiblePages"
                :key="p"
                @click="goToPage(p)"
                type="button"
                class="size-11 rounded-full flex items-center justify-center text-base font-bold transition-all cursor-pointer shadow-md"
                :class="p === currentPage ? 'bg-action text-tertiary scale-110' : 'bg-tertiary text-white hover:bg-secondary'">
                {{ p }}
            </button>

            <!-- Next Button -->
            <button
                @click="goToPage(currentPage + 1)"
                :disabled="currentPage >= totalPages"
                type="button"
                class="size-11 rounded-full flex items-center justify-center transition-all cursor-pointer font-bold border"
                :class="currentPage >= totalPages ? 'border-white/10 text-white/30 cursor-not-allowed bg-tertiary/40' : 'border-action/40 bg-tertiary text-white hover:bg-secondary hover:border-action'">
                <ChevronRight class="size-5"/>
            </button>
        </div>

        <!-- Jump to page & Counter row -->
        <div class="flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-sm text-white/80">
            <p v-if="totalItems > 0">
                Displaying <span class="text-action font-bold">{{ itemCount }}</span> out of <span class="text-action font-bold">{{ totalItems }}</span> games
            </p>

            <form @submit.prevent="handleJump" class="flex flex-row items-center gap-x-2">
                <span>Go to page:</span>
                <input
                    v-model.number="jumpPage"
                    type="number"
                    min="1"
                    :max="totalPages"
                    placeholder="#"
                    class="w-14 bg-tertiary text-white rounded-lg px-2 py-1 text-center text-sm border border-white/20 focus:outline-action"
                />
                <button
                    type="submit"
                    class="bg-action text-tertiary font-bold px-3 py-1 rounded-lg text-xs hover:bg-lime-400 transition-all cursor-pointer">
                    GO
                </button>
            </form>
        </div>
    </div>
</template>

