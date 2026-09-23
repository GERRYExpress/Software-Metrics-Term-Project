<script setup lang="ts">
import { Search } from '@lucide/vue';
import BaseInput from './BaseInput.vue';
import AppNavItem from './AppNavItem.vue';

const navItems = [
    { name: 'cart', routeName: 'cart-view' },
    { name: 'wishlist', routeName: 'wishlist-view' },
    { name: 'profile', routeName: 'profile-view' }
]

</script>

<template>
    <nav class="w-full items-center justify-center flex flex-row gap-x-2.5 my-5">
        <div class="relative flex flex-col items-start justify-center w-1/2">
            <Search class="absolute left-1.5"/>
            <BaseInput/>
        </div>
        <AppNavItem v-for="navItem in navItems" :name="navItem.name" :routeName="navItem.routeName" :key="navItem.name"/>
    </nav>
</template>