<script setup lang="ts">
import AppCard from '@/components/AppCard.vue';
import BaseButton from '@/components/BaseButton.vue';
import ProductMainContainer from '@/components/ProductMainContainer.vue';

</script>


<template>
    <ProductMainContainer>
        <h1 class="w-2xl">Profile</h1>
        <AppCard class="flex flex-row gap-x-4 w-2xl">
            <div class="flex flex-col w-full place-self-start gap-y-2 **:px-5 **:py-2 text-start font-bold">
                <p>Username</p>
                <p>Account Created</p>
                <p>Total Spend</p>
                <p>Game Purchased</p>
            </div>
            <div class="flex flex-col w-full place-self-start gap-y-2 text-start **:bg-tertiary **:px-5 **:py-2 **:rounded-lg">
                <p>Name</p>
                <p>25</p>
                <p>$0.00</p>
                <p>0</p>
            </div>
            <img src="https://picsum.photos/128" class="rounded-lg"/>
        </AppCard>
        <div class="w-2xl flex flex-row gap-x-3">
            <BaseButton type="neutral">Logout</BaseButton>
            <BaseButton type="danger-link">Delete My Account</BaseButton>
        </div>
    </ProductMainContainer>
</template>