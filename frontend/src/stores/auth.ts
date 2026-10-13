import type { Account } from '@/types'
import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
    const account = ref<Account | null>(readAccount())
    const token = ref(localStorage.getItem('kd5f_token'))
    const isLoggedIn = computed(() => Boolean(token.value && account.value))

    function setSession(nextAccount: Account, nextToken: string) {
        account.value = nextAccount
        token.value = nextToken
        localStorage.setItem('kd5f_token', nextToken)
        localStorage.setItem('kd5f_account', JSON.stringify(nextAccount))
    }

    function logout() {
        account.value = null
        token.value = null
        localStorage.removeItem('kd5f_token')
        localStorage.removeItem('kd5f_account')
    }

    return { account, isLoggedIn, setSession, logout }
})

function readAccount(): Account | null {
    const raw = localStorage.getItem('kd5f_account')
    if (!raw) return null

    try {
        const parsed = JSON.parse(raw) as Partial<Account>
        return parsed.id ? parsed as Account : null
    } catch {
        localStorage.removeItem('kd5f_account')
        return null
    }
}
