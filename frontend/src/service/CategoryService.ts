import apiClient from "./AxiosClient"

export default {
    getCategories() {
        return apiClient.get('/api/v1/categories')
    },
    createCategory(name: string, isHidden = false) {
        return apiClient.post('/api/v1/categories', { name, isHidden })
    },
    updateCategory(id: number, name: string, isHidden: boolean) {
        return apiClient.put(`/api/v1/categories/${id}`, { name, isHidden })
    },
    deleteCategory(id: number) {
        return apiClient.delete(`/api/v1/categories/${id}`)
    }
}