<script setup lang="ts">
import { Copyright, InfoIcon } from '@lucide/vue';

</script>


<template>
    <footer class="flex flex-row justify-between items-center px-12 py-16 h-96 w-full bg-tertiary">
        <div class="flex gap-x-1.5">
            <Copyright/><p>2027 KaiDao-5-Fong</p>
        </div>
        <RouterLink class="flex gap-x-1.5" :to="{ name: 'contact-view' }">
            <InfoIcon/>CONTACT US
        </RouterLink>
    </footer>
</template>