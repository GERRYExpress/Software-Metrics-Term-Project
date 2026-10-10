<script setup lang="ts">
import PaymentModal from '@/components/PaymentModal.vue';
import PaymentSuccessModal from '@/components/PaymentSuccessModal.vue';
import PriceSummaryCard from '@/components/PriceSummaryCard.vue';
import ProductListCard from '@/components/ProductListCard.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';
import CartService from '@/service/CartService';
import { useUserActionStore } from '@/stores/userAction';
import type { Cart } from '@/types';
import { computed, onMounted, ref, watchEffect } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter()
const carts = ref<Cart[] | null>(null)
const showPaymentModal = ref(false)
const showSuccessModal = ref(false)

const cartGames = computed(() => {
    return carts.value?.map(cart => cart.game) || []
})

const totalAmount = computed(() => {
    const subtotal = cartGames.value.reduce((sum, g) => sum + (Number(g.price) || 0), 0)
    const vat = Math.ceil(subtotal * 7) / 100
    return subtotal + vat
})

onMounted(() => {
    watchEffect(() => {
        CartService.getCarts(1)
            .then(res => carts.value = res.data)
            .catch(err => console.log(err))
    })
})

const userActionStore = useUserActionStore()

const handleDelete = (userId: number, gameId: number) => {
    CartService.deleteCart(userId, gameId)
        .then(() => {
            if (carts.value) {
                carts.value = carts.value.filter(c => c.productId !== gameId)
            }
            userActionStore.cartIds = userActionStore.cartIds.filter(id => id !== gameId)
            userActionStore.notify('Removed item from cart', 'info')
        })
        .catch(err => console.error('Failed to delete cart item:', err))
}

const handleCheckout = () => {
    if (cartGames.value.length > 0) {
        showPaymentModal.value = true
    }
}

const handlePaymentProceed = () => {
    showPaymentModal.value = false
    showSuccessModal.value = true
    // Clear carts locally
    carts.value = []
}

const handleSuccessDone = () => {
    showSuccessModal.value = false
    router.push({ name: 'landing-page-view' })
}
</script>

<template>
    <ProductMainContainer class="w-4xl place-self-center">
        <h1 class="w-full">Shopping Cart</h1>
        <div class="grid grid-cols-4 w-full gap-x-5">
            <div v-if="carts?.length" class="w-full flex flex-col col-span-3 gap-y-2.5">
                <ProductListCard @delete-item="handleDelete" v-for="cart in carts" :product="cart" :key="cart.game.id"/>
            </div>
            <div v-else class="col-span-3 flex flex-col items-center justify-center p-8 opacity-75 bg-tertiary rounded-xl text-center">
                <h2 class="text-xl">Your cart is currently empty</h2>
                <h3 class="text-sm mt-2">Find more games <RouterLink :to="{ name: 'search-result-view' }">here</RouterLink></h3>
            </div>
            <PriceSummaryCard :games="cartGames" @checkout="handleCheckout"/>
        </div>

        <!-- Checkout Modals -->
        <PaymentModal
            :show="showPaymentModal"
            :total-amount="totalAmount"
            @close="showPaymentModal = false"
            @proceed="handlePaymentProceed"
        />

        <PaymentSuccessModal
            :show="showSuccessModal"
            @done="handleSuccessDone"
        />
    </ProductMainContainer>
</template>