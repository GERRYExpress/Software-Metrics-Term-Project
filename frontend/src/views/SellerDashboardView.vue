<script setup lang="ts">
import AppCard from '@/components/AppCard.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseInput from '@/components/BaseInput.vue'
import ProductMainContainer from '@/components/ProductMainContainer.vue'
import CategoryService from '@/service/CategoryService'
import GameService from '@/service/GameService'
import CategoryTable from '@/components/CategoryTable.vue'
import { useSellerStore } from '@/stores/seller'
import { useUserActionStore } from '@/stores/userAction'
import type { Category, Game } from '@/types'
import { Plus } from '@lucide/vue'
import { Image } from '@unpic/vue'
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const sellerStore = useSellerStore()
const userActionStore = useUserActionStore()
const categories = ref<Category[]>([])
const products = ref<Game[]>([])
const newCategoryName = ref('')
const editingCategoryId = ref<number | null>(null)
const isLoading = ref(true)
const isSavingCategory = ref(false)
const deletingProductId = ref<number | null>(null)
const errorMessage = ref('')
const errorBanner = ref<HTMLElement | null>(null)
const currentPage = ref(1)
const totalProducts = ref(0)
const productsPerPage = 6
const totalPages = computed(() => Math.ceil(totalProducts.value / productsPerPage))
const optimizedPoster = (url: string) =>
    `https://wsrv.nl/?url=${encodeURIComponent(url)}&w=288&h=384&fit=outside&output=webp&q=70`

onMounted(loadDashboard)
watch(errorMessage, async message => {
    if (message) {
        await nextTick()
        errorBanner.value?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }
})

async function loadDashboard() {
    isLoading.value = true
    errorMessage.value = ''
    try {
        const seller = await sellerStore.loadSeller()
    if (!seller) throw new Error('No seller account is logged in.')
    const [categoryResponse, productResponse] = await Promise.all([
            CategoryService.getCategories(),
            GameService.getGamesBySeller(seller.id, productsPerPage, currentPage.value)
        ])
        categories.value = categoryResponse.data
        products.value = productResponse.data
        totalProducts.value = Number(productResponse.headers['x-total-count'] ?? 0)
    } catch (error) {
        errorMessage.value = getErrorMessage(error, 'Unable to load the seller dashboard.')
    } finally {
        isLoading.value = false
    }
}

async function addCategory() {
    const name = newCategoryName.value.trim()
    if (!name) {
        errorMessage.value = 'Enter a category name before adding it.'
        return
    }

    isSavingCategory.value = true
    errorMessage.value = ''
    try {
        await CategoryService.createCategory(name)
        newCategoryName.value = ''
        await refreshCategories()
        notify('Category added successfully.')
    } catch (error) {
        errorMessage.value = getErrorMessage(error, 'Unable to add the category.')
    } finally {
        isSavingCategory.value = false
    }
}

async function saveCategory(category: Category, rawName: string) {
    const name = rawName.trim()
    if (!name) {
        errorMessage.value = 'Category name cannot be empty.'
        return
    }

    isSavingCategory.value = true
    errorMessage.value = ''
    try {
        await CategoryService.updateCategory(category.id, name, category.isHidden)
        await refreshCategories()
        editingCategoryId.value = null
        notify('Category updated successfully.')
    } catch (error) {
        errorMessage.value = getErrorMessage(error, 'Unable to update the category.')
    } finally {
        isSavingCategory.value = false
    }
}

async function toggleCategory(category: Category) {
    isSavingCategory.value = true
    errorMessage.value = ''
    try {
        await CategoryService.updateCategory(category.id, category.name, !category.isHidden)
        await refreshCategories()
        notify(category.isHidden ? 'Category is visible again.' : 'Category hidden from buyers.')
    } catch (error) {
        errorMessage.value = getErrorMessage(error, 'Unable to change category visibility.')
    } finally {
        isSavingCategory.value = false
    }
}

async function deleteCategory(category: Category) {
    if (!window.confirm(`Delete the "${category.name}" category?`)) return

    isSavingCategory.value = true
    errorMessage.value = ''
    try {
        await CategoryService.deleteCategory(category.id)
        await refreshCategories()
        notify('Category deleted successfully.')
    } catch (error) {
        errorMessage.value = getErrorMessage(error, 'Unable to delete the category.')
    } finally {
        isSavingCategory.value = false
    }
}

async function deleteProduct(product: Game) {
    const seller = sellerStore.account
    if (!seller || !window.confirm(`Delete "${product.name}"? This cannot be undone.`)) return

    deletingProductId.value = product.id
    errorMessage.value = ''
    try {
        await GameService.deleteGame(product.id, seller.id)
        products.value = products.value.filter(item => item.id !== product.id)
        totalProducts.value = Math.max(0, totalProducts.value - 1)
        if (products.value.length === 0 && currentPage.value > 1) {
            currentPage.value--
            await loadDashboard()
        }
        notify('Product deleted successfully.')
    } catch (error) {
        errorMessage.value = getErrorMessage(error, 'Unable to delete the product.')
    } finally {
        deletingProductId.value = null
    }
}

async function refreshCategories() {
    const response = await CategoryService.getCategories()
    categories.value = response.data
}

function notify(message: string) {
    userActionStore.notify(message)
}

function getErrorMessage(error: unknown, fallback: string) {
    const responseMessage = (error as { response?: { data?: { error?: string } } })?.response?.data?.error
    return responseMessage || (error instanceof Error ? error.message : fallback)
}

async function changePage(page: number) {
    if (page < 1 || page > totalPages.value || page === currentPage.value) return
    currentPage.value = page
    await loadDashboard()
}
</script>

<template>
    <ProductMainContainer class="mx-auto flex w-2xl *:w-full flex-col gap-y-6">
        <div class="flex flex-row items-end justify-between">
            <div>
                <p class="text-action font-bold uppercase tracking-widest">Seller center</p>
                <h1>Dashboard</h1>
                <p class="text-white/70">Manage your catalog and storefront categories.</p>
            </div>
            <BaseButton class="w-auto!" @click="router.push({ name: 'seller-product-add-view' })">
                <Plus class="mr-2 size-5"/>
                Add product
            </BaseButton>
        </div>

        <div v-if="errorMessage" ref="errorBanner" class="rounded-xl border border-red/60 bg-red/10 p-4 text-red-100" role="alert">
            {{ errorMessage }}
        </div>

        <div v-if="isLoading" class="rounded-xl bg-tertiary p-10 text-center text-white/70">
            Loading seller dashboard...
        </div>

        <template v-else>


            <AppCard class="gap-y-4!">
                <div class="flex flex-wrap items-center w-full justify-between gap-3">
                    <div class="">
                        <h2 class="text-2xl">Your products</h2>
                        <p class="text-sm text-white/60">{{ totalProducts }} product{{ totalProducts === 1 ? '' : 's' }}</p>
                    </div>
                </div>

                <div v-if="products.length === 0" class="rounded-lg bg-tertiary p-8 text-center text-white/60">
                    You have no products yet. Add your first game to start selling.
                </div>
                <div v-else class="grid grid-cols-1 gap-4 lg:grid-cols-2">
                    <article v-for="product in products" :key="product.id" class="flex gap-4 rounded-xl bg-tertiary p-4">
                        <Image
                            layout="constrained"
                            loading="lazy"
                            decoding="async"
                            fetchpriority="low"
                            :src="optimizedPoster(product.posterUrl)"
                            :alt="`${product.name} poster`"
                            class="w-24! h-32 rounded-lg object-cover"
                        />
                        <div class="flex min-w-0 flex-1 flex-col">
                            <h3 class="truncate">{{ product.name }}</h3>
                            <p class="text-action">${{ Number(product.price).toFixed(2) }}</p>
                            <p class="mt-1 line-clamp-2 text-sm text-white/60">{{ product.description }}</p>
                            <div class="mt-auto flex gap-2 pt-3">
                                <RouterLink :to="{ name: 'seller-product-edit-view', params: { id: product.id } }" class="rounded-lg bg-secondary px-3 py-2 text-sm text-white! no-underline">Edit</RouterLink>
                                <button type="button" class="rounded-lg bg-red px-3 py-2 text-sm font-bold cursor-pointer disabled:opacity-50" :disabled="deletingProductId === product.id" @click="deleteProduct(product)">
                                    {{ deletingProductId === product.id ? 'Deleting...' : 'Delete' }}
                                </button>
                            </div>
                        </div>
                    </article>
                </div>
                <div v-if="totalPages > 1" class="flex items-center justify-center gap-3 text-xl">
                    <button
                        type="button"
                        class="rounded-full bg-tertiary px-4 py-2 text-base cursor-pointer disabled:pointer-events-none disabled:opacity-20"
                        :disabled="currentPage === 1"
                        @click="changePage(currentPage - 1)"
                    >
                        Prev
                    </button>
                    <span class="text-base text-white/70">Page {{ currentPage }} of {{ totalPages }}</span>
                    <button
                        type="button"
                        class="rounded-full bg-tertiary px-4 py-2 text-base cursor-pointer disabled:pointer-events-none disabled:opacity-20"
                        :disabled="currentPage === totalPages"
                        @click="changePage(currentPage + 1)"
                    >
                        Next
                    </button>
                </div>
            </AppCard>
                        <AppCard class="w-2xl gap-y-4!">
                <div class="flex flex-col w-full items-center justify-between gap-3">
                    <div class="w-full">
                        <h2 class="text-2xl">Categories</h2>
                        <p class="text-sm text-white/60">Control which categories buyers can see.</p>
                    </div>
                    <form class="flex w-full gap-2" @submit.prevent="addCategory">
                        <BaseInput class="w-full!" v-model="newCategoryName" aria-label="New category name" placeholder="New category name"/>
                        <BaseButton class="w-auto!" :disabled="isSavingCategory">
                            <Plus class="size-5"/>
                            Add
                        </BaseButton>
                    </form>
                </div>

                <div v-if="categories.length === 0" class="rounded-lg bg-tertiary p-5 text-center text-white/60">
                    No categories have been created yet.
                </div>
                    <CategoryTable
                        v-else
                        v-model:editing-id="editingCategoryId"
                        :categories="categories"
                        :is-loading="isSavingCategory"
                        @save="saveCategory"
                        @toggle="toggleCategory"
                        @remove="deleteCategory"
                    />
            </AppCard>
        </template>
    </ProductMainContainer>
</template>
