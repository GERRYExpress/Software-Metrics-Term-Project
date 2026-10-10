<script setup lang="ts">
import AppFooter from '@/components/AppFooter.vue'
import NotificationToast from '@/components/NotificationToast.vue'
import SellerNavbar from '@/components/SellerNavbar.vue'
import SessionExpiredModal from '@/components/SessionExpiredModal.vue'
import { useSellerStore } from '@/stores/seller'
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const sellerStore = useSellerStore()

const isSessionExpired = computed(() => route.query.expired === 'true')

onMounted(() => {
    sellerStore.loadSeller().catch(() => {
        // Child views display the actionable store error.
    })
})

const handleSessionExpiredConfirm = () => {
    sellerStore.clear()
    router.push({ name: 'login-view' })
}
</script>

<template>
    <div class="flex min-h-screen w-full flex-col">
        <SellerNavbar/>
        <main class="flex-1">
            <RouterView/>
        </main>
        <AppFooter/>

        <SessionExpiredModal
            :show="isSessionExpired"
            @confirm="handleSessionExpiredConfirm"
        />
        <NotificationToast/>
    </div>
</template>
