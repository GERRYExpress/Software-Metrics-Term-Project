<script setup lang="ts">
import AppCard from '@/components/AppCard.vue'
import BaseButton from '@/components/BaseButton.vue'
import ProductMainContainer from '@/components/ProductMainContainer.vue'
import AccountService from '@/service/AccountService'
import { useSellerStore } from '@/stores/seller'
import { useUserActionStore } from '@/stores/userAction'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const sellerStore = useSellerStore()
const userActionStore = useUserActionStore()
const showDeleteModal = ref(false)
const isDeleting = ref(false)
const errorMessage = ref('')

onMounted(() => {
    sellerStore.loadSeller().catch(error => {
        errorMessage.value = getErrorMessage(error, 'Unable to load seller profile.')
    })
})

const formattedCreatedDate = computed(() => {
    const value = sellerStore.account?.createdDate
    if (!value) return '-'
    return new Date(value).toLocaleDateString('en-GB', {
        day: '2-digit',
        month: 'long',
        year: 'numeric'
    })
})

function handleLogout() {
    sellerStore.clear()
    userActionStore.notify('Logged out successfully.', 'info')
    router.push({ name: 'login-view', query: { role: 'true' } })
}

async function deleteAccount() {
    const account = sellerStore.account
    if (!account) return

    isDeleting.value = true
    errorMessage.value = ''
    try {
        await AccountService.deleteAccount(account.id)
        sellerStore.clear()
        userActionStore.notify('Seller account deleted.', 'info')
        await router.push({ name: 'login-view' })
    } catch (error) {
        errorMessage.value = getErrorMessage(error, 'Unable to delete the seller account.')
    } finally {
        isDeleting.value = false
        showDeleteModal.value = false
    }
}

function getErrorMessage(error: unknown, fallback: string) {
    const responseMessage = (error as { response?: { data?: { error?: string } } })?.response?.data?.error
    return responseMessage || (error instanceof Error ? error.message : fallback)
}
</script>

<template>
    <ProductMainContainer class="mx-auto flex w-full max-w-4xl flex-col gap-y-6 px-4 md:px-8">
        <div>
            <p class="text-action font-bold uppercase tracking-widest">Seller center</p>
            <h1>Seller profile</h1>
        </div>

        <div v-if="errorMessage" class="rounded-xl border border-red/60 bg-red/10 p-4 text-red-100" role="alert">
            {{ errorMessage }}
        </div>

        <AppCard v-if="sellerStore.account" class="w-full gap-y-4!">
            <div class="flex flex-wrap items-center gap-5">
                <div class="flex size-24 items-center justify-center rounded-full bg-action text-4xl font-bold text-tertiary">
                    {{ sellerStore.account.username.charAt(0).toUpperCase() }}
                </div>
                <div>
                    <p class="text-sm uppercase tracking-widest text-white/60">Seller account</p>
                    <h2>{{ sellerStore.account.username }}</h2>
                    <p class="text-white/60">Seller role active</p>
                </div>
            </div>

            <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
                <div class="rounded-xl bg-tertiary p-4">
                    <p class="text-sm text-white/60">Account created</p>
                    <p class="font-bold">{{ formattedCreatedDate }}</p>
                </div>
                <div class="rounded-xl bg-tertiary p-4">
                    <p class="text-sm text-white/60">Products</p>
                    <p class="font-bold">{{ sellerStore.account.totalGame }}</p>
                </div>
                <div class="rounded-xl bg-tertiary p-4">
                    <p class="text-sm text-white/60">Total sales</p>
                    <p class="font-bold">${{ Number(sellerStore.account.totalSale).toFixed(2) }}</p>
                </div>
            </div>
        </AppCard>

        <div v-else-if="sellerStore.isLoading" class="rounded-xl bg-tertiary p-10 text-center text-white/70">
            Loading seller profile...
        </div>

        <div class="flex flex-wrap gap-3">
            <BaseButton class="w-auto!" type="neutral" @click="handleLogout">Logout</BaseButton>
            <BaseButton class="w-auto!" type="danger-link" @click="showDeleteModal = true">Delete my account</BaseButton>
        </div>

        <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4" role="dialog" aria-modal="true" aria-labelledby="delete-seller-title">
            <div class="flex w-full max-w-md flex-col gap-y-4 rounded-2xl bg-secondary p-6 text-center shadow-2xl">
                <h2 id="delete-seller-title" class="text-2xl">Delete seller account?</h2>
                <p class="text-white/80">This permanently removes your account and its products. This action cannot be undone.</p>
                <div class="flex gap-3">
                    <BaseButton type="neutral" @click="showDeleteModal = false">Cancel</BaseButton>
                    <BaseButton type="danger" :disabled="isDeleting" @click="deleteAccount">
                        {{ isDeleting ? 'Deleting...' : 'Yes, delete' }}
                    </BaseButton>
                </div>
            </div>
        </div>
    </ProductMainContainer>
</template>
