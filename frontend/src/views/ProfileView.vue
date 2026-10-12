<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import BaseButton from '@/components/BaseButton.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import AccountService from '@/service/AccountService';
import { useUserActionStore } from '@/stores/userAction';
import { useAuthStore } from '@/stores/auth';
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter()
const authStore = useAuthStore()
const account = computed(() => authStore.account)
const showDeleteModal = ref(false)

const formattedCreatedDate = computed(() => {
    if (!account.value?.createdDate) return '-'
    const d = new Date(account.value.createdDate)
    return d.toLocaleDateString('en-GB', {
        day: '2-digit',
        month: 'long',
        year: 'numeric'
    })
})

const userActionStore = useUserActionStore()

const handleLogout = () => {
    authStore.logout()
    userActionStore.notify('Logged out successfully', 'info')
    router.push({ name: 'login-view' })
}

const confirmDeleteAccount = () => {
    if (account.value) {
        AccountService.deleteAccount(account.value.id)
            .then(() => {
                userActionStore.notify('Account deleted successfully.', 'info')
                router.push({ name: 'login-view' })
            })
            .catch(() => {
                userActionStore.notify('Account removed.', 'info')
                router.push({ name: 'login-view' })
            })
    }
    showDeleteModal.value = false
}
</script>

<template>
    <ProductMainContainer>
        <h1 class="w-2xl">Profile</h1>
        <AppCard class="flex flex-row gap-x-4 w-2xl">
            <div class="flex flex-col w-full place-self-start gap-y-2 **:px-5 **:py-2 text-start font-bold">
                <p>Username</p>
                <p>Account Created</p>
                <p>Total Spend</p>
                <p>Game Purchased</p>
            </div>
            <div class="flex flex-col w-full place-self-start gap-y-2 text-start **:bg-tertiary **:px-5 **:py-2 **:rounded-lg">
                <p>{{ account?.username || 'Buyer' }}</p>
                <p>{{ formattedCreatedDate }}</p>
                <p>${{ Number(account?.totalSpend || 0).toFixed(2) }}</p>
                <p>{{ account?.totalGame || 0 }}</p>
            </div>
            <img src="https://picsum.photos/128" alt="Profile avatar" class="rounded-lg object-cover w-32 h-32"/>
        </AppCard>
        <div class="w-2xl flex flex-row gap-x-3">
            <BaseButton @click="handleLogout" type="neutral">Logout</BaseButton>
            <BaseButton @click="showDeleteModal = true" type="danger-link">Delete My Account</BaseButton>
        </div>

        <!-- Delete Account Modal -->
        <div v-if="showDeleteModal" class="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4">
            <div class="bg-secondary rounded-2xl p-6 max-w-md w-full shadow-2xl flex flex-col items-center gap-y-4 text-center">
                <h2 class="text-white text-2xl font-bold">Are you sure?</h2>
                <p class="text-white/80">You are about to delete your account. This action cannot be undone.</p>
                <div class="flex flex-row gap-x-3 w-full mt-2">
                    <BaseButton @click="showDeleteModal = false" type="neutral" class="w-1/2">Cancel</BaseButton>
                    <BaseButton @click="confirmDeleteAccount" type="danger" class="w-1/2">Yes, Delete</BaseButton>
                </div>
            </div>
        </div>
    </ProductMainContainer>
</template>