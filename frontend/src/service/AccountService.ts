import type { LoginForm, RegisterForm } from "@/types";
import apiClient from "./AxiosClient";

export default {
    authenticate(account: LoginForm) {
        return apiClient.post('/api/v1/accounts/login', {
            username: account.username,
            password: account.password
        })
    },
    register(account: RegisterForm) {
        return apiClient.post('/api/v1/accounts', {
            username: account.username,
            password: account.password
        })
    },
    getAccounts() {
        return apiClient.get('/api/v1/accounts')
    },
    getAccount(id: number) {
        return apiClient.get(`/api/v1/accounts/${id}`)
    },
    deleteAccount(id: number) {
        return apiClient.delete(`/api/v1/accounts/${id}`)
    }
}