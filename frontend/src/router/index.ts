import BuyerLayoutView from '@/views/BuyerLayoutView.vue'
import CartView from '@/views/CartView.vue'
import ContactView from '@/views/ContactView.vue'
import LandingView from '@/views/LandingView.vue'
import LoginFormView from '@/views/LoginFormView.vue'
import ProfileView from '@/views/ProfileView.vue'
import RegisterFormView from '@/views/RegisterFormView.vue'
import SearchResultView from '@/views/SearchResultView.vue'
import WishlistView from '@/views/WishlistView.vue'
import nProgress from 'nprogress'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login ',
      name: 'login-view',
      component: LoginFormView,
    },
    {
      path: '/register',
      name: 'register-view',
      component: RegisterFormView
    },
    {
      path: '',
      name: 'buyer-layout-view',
      component: BuyerLayoutView,
      children: [
        {
          path: '',
          name: 'landing-page-view',
          component: LandingView
        },
        {
          path: '/search',
          name: 'search-result-view',
          component: SearchResultView,
          props: true
        },
        {
          path: '/profile',
          name: 'profile-view',
          component: ProfileView
        },
        {
          path: '/cart',
          name: 'cart-view',
          component: CartView
        },
        {
          path: '/wishlists',
          name: 'wishlist-view',
          component: WishlistView
        },
        {
          path: '/contacts',
          name: 'contact-view',
          component: ContactView
        }
      ]
    }
  ],
})

router.beforeEach(() => {
  nProgress.start()
})

router.afterEach(() => {
  nProgress.done()
})

export default router
