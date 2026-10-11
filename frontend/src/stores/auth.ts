import type { Account } from '@/types'
import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
    const account = ref<Account | null>(readAccount())
    const isLoggedIn = computed(() => Boolean(sessionStorage.getItem('kd5f_token') && account.value))

    function setSession(nextAccount: Account, token: string) {
        account.value = nextAccount
        sessionStorage.setItem('kd5f_token', token)
        sessionStorage.setItem('kd5f_account', JSON.stringify(nextAccount))
    }

    function logout() {
        account.value = null
        sessionStorage.removeItem('kd5f_token')
        sessionStorage.removeItem('kd5f_account')
    }

    return { account, isLoggedIn, setSession, logout }
})

function readAccount(): Account | null {
    const raw = sessionStorage.getItem('kd5f_account')
    if (!raw) return null

    try {
        const parsed = JSON.parse(raw) as Partial<Account>
        return parsed.id ? parsed as Account : null
    } catch {
        sessionStorage.removeItem('kd5f_account')
        return null
    }
}
