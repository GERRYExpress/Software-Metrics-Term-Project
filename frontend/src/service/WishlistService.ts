import apiClient from "./AxiosClient"

export default {
    getWishlists(userId: number) {
        return apiClient.get('/api/v1/wishlists/' + userId)
    },
    addWishlist(userId: number, productId: number) {
        return apiClient.post('/api/v1/wishlists/' + userId + '/' + productId)
    },
    deleteWishlist(userId: number, productId: number) {
        return apiClient.delete('/api/v1/wishlists/' + userId + '/' + productId)
    }
}