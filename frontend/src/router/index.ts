import GameService from '@/service/GameService'
import { useAuthStore } from '@/stores/auth'
import { useGameStore } from '@/stores/game'
import BuyerLayoutView from '@/views/BuyerLayoutView.vue'
import CartView from '@/views/CartView.vue'
import ContactView from '@/views/ContactView.vue'
import LandingView from '@/views/LandingView.vue'
import LoginFormView from '@/views/LoginFormView.vue'
import ProductDetailView from '@/views/ProductDetailView.vue'
import ProfileView from '@/views/ProfileView.vue'
import RegisterFormView from '@/views/RegisterFormView.vue'
import SearchResultView from '@/views/SearchResultView.vue'
import SellerDashboardView from '@/views/SellerDashboardView.vue'
import SellerLayoutView from '@/views/SellerLayoutView.vue'
import SellerProductFormView from '@/views/SellerProductFormView.vue'
import SellerProfileView from '@/views/SellerProfileView.vue'
import WishlistView from '@/views/WishlistView.vue'
import nProgress from 'nprogress'
import type { Component } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'

// Leaf pages sharing { path, name, component, meta: { title, requiresAuth } }.
function authedPage(path: string, name: string, component: Component, title: string) {
  return { path, name, component, meta: { title, requiresAuth: true } }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login-view',
      component: LoginFormView,
      meta: { title: 'Login' }
    },
    {
      path: '/register',
      name: 'register-view',
      component: RegisterFormView,
      meta: { title: 'Register' }
    },
    {
      path: '',
      name: 'buyer-layout-view',
      component: BuyerLayoutView,
      children: [
        {
          path: '',
          name: 'landing-page-view',
          component: LandingView,
          meta: { title: 'Home' }
        },
        {
          path: '/search',
          name: 'search-result-view',
          component: SearchResultView,
          meta: { title: 'Search' },
          props: (route) => ({
            page: Number.parseInt(route.query.page ? route.query.page.toString() : '1'),
            category: route.query.category ? route.query.category : ''
          })
        },
        authedPage('/profile', 'profile-view', ProfileView, 'Profile'),
        authedPage('/carts', 'cart-view', CartView, 'Carts'),
        authedPage('/wishlists', 'wishlist-view', WishlistView, 'Wishlists'),
        authedPage('/contacts', 'contact-view', ContactView, 'Contacts'),
        {
          path: '/games/:id',
          name: 'game-detail-view',
          component: ProductDetailView,
          beforeEnter: (route) => {
            const gameStore = useGameStore()
            GameService.getGame(Number(route.params.id))
              .then(res => {
                gameStore.save(res.data)
              })
              .catch(() => {
              })
          },
          props: true
        }
      ]
    },
    {
      path: '/seller',
      name: 'seller-layout-view',
      component: SellerLayoutView,
      meta: { title: 'Seller Center', requiresAuth: true },
      children: [
        {
          path: '',
          redirect: { name: 'seller-dashboard-view' }
        },
        authedPage('dashboard', 'seller-dashboard-view', SellerDashboardView, 'Seller Dashboard'),
        authedPage('products/add', 'seller-product-add-view', SellerProductFormView, 'Add Product'),
        authedPage('products/:id/edit', 'seller-product-edit-view', SellerProductFormView, 'Edit Product'),
        authedPage('profile', 'seller-profile-view', SellerProfileView, 'Seller Profile')
      ]
    }
  ],
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.name === from.name && to.query.page !== from.query.page) return false
    return { top: 0 }
  }
})

router.beforeEach(to => {
  const authStore = useAuthStore()
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return {
      name: 'login-view',
      query: { redirect: to.fullPath }
    }
  }

  const defaultTitle = 'KD5F'
  document.title = to.meta.title ? `KD5F - ${String(to.meta.title)}` : defaultTitle
  nProgress.start()
})

router.afterEach(() => {
  nProgress.done()
})

export default router
