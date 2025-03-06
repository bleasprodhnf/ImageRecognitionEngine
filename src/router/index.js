import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '../layout/MainLayout.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/500',
      name: 'ServerError',
      component: () => import('../views/error/500.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('../views/error/404.vue')
    },
    {
      path: '/admin',
      component: MainLayout,
      children: [
        {
          path: 'dashboard',
          component: () => import('../views/admin/dashboard.vue')
        },
        {
          path: 'system/params',
          component: () => import('../views/system/params.vue')
        },
        {
          path: 'system/monitor',
          component: () => import('../views/system/monitor.vue')
        },
        {
          path: 'system/logs',
          component: () => import('../views/system/logs.vue')
        },
        {
          path: 'auth/admins',
          component: () => import('../views/auth/admins.vue')
        },
        {
          path: 'auth/roles',
          component: () => import('../views/auth/roles.vue')
        },
        {
          path: 'auth/permissions',
          component: () => import('../views/auth/permissions.vue')
        },
        {
          path: 'model/versions',
          component: () => import('../views/model/versions.vue')
        },
        {
          path: 'model/params',
          component: () => import('../views/model/params.vue')
        },
        {
          path: 'model/monitor',
          component: () => import('../views/model/monitor.vue')
        },
        {
          path: 'stats/system',
          component: () => import('../views/stats/system.vue')
        },
        {
          path: 'stats/accuracy',
          component: () => import('../views/stats/accuracy.vue')
        },
        {
          path: 'stats/resources',
          component: () => import('../views/stats/resources.vue')
        },
        {
          path: 'customer/accounts',
          component: () => import('../views/customer/accounts.vue')
        },
        {
          path: 'customer/packages',
          component: () => import('../views/customer/packages.vue')
        },
        {
          path: 'customer/service',
          component: () => import('../views/customer/service.vue')
        },
      ]
    },
    {
      path: '/admin/login',
      name: 'AdminLogin',
      component: () => import('../views/login/index.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/client',
      component: () => import('../layout/ClientLayout.vue'),
      children: [
        {
          path: 'dashboard',
          component: () => import('../views/client/dashboard.vue')
        },
        {
          path: 'profile',
          component: () => import('../views/client/profile.vue')
        },
        {
          path: 'usage',
          component: () => import('../views/client/usage.vue')
        },
        {
          path: 'api-settings',
          component: () => import('../views/client/api-settings.vue')
        },
        {
          path: 'image-labels',
          component: () => import('../views/client/image-labels.vue')
        },
        {
          path: 'api-documentation',
          component: () => import('../views/client/api-documentation.vue')
        }
      ]
    },
    {
      path: '/client/login',
      name: 'ClientLogin',
      component: () => import('../views/client/login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/home',
      name: 'Home',
      component: () => import('../views/home.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      redirect: '/home'
    }
  ]
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const clientToken = localStorage.getItem('client-token')

  if (to.meta.requiresAuth !== false) {
    if (to.path.startsWith('/admin') && !token) {
      next('/admin/login')
    } else if (to.path.startsWith('/client') && !clientToken) {
      next('/client/login')
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router