import type { LoginForm } from "@/types";
import apiClient from "./AxiosClient";

export default {
    authenticate(account: LoginForm) {
        return apiClient.post('/api/v1/accounts/authenticate', account)
    }
}