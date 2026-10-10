<script setup lang="ts">
import { Search } from '@lucide/vue';
import BaseInput from './BaseInput.vue';
import BaseSelect from './BaseSelect.vue';
import AppNavItem from './AppNavItem.vue';
import type { Category } from '@/types.ts';
import { computed, onMounted, ref } from 'vue';
import CategoryService from '@/service/CategoryService.ts';
import { useUserActionStore } from '@/stores/userAction';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter()
const route = useRoute()
const userActionStore = useUserActionStore()

const navItems = [
    { name: 'cart', routeName: 'cart-view' },
    { name: 'wishlist', routeName: 'wishlist-view' },
    { name: 'profile', routeName: 'profile-view' }
]

const categories = ref<Category[] | null>()
const searchKeyword = ref<string>('')
const selectedCategory = ref<Category>()

const cartCount = computed(() => userActionStore.cartIds.length)
const wishlistCount = computed(() => userActionStore.wishlistIds.length)

const getItemCount = (name: string) => {
    if (name === 'cart') return cartCount.value
    if (name === 'wishlist') return wishlistCount.value
    return 0
}

onMounted(() => {
    CategoryService.getCategories()
        .then(res => {
            categories.value = res.data
        })
        .catch(err => {
            console.error('Failed to load categories:', err)
        })
})

const handleSearch = () => {
    router.push({ name: 'search-result-view', query: { ...route.query, keyword: searchKeyword.value } })
}

const handleFilter = () => {
    router.push({ name: 'search-result-view', query: { ...route.query, category: selectedCategory.value?.name, page: 1 } })
}
</script>

<template>
    <nav class="w-5xl mx-auto items-center justify-center flex flex-row gap-x-2.5 my-5">
        <AppNavItem name="home" routeName="landing-page-view"/>
        <div class="relative flex flex-col items-start justify-center w-1/2">
            <Search class="absolute left-1.5 text-white/50"/>
            <BaseInput @input="handleSearch" v-model="searchKeyword" class="px-8" placeholder="Search games..."/>
            <BaseSelect @change="handleFilter" v-model="selectedCategory" class="absolute top-0 right-0 rounded-r-lg text-xs" v-if="categories" :options="categories"/>
        </div>
        <AppNavItem
            v-for="navItem in navItems"
            :key="navItem.routeName"
            :name="navItem.name"
            :routeName="navItem.routeName"
            :count="getItemCount(navItem.name)"
        />
    </nav>
</template>