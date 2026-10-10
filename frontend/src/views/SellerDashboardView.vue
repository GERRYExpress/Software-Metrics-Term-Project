<script setup lang="ts">
import AppCard from '@/components/AppCard.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseInput from '@/components/BaseInput.vue'
import ProductMainContainer from '@/components/ProductMainContainer.vue'
import CategoryService from '@/service/CategoryService'
import GameService from '@/service/GameService'
import { useSellerStore } from '@/stores/seller'
import { useUserActionStore } from '@/stores/userAction'
import type { Category, Game } from '@/types'
import { Edit3, Eye, EyeOff, Plus, Trash2 } from '@lucide/vue'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const sellerStore = useSellerStore()
const userActionStore = useUserActionStore()
const categories = ref<Category[]>([])
const products = ref<Game[]>([])
const newCategoryName = ref('')
const editingCategoryId = ref<number | null>(null)
const editingCategoryName = ref('')
const isLoading = ref(true)
const isSavingCategory = ref(false)
const deletingProductId = ref<number | null>(null)
const errorMessage = ref('')

onMounted(loadDashboard)

async function loadDashboard() {
    isLoading.value = true
    errorMessage.value = ''
    try {
        const seller = await sellerStore.loadSeller()
        const [categoryResponse, productResponse] = await Promise.all([
            CategoryService.getCategories(),
            GameService.getGamesBySeller(seller.id)
        ])
        categories.value = categoryResponse.data
        products.value = productResponse.data
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

function beginEditCategory(category: Category) {
    editingCategoryId.value = category.id
    editingCategoryName.value = category.name
}

function cancelEditCategory() {
    editingCategoryId.value = null
    editingCategoryName.value = ''
}

async function saveCategory(category: Category) {
    const name = editingCategoryName.value.trim()
    if (!name) {
        errorMessage.value = 'Category name cannot be empty.'
        return
    }

    isSavingCategory.value = true
    errorMessage.value = ''
    try {
        await CategoryService.updateCategory(category.id, name, category.isHidden)
        await refreshCategories()
        cancelEditCategory()
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
</script>

<template>
    <ProductMainContainer class="mx-auto flex w-full max-w-6xl flex-col gap-y-6 px-4 md:px-8">
        <div class="flex flex-wrap items-end justify-between gap-3">
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

        <div v-if="errorMessage" class="rounded-xl border border-red/60 bg-red/10 p-4 text-red-100" role="alert">
            {{ errorMessage }}
        </div>

        <div v-if="isLoading" class="rounded-xl bg-tertiary p-10 text-center text-white/70">
            Loading seller dashboard...
        </div>

        <template v-else>
            <AppCard class="w-full gap-y-4!">
                <div class="flex flex-wrap items-center justify-between gap-3">
                    <div>
                        <h2 class="text-2xl">Categories</h2>
                        <p class="text-sm text-white/60">Control which categories buyers can see.</p>
                    </div>
                    <form class="flex w-full max-w-md gap-2" @submit.prevent="addCategory">
                        <BaseInput v-model="newCategoryName" aria-label="New category name" placeholder="New category name"/>
                        <BaseButton class="w-auto!" :disabled="isSavingCategory">
                            <Plus class="mr-1 size-4"/>
                            Add
                        </BaseButton>
                    </form>
                </div>

                <div v-if="categories.length === 0" class="rounded-lg bg-tertiary p-5 text-center text-white/60">
                    No categories have been created yet.
                </div>
                <div v-else class="overflow-x-auto">
                    <table class="w-full min-w-[640px] text-left">
                        <thead class="border-b border-white/10 text-sm text-white/60">
                            <tr>
                                <th class="p-3">Name</th>
                                <th class="p-3">Visibility</th>
                                <th class="p-3 text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="category in categories" :key="category.id" class="border-b border-white/5">
                                <td class="p-3">
                                    <input
                                        v-if="editingCategoryId === category.id"
                                        v-model="editingCategoryName"
                                        class="rounded-lg bg-tertiary p-2 text-white outline-none focus:ring-2 focus:ring-action"
                                        @keyup.enter="saveCategory(category)"
                                        @keyup.esc="cancelEditCategory"
                                    />
                                    <span v-else>{{ category.name }}</span>
                                </td>
                                <td class="p-3">
                                    <span :class="category.isHidden ? 'text-white/50' : 'text-action'">
                                        {{ category.isHidden ? 'Hidden' : 'Visible' }}
                                    </span>
                                </td>
                                <td class="flex justify-end gap-2 p-3">
                                    <template v-if="editingCategoryId === category.id">
                                        <button type="button" class="rounded-lg bg-action px-3 py-2 text-sm font-bold text-tertiary cursor-pointer" @click="saveCategory(category)">Save</button>
                                        <button type="button" class="rounded-lg bg-secondary px-3 py-2 text-sm cursor-pointer" @click="cancelEditCategory">Cancel</button>
                                    </template>
                                    <template v-else>
                                        <button type="button" class="rounded-lg bg-secondary p-2 cursor-pointer" :title="`Edit ${category.name}`" @click="beginEditCategory(category)">
                                            <Edit3 class="size-4"/>
                                        </button>
                                        <button type="button" class="rounded-lg bg-secondary p-2 cursor-pointer" :title="category.isHidden ? 'Show category' : 'Hide category'" @click="toggleCategory(category)">
                                            <EyeOff v-if="category.isHidden" class="size-4"/>
                                            <Eye v-else class="size-4"/>
                                        </button>
                                        <button type="button" class="rounded-lg bg-red p-2 cursor-pointer" :title="`Delete ${category.name}`" @click="deleteCategory(category)">
                                            <Trash2 class="size-4"/>
                                        </button>
                                    </template>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </AppCard>

            <AppCard class="w-full gap-y-4!">
                <div class="flex flex-wrap items-center justify-between gap-3">
                    <div>
                        <h2 class="text-2xl">Your products</h2>
                        <p class="text-sm text-white/60">{{ products.length }} product{{ products.length === 1 ? '' : 's' }}</p>
                    </div>
                </div>

                <div v-if="products.length === 0" class="rounded-lg bg-tertiary p-8 text-center text-white/60">
                    You have no products yet. Add your first game to start selling.
                </div>
                <div v-else class="grid grid-cols-1 gap-4 lg:grid-cols-2">
                    <article v-for="product in products" :key="product.id" class="flex gap-4 rounded-xl bg-tertiary p-4">
                        <img :src="product.posterUrl" :alt="`${product.name} poster`" class="h-32 w-24 rounded-lg object-cover"/>
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
            </AppCard>
        </template>
    </ProductMainContainer>
</template>
