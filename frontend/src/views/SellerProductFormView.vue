<script setup lang="ts">
import AppCard from '@/components/AppCard.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseInput from '@/components/BaseInput.vue'
import GalleryUploadCard from '@/components/GalleryUploadCard.vue'
import ProductMainContainer from '@/components/ProductMainContainer.vue'
import CategoryService from '@/service/CategoryService'
import GameService from '@/service/GameService'
import { useSellerStore } from '@/stores/seller'
import { useUserActionStore } from '@/stores/userAction'
import type { Category, GameDetail, GameInput, Hardware, HardwareInput } from '@/types'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const sellerStore = useSellerStore()
const userActionStore = useUserActionStore()
const categories = ref<Category[]>([])
const isLoading = ref(true)
const isSaving = ref(false)
const errorMessage = ref('')

const form = ref<GameInput>(createEmptyForm())
const selectedCategoryIds = ref<number[]>([])
const hardwareRequirements = ref<HardwareInput[]>([
    createHardware('min'),
    createHardware('recommended')
])

const productId = computed(() => {
    const value = Number(route.params.id)
    return Number.isInteger(value) && value > 0 ? value : null
})
const isEditMode = computed(() => productId.value !== null)
const pageTitle = computed(() => isEditMode.value ? 'Edit product' : 'Add product')

onMounted(loadForm)

async function loadForm() {
    isLoading.value = true
    errorMessage.value = ''
    try {
        const seller = await sellerStore.loadSeller()
        if (!seller) throw new Error('No seller account is logged in.')
        const categoryResponse = await CategoryService.getCategories()
        categories.value = categoryResponse.data

        if (productId.value) {
            const response = await GameService.getGame(productId.value)
            const detail = response.data as GameDetail
            if (detail.account?.id !== seller.id) {
                throw new Error('You can only edit products that belong to your seller account.')
            }
            populateForm(detail, seller.id)
        } else {
            form.value.userId = seller.id
        }
    } catch (error) {
        errorMessage.value = getErrorMessage(error, 'Unable to load the product form.')
    } finally {
        isLoading.value = false
    }
}

function populateForm(detail: GameDetail, sellerId: number) {
    form.value = {
        userId: sellerId,
        name: detail.name || '',
        description: detail.description || '',
        releaseDate: detail.releaseDate || new Date(),
        price: Number(detail.price || 0),
        posterUrl: detail.posterUrl || '',
        categoryIds: detail.categories?.map(category => category.id) ?? [],
        hardwareRequirements: [],
        mediaUrls: detail.media?.map(media => media.mediaUrl).filter(Boolean) ?? []
    }
    selectedCategoryIds.value = [...form.value.categoryIds]
    hardwareRequirements.value = [
        mapHardware(detail.hardwareRequirement?.find(item => item.type === 'min'), 'min'),
        mapHardware(detail.hardwareRequirement?.find(item => item.type === 'recommended'), 'recommended')
    ]
}

async function submitForm() {
    const validationError = validateForm()
    if (validationError) {
        errorMessage.value = validationError
        return
    }

    isSaving.value = true
    errorMessage.value = ''
    form.value.categoryIds = [...selectedCategoryIds.value]
    form.value.hardwareRequirements = hardwareRequirements.value
        .filter(requirement => hasHardwareValue(requirement))

    try {
        if (productId.value) {
            await GameService.updateGame(productId.value, form.value)
            userActionStore.notify('Product updated successfully.')
        } else {
            await GameService.createGame(form.value)
            userActionStore.notify('Product created successfully.')
        }
        await router.push({ name: 'seller-dashboard-view' })
    } catch (error) {
        errorMessage.value = getErrorMessage(error, 'Unable to save the product. Your changes are still on this form.')
    } finally {
        isSaving.value = false
    }
}

function validateForm() {
    if (!form.value.name.trim()) return 'Product name is required.'
    if (!form.value.description.trim()) return 'Description is required.'
    if (!Number.isFinite(form.value.price) || form.value.price < 0) return 'Price must be zero or greater.'
    if (!form.value.posterUrl.trim()) return 'A poster URL is required.'
    if (selectedCategoryIds.value.length === 0) return 'Select at least one category.'
    if (form.value.mediaUrls.some(url => !isValidUrl(url))) return 'Every gallery URL must be a valid URL.'
    return ''
}

function createEmptyForm(): GameInput {
    return {
        userId: 0,
        name: '',
        description: '',
        releaseDate: new Date(),
        price: 0,
        posterUrl: '',
        categoryIds: [],
        hardwareRequirements: [],
        mediaUrls: []
    }
}

function createHardware(type: HardwareInput['type']): HardwareInput {
    return {
        operatingSystem: '',
        processor: '',
        memory: '',
        graphicsCard: '',
        directXVersion: '',
        network: '',
        storage: '',
        soundcard: '',
        additionalNote: '',
        type
    }
}

function mapHardware(hardware: Hardware | undefined, type: HardwareInput['type']): HardwareInput {
    return {
        operatingSystem: hardware?.operatingSystem || '',
        processor: hardware?.processor || '',
        memory: hardware?.memory || '',
        graphicsCard: hardware?.graphicsCard || '',
        directXVersion: hardware?.directXVersion || '',
        network: hardware?.network || '',
        storage: hardware?.storage || '',
        soundcard: hardware?.soundcard || '',
        additionalNote: hardware?.additionalNote || '',
        type
    }
}

function hasHardwareValue(requirement: HardwareInput) {
    return Object.entries(requirement)
        .filter(([key]) => key !== 'type')
        .some(([, value]) => String(value).trim())
}

function isValidUrl(value: string) {
    try {
        const url = new URL(value)
        return url.protocol === 'http:' || url.protocol === 'https:'
    } catch {
        return false
    }
}

function getErrorMessage(error: unknown, fallback: string) {
    const responseMessage = (error as { response?: { data?: { error?: string } } })?.response?.data?.error
    return responseMessage || (error instanceof Error ? error.message : fallback)
}
</script>

<template>
    <ProductMainContainer class="mx-auto flex w-2xl *:w-full flex-col gap-y-6 px-4 md:px-8">
        <div class="flex flex-wrap items-end justify-between gap-3">
            <div>
                <p class="text-action font-bold uppercase tracking-widest">Seller center</p>
                <h1>{{ pageTitle }}</h1>
                <p class="text-white/70">Keep your storefront information complete and accurate.</p>
            </div>
            <button type="button" class="rounded-lg bg-secondary px-4 py-2 font-bold cursor-pointer" @click="router.push({ name: 'seller-dashboard-view' })">
                Cancel
            </button>
        </div>

        <div v-if="errorMessage" class="rounded-xl border border-red/60 bg-red/10 p-4 text-red-100" role="alert">
            {{ errorMessage }}
        </div>

        <div v-if="isLoading" class="rounded-xl bg-tertiary p-10 text-center text-white/70">
            Loading product...
        </div>

        <form v-else class="flex flex-col gap-y-6" @submit.prevent="submitForm">
            <AppCard class="w-full">
                <h2 class="text-2xl">Basic information</h2>
                <div class="grid grid-cols-1 gap-4 w-full md:grid-cols-2">
                    <BaseInput v-model="form.name" label="Product name" placeholder="Enter a game name"/>
                    <BaseInput v-model.number="form.price" label="Price" type="number" min="0" step="0.01" placeholder="0.00"/>
                    <BaseInput v-model="form.releaseDate" label="Release date" type="date"/>
                    <BaseInput v-model="form.posterUrl" label="Poster URL" type="url" placeholder="https://..."/>
                </div>
                <div class="flex flex-col gap-y-2 w-full">
                    <label for="description" class="font-bold">Description</label>
                    <textarea id="description" v-model="form.description" rows="6" class="w-full border rounded-lg bg-tertiary p-3 text-white outline-none focus:ring-2 focus:ring-action" placeholder="Describe your game"/>
                </div>
                <div v-if="form.posterUrl" class="flex flex-col gap-y-2">
                    <span class="font-bold">Poster preview</span>
                    <img :src="form.posterUrl" alt="Product poster preview" class="h-52 w-36 rounded-lg object-cover" @error="errorMessage = 'The poster URL could not be loaded.'"/>
                </div>
            </AppCard>

            <AppCard class="w-full gap-y-4!">
                <div class="w-full">
                    <h2 class="text-2xl">Categories</h2>
                    <p class="text-sm text-white/60">Select at least one category.</p>
                </div>
                <div v-if="categories.length" class="grid grid-cols-1 gap-2 sm:grid-cols-2 lg:grid-cols-3 w-full">
                    <label v-for="category in categories" :key="category.id" class="flex items-center gap-2 rounded-lg bg-tertiary p-3 cursor-pointer">
                        <input v-model="selectedCategoryIds" type="checkbox" :value="category.id" class="accent-action"/>
                        <span>{{ category.name }}</span>
                    </label>
                </div>
                <p v-else class="text-white/60">No categories are available. Create one from the seller dashboard first.</p>
            </AppCard>

            <AppCard class="w-full *:w-full gap-y-4!">
                <GalleryUploadCard v-model="form.mediaUrls"/>
            </AppCard>

            <AppCard class="w-full *:w-full gap-y-4!">
                <div>
                    <h2 class="text-2xl">System requirements</h2>
                    <p class="text-sm text-white/60">Both sections are optional. Empty sections are not sent.</p>
                </div>
                <div v-for="requirement in hardwareRequirements" :key="requirement.type" class="rounded-xl bg-tertiary p-4">
                    <h3 class="mb-3 capitalize">{{ requirement.type === 'min' ? 'Minimum' : 'Recommended' }} requirements</h3>
                    <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
                        <BaseInput v-model="requirement.operatingSystem" label="Operating system" placeholder="Windows 11"/>
                        <BaseInput v-model="requirement.processor" label="Processor" placeholder="Intel Core i5"/>
                        <BaseInput v-model="requirement.memory" label="Memory" placeholder="8 GB RAM"/>
                        <BaseInput v-model="requirement.graphicsCard" label="Graphics card" placeholder="GTX 1060"/>
                        <BaseInput v-model="requirement.directXVersion" label="DirectX" placeholder="Version 12"/>
                        <BaseInput v-model="requirement.network" label="Network" placeholder="Broadband connection"/>
                        <BaseInput v-model="requirement.storage" label="Storage" placeholder="50 GB"/>
                        <BaseInput v-model="requirement.soundcard" label="Sound card" placeholder="DirectX compatible"/>
                    </div>
                    <BaseInput v-model="requirement.additionalNote" class="mt-3" label="Additional notes" placeholder="Optional notes"/>
                </div>
            </AppCard>

            <BaseButton :disabled="isSaving" type="action">
                {{ isSaving ? 'Saving...' : isEditMode ? 'Save changes' : 'Create product' }}
            </BaseButton>
        </form>
    </ProductMainContainer>
</template>
