<script setup lang="ts">
import AppFooter from '@/components/AppFooter.vue';
import AppNavbar from '@/components/AppNavbar.vue';

</script>

<template>
    <div class="flex flex-col w-full h-screen">
        <AppNavbar/>
        <RouterView/>
        <AppFooter/>
    </div>
</template>