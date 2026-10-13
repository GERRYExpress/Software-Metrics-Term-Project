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
import { createRouter, createWebHistory } from 'vue-router'

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
        {
          path: '/profile',
          name: 'profile-view',
          component: ProfileView,
          meta: { title: 'Profile', requiresAuth: true }
        },
        {
          path: '/carts',
          name: 'cart-view',
          component: CartView,
          meta: { title: 'Carts', requiresAuth: true }
        },
        {
          path: '/wishlists',
          name: 'wishlist-view',
          component: WishlistView,
          meta: { title: 'Wishlists', requiresAuth: true }
        },
        {
          path: '/contacts',
          name: 'contact-view',
          component: ContactView,
          meta: { title: 'Contacts', requiresAuth: true }
        },
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
        {
          path: 'dashboard',
          name: 'seller-dashboard-view',
          component: SellerDashboardView,
          meta: { title: 'Seller Dashboard', requiresAuth: true }
        },
        {
          path: 'products/add',
          name: 'seller-product-add-view',
          component: SellerProductFormView,
          meta: { title: 'Add Product', requiresAuth: true }
        },
        {
          path: 'products/:id/edit',
          name: 'seller-product-edit-view',
          component: SellerProductFormView,
          meta: { title: 'Edit Product', requiresAuth: true }
        },
        {
          path: 'profile',
          name: 'seller-profile-view',
          component: SellerProfileView,
          meta: { title: 'Seller Profile', requiresAuth: true }
        }
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
