import CartService from '@/service/CartService';
import WishlistService from '@/service/WishlistService';
import { defineStore } from 'pinia';
import { ref } from 'vue';

export interface NotificationPayload {
    id: number;
    message: string;
    type: 'cart' | 'wishlist' | 'success' | 'info';
}

export const useUserActionStore = defineStore('userAction', () => {
    const cartIds = ref<number[]>([])
    const wishlistIds = ref<number[]>([])
    const activeNotification = ref<NotificationPayload | null>(null)
    const isInitialized = ref(false)
    let timeoutId: ReturnType<typeof setTimeout> | null = null

    const notify = (message: string, type: 'cart' | 'wishlist' | 'success' | 'info' = 'success') => {
        if (timeoutId) {
            clearTimeout(timeoutId)
        }
        activeNotification.value = {
            id: Date.now(),
            message,
            type
        }
        timeoutId = setTimeout(() => {
            activeNotification.value = null
        }, 3200)
    }

    const clearNotification = () => {
        if (timeoutId) {
            clearTimeout(timeoutId)
        }
        activeNotification.value = null
    }

    const fetchInitialData = () => {
        if (isInitialized.value) return
        isInitialized.value = true

        CartService.getCarts(1)
            .then(res => {
                if (Array.isArray(res.data)) {
                    cartIds.value = res.data.map(item => item.productId ?? item.game?.id).filter(Boolean)
                }
            })
            .catch(() => {
                // If backend is offline, maintain local state
            })

        WishlistService.getWishlists(1)
            .then(res => {
                if (Array.isArray(res.data)) {
                    wishlistIds.value = res.data.map(item => item.productId ?? item.game?.id).filter(Boolean)
                }
            })
            .catch(() => {
                // If backend is offline, maintain local state
            })
    }

    const isInCart = (gameId: number) => {
        return cartIds.value.includes(gameId)
    }

    const isInWishlist = (gameId: number) => {
        return wishlistIds.value.includes(gameId)
    }

    const addToCart = (gameId: number, gameName: string) => {
        if (cartIds.value.includes(gameId)) {
            notify(`"${gameName}" is already in your cart!`, 'info')
            return
        }

        CartService.addCart(1, gameId)
            .then(() => {
                if (!cartIds.value.includes(gameId)) {
                    cartIds.value.push(gameId)
                }
                notify(`Added "${gameName}" to cart!`, 'cart')
            })
            .catch(() => {
                // Optimistically update even if offline/mock
                if (!cartIds.value.includes(gameId)) {
                    cartIds.value.push(gameId)
                }
                notify(`Added "${gameName}" to cart!`, 'cart')
            })
    }

    const toggleWishlist = (gameId: number, gameName: string) => {
        if (wishlistIds.value.includes(gameId)) {
            WishlistService.deleteWishlist(1, gameId)
                .then(() => {
                    wishlistIds.value = wishlistIds.value.filter(id => id !== gameId)
                    notify(`Removed "${gameName}" from wishlist`, 'info')
                })
                .catch(() => {
                    wishlistIds.value = wishlistIds.value.filter(id => id !== gameId)
                    notify(`Removed "${gameName}" from wishlist`, 'info')
                })
        } else {
            WishlistService.addWishlist(1, gameId)
                .then(() => {
                    wishlistIds.value.push(gameId)
                    notify(`Added "${gameName}" to wishlist!`, 'wishlist')
                })
                .catch(() => {
                    wishlistIds.value.push(gameId)
                    notify(`Added "${gameName}" to wishlist!`, 'wishlist')
                })
        }
    }

    return {
        cartIds,
        wishlistIds,
        activeNotification,
        notify,
        clearNotification,
        fetchInitialData,
        isInCart,
        isInWishlist,
        addToCart,
        toggleWishlist
    }
})

