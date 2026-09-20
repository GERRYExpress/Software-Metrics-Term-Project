import LandingView from '@/views/LandingView.vue'
import LoginFormView from '@/views/LoginFormView.vue'
import RegisterFormView from '@/views/RegisterFormView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
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
      name: 'landing-page-view',
      component: LandingView
    }
  ],
})

export default router
