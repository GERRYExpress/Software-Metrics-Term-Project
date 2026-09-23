<script setup lang="ts">
import { Search } from '@lucide/vue';
import BaseInput from './BaseInput.vue';
import BaseSelect from './BaseSelect.vue';
import AppNavItem from './AppNavItem.vue';
import type { Category } from '@/types.ts';
import { onMounted, ref } from 'vue';
import CategoryService from '@/service/CategoryService.ts';

const navItems = [
    { name: 'cart', routeName: 'cart-view' },
    { name: 'wishlist', routeName: 'wishlist-view' },
    { name: 'profile', routeName: 'profile-view' }
]

const categories = ref<Category[] | null>()


onMounted(() => {
    CategoryService.getCategories()
        .then(res => {
            categories.value = res.data
        })
})

</script>

<template>
    <nav class="w-5xl mx-auto items-center justify-center flex flex-row gap-x-2.5 my-5">
        <div class="relative flex flex-col items-start justify-center w-1/2">
            <Search class="absolute left-1.5"/>
            <BaseInput class="px-8"/>
            <BaseSelect class="absolute top-0 right-0 rounded-r-lg text-xs" v-if="categories" :options="categories"/>
        </div>
        <AppNavItem v-for="navItem in navItems" :name="navItem.name" :routeName="navItem.routeName" :key="navItem.name"/>
    </nav>
</template>