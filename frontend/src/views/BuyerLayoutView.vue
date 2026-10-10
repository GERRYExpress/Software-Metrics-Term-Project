<script setup lang="ts">
import AppFooter from '@/components/AppFooter.vue';
import AppNavbar from '@/components/AppNavbar.vue';
import NotificationToast from '@/components/NotificationToast.vue';
import SessionExpiredModal from '@/components/SessionExpiredModal.vue';
import { useUserActionStore } from '@/stores/userAction';
import { computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute()
const router = useRouter()
const userActionStore = useUserActionStore()

onMounted(() => {
    userActionStore.fetchInitialData()
})

const isSessionExpired = computed(() => {
    return route.query.expired === 'true'
})

const handleSessionExpiredConfirm = () => {
    router.push({ name: 'login-view' })
}
</script>

<template>
    <div class="flex flex-col w-full min-h-screen">
        <AppNavbar/>
        <div class="flex-1">
            <RouterView/>
        </div>
        <AppFooter/>

        <SessionExpiredModal
            :show="isSessionExpired"
            @confirm="handleSessionExpiredConfirm"
        />
        <NotificationToast/>
    </div>
</template>