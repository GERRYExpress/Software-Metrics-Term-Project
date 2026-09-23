<script setup lang="ts">
import { Heart, ShoppingCart, UserRound } from '@lucide/vue';
interface NavItemProps {
    name: string
    routeName: string
}
withDefaults(defineProps<NavItemProps>(), {
    name: '',
    routeName: ''
})
</script>


<template>
    <RouterLink :to="{ name: routeName }" class="flex items-center justify-center bg-secondary text-white! rounded-full size-12 shadow-md">
        <ShoppingCart v-if="name === 'cart'"/>
        <Heart v-else-if="name === 'wishlist'"/>
        <UserRound v-else-if="name === 'profile'"/>
    </RouterLink>
</template>