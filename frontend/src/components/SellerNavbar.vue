<script setup lang="ts">
import { LayoutDashboard, PackagePlus, Store, UserRound } from '@lucide/vue'
import { useSellerStore } from '@/stores/seller'
import { onMounted } from 'vue'

const sellerStore = useSellerStore()

onMounted(() => {
    sellerStore.loadSeller().catch(() => {
        // The layout renders its error state from the store.
    })
})
</script>

<template>
    <nav class="mx-auto flex w-full max-w-6xl flex-wrap items-center justify-between gap-3 px-4 py-5 md:px-8">
        <RouterLink :to="{ name: 'seller-dashboard-view' }" class="flex items-center gap-x-2 text-xl no-underline">
            <Store class="size-7"/>
            <span>KD5F Seller</span>
        </RouterLink>

        <div class="flex flex-wrap items-center justify-center gap-2">
            <RouterLink
                :to="{ name: 'seller-dashboard-view' }"
                class="flex items-center gap-x-2 rounded-full bg-secondary px-4 py-2 text-sm text-white! no-underline hover:bg-tertiary">
                <LayoutDashboard class="size-4"/>
                Dashboard
            </RouterLink>
            <RouterLink
                :to="{ name: 'seller-product-add-view' }"
                class="flex items-center gap-x-2 rounded-full bg-secondary px-4 py-2 text-sm text-white! no-underline hover:bg-tertiary">
                <PackagePlus class="size-4"/>
                Add product
            </RouterLink>
            <RouterLink
                :to="{ name: 'seller-profile-view' }"
                class="flex items-center gap-x-2 rounded-full bg-secondary px-4 py-2 text-sm text-white! no-underline hover:bg-tertiary">
                <UserRound class="size-4"/>
                {{ sellerStore.account?.username || 'Profile' }}
            </RouterLink>
            <RouterLink
                :to="{ name: 'landing-page-view' }"
                class="rounded-full border border-action/50 px-4 py-2 text-sm text-action! no-underline hover:bg-action hover:text-tertiary!">
                Buyer storefront
            </RouterLink>
        </div>
    </nav>
</template>
