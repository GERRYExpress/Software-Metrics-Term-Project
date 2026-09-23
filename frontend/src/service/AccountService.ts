import type { LoginForm } from "@/types";
import axios from "axios";

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_BACKEND_URL,
    headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
    }
})

export default {
    authenticate(account: LoginForm) {
        return apiClient.post('/api/v1/accounts/authenticate', account)
    }
}