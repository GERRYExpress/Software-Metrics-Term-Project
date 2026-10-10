import type { LoginForm } from "@/types";
import apiClient from "./AxiosClient";

export default {
    authenticate(account: LoginForm) {
        return apiClient.post('/api/v1/accounts/authenticate', account)
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