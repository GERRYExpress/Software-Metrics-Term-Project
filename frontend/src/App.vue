<script setup lang="ts">
import Loading from 'vue-loading-overlay'
import { onMounted, ref , onUnmounted} from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isLoading = ref(true)
let timer: ReturnType<typeof setTimeout> | undefined

router.beforeEach(() => {
    isLoading.value = true
})

router.afterEach(() => {
    isLoading.value = false
})

router.onError(() => {
    isLoading.value = false
})

function show() {
    clearTimeout(timer)
    timer = setTimeout(() => { isLoading.value = true }, 200)
}

function hide() {
    clearTimeout(timer)
    isLoading.value = false
}

onMounted(async () => {
    await router.isReady()
    isLoading.value = false
})

const removeBefore = router.beforeEach(show)
const removeAfter = router.afterEach(hide)
const removeError = router.onError(hide)

onUnmounted(() => {
    removeBefore()
    removeAfter()
    removeError()
    clearTimeout(timer)
})

onMounted(async () => {
    await router.isReady()
    hide()
})
</script>

<template>
    <Loading :active="isLoading" color="#B7EF00" />
    <RouterView/>
</template>
