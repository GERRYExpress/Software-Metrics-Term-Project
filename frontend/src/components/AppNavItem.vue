<script setup lang="ts">
import { Heart, Home, ShoppingCart, UserRound } from '@lucide/vue';

interface NavItemProps {
    name: string
    routeName: string
    count?: number
}

withDefaults(defineProps<NavItemProps>(), {
    name: '',
    routeName: '',
    count: 0
})
</script>

<template>
    <RouterLink
        :to="{ name: routeName }"
        class="flex items-center justify-center bg-secondary text-white! rounded-full size-12 shadow-md relative transition-all hover:scale-105 hover:shadow-lg">
        <ShoppingCart v-if="name === 'cart'"/>
        <Heart v-else-if="name === 'wishlist'"/>
        <UserRound v-else-if="name === 'profile'"/>
        <Home v-else-if="name === 'home'"/>

        <!-- Counter Badge -->
        <span
            v-if="count > 0"
            class="absolute -top-1 -right-1 bg-action text-tertiary text-xs font-extrabold rounded-full size-5 flex items-center justify-center shadow-md animate-in zoom-in duration-200">
            {{ count > 9 ? '9+' : count }}
        </span>
    </RouterLink>
</template>