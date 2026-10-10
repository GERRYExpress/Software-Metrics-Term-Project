import GameService from '@/service/GameService'
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
          props: true
        },
        {
          path: '/profile',
          name: 'profile-view',
          component: ProfileView,
          meta: { title: 'Profile' }
        },
        {
          path: '/carts',
          name: 'cart-view',
          component: CartView,
          meta: { title: 'Carts' }
        },
        {
          path: '/wishlists',
          name: 'wishlist-view',
          component: WishlistView,
          meta: { title: 'Wishlists' }
        },
        {
          path: '/contacts',
          name: 'contact-view',
          component: ContactView,
          meta: { title: 'Contacts' }
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
    }
  ],
})

router.beforeEach(to => {
  const defaultTitle = 'KD5F'
  document.title = to.meta.title ? `KD5F - ${String(to.meta.title)}` : defaultTitle
  nProgress.start()
})

router.afterEach(() => {
  nProgress.done()
})

export default router
