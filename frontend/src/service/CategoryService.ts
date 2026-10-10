import apiClient from "./AxiosClient"

export default {
    getCategories() {
        return apiClient.get('api/v1/categories')
    }
}