import BuyerLayoutView from '@/views/BuyerLayoutView.vue'
import LandingView from '@/views/LandingView.vue'
import LoginFormView from '@/views/LoginFormView.vue'
import RegisterFormView from '@/views/RegisterFormView.vue'
import SearchResultView from '@/views/SearchResultView.vue'
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
        }
      ]
    }
  ],
})

export default router
