import AccountService from '@/service/AccountService'
import type { Account } from '@/types'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSellerStore = defineStore('seller', () => {
    const account = ref<Account | null>(null)
    const isLoading = ref(false)
    const errorMessage = ref('')

    const loadSeller = async () => {
        if (account.value) return account.value

        isLoading.value = true
        errorMessage.value = ''
        try {
            const response = await AccountService.getAccounts()
            const accounts = (response.data as Array<Partial<Account> & Record<string, unknown>>).map(normalizeAccount)
            account.value = accounts.find(candidate => candidate.isSeller) ?? null
            if (!account.value) {
                throw new Error('No seller account is available.')
            }
            return account.value
        } catch (error) {
            errorMessage.value = getErrorMessage(error)
            throw error
        } finally {
            isLoading.value = false
        }
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

function normalizeAccount(raw: Partial<Account> & Record<string, unknown>): Account {
    return {
        id: Number(raw.id),
        username: String(raw.username ?? ''),
        password: '',
        createdDate: raw.createdDate ?? (raw.created_date as Date),
        isBuyer: raw.isBuyer ?? Boolean(raw.is_buyer),
        isSeller: raw.isSeller ?? Boolean(raw.is_seller),
        totalGame: raw.totalGame ?? Number(raw.total_game ?? 0),
        totalSpend: raw.totalSpend ?? Number(raw.total_spend ?? 0),
        totalSale: raw.totalSale ?? Number(raw.total_sale ?? 0)
    }
}

function getErrorMessage(error: unknown) {
    if (error instanceof Error) return error.message
    return 'Unable to load the seller account.'
}
