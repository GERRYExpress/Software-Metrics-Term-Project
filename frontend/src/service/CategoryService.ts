import axios from "axios";

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_BACKEND_URL,
    headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
    }
})

export default {
    getCategories() {
        return apiClient.get('api/v1/categories')
    }
}