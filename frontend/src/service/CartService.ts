import apiClient from "./AxiosClient"

export default {
    getCarts(userId: number) {
        return apiClient.get('/api/v1/carts/' + userId)
    },
    addCart(userId: number, productId: number) {
        return apiClient.post('/api/v1/carts/' + userId + '/' + productId)
    },
    deleteCart(userId: number, productId: number) {
        return apiClient.delete('/api/v1/carts/' + userId + '/' + productId)
    }
}