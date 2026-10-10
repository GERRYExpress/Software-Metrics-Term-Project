<script setup lang="ts">
import { ref } from 'vue';
import BaseButton from './BaseButton.vue';
import { QrCode, Smartphone, X } from '@lucide/vue';

defineProps<{
    show: boolean,
    totalAmount: number
}>()

const emit = defineEmits<{
    (e: 'close'): void,
    (e: 'proceed'): void
}>()

const selectedMethod = ref<'qr' | 'truemoney'>('qr')
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-black/75 backdrop-blur-xs flex items-center justify-center z-50 p-4">
        <div class="bg-secondary rounded-2xl p-6 max-w-md w-full shadow-2xl flex flex-col items-center gap-y-4 text-center border border-white/10 relative animate-in fade-in zoom-in duration-200">
            <button @click="emit('close')" class="absolute top-4 right-4 text-white/70 hover:text-white cursor-pointer">
                <X class="size-6"/>
            </button>

            <h2 class="text-white text-2xl font-bold">
                {{ selectedMethod === 'qr' ? 'Scan QR to Pay' : 'Truemoney Payment' }}
            </h2>

            <!-- Payment Method Tabs -->
            <div class="flex flex-row gap-x-2 w-full bg-tertiary p-1 rounded-xl">
                <button
                    @click="selectedMethod = 'qr'"
                    type="button"
                    class="flex items-center justify-center gap-x-2 flex-1 py-2 rounded-lg text-sm font-bold transition-all cursor-pointer"
                    :class="selectedMethod === 'qr' ? 'bg-action text-tertiary' : 'text-white/70 hover:text-white'">
                    <QrCode class="size-4"/>
                    QR PromptPay
                </button>
                <button
                    @click="selectedMethod = 'truemoney'"
                    type="button"
                    class="flex items-center justify-center gap-x-2 flex-1 py-2 rounded-lg text-sm font-bold transition-all cursor-pointer"
                    :class="selectedMethod === 'truemoney' ? 'bg-action text-tertiary' : 'text-white/70 hover:text-white'">
                    <Smartphone class="size-4"/>
                    Truemoney
                </button>
            </div>

            <!-- Payment Method Visual -->
            <div v-if="selectedMethod === 'qr'" class="bg-white p-4 rounded-xl flex flex-col items-center gap-y-2">
                <img
                    src="https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=PROMPTPAY_KD5F_DEMO"
                    alt="PromptPay QR Code"
                    class="size-44 object-contain rounded-md"
                />
                <p class="text-xs text-tertiary font-bold">Scan via any Mobile Banking App</p>
            </div>

            <div v-else class="bg-tertiary p-6 rounded-xl flex flex-col items-center gap-y-2 w-full">
                <p class="text-white/80 text-sm">Truemoney Wallet</p>
                <p class="text-action text-lg font-mono font-bold tracking-widest">08X-XXX-XXXX</p>
                <p class="text-xs text-white/60">KaiDao-5-Fong Official Wallet</p>
            </div>

            <!-- Amount -->
            <div class="bg-tertiary/60 px-4 py-2 rounded-lg w-full">
                <p class="text-white/80 text-sm">Total Amount</p>
                <p class="text-action text-2xl font-bold">${{ totalAmount.toFixed(2) }}</p>
            </div>

            <!-- Actions -->
            <div class="flex flex-row gap-x-3 w-full mt-2">
                <BaseButton @click="emit('close')" type="neutral" class="w-1/2">
                    Cancel
                </BaseButton>
                <BaseButton @click="emit('proceed')" class="w-1/2">
                    Proceed
                </BaseButton>
            </div>
        </div>
    </div>
</template>
