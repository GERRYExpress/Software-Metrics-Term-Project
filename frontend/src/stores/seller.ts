import type { Account } from '@/types'
import { useAuthStore } from './auth'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSellerStore = defineStore('seller', () => {
    const account = ref<Account | null>(null)
    const isLoading = ref(false)
    const errorMessage = ref('')
    const authStore = useAuthStore()

    const loadSeller = async () => {
        isLoading.value = false
        errorMessage.value = ''
        account.value = authStore.account
        if (!account.value) {
            errorMessage.value = 'Please sign in before opening the seller center.'
            throw new Error(errorMessage.value)
        }
        return account.value
    }

    const clear = () => {
        account.value = null
        errorMessage.value = ''
    }

    return {
        account,
        isLoading,
        errorMessage,
        loadSeller,
        clear
    }
})
