import { createRouter, createWebHistory } from 'vue-router';
import UserProfile from '../views/UserProfile.vue';
import { useAuthStore } from '../store/authStore';

const routes = [
  {
    path: '/',
    name: 'Home',
    redirect: '/profile'
  },
  {
    path: '/profile',
    name: 'UserProfile',
    component: UserProfile,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    // In a real app, this would be a Login component
    component: { template: '<div>Login Page</div>' }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Global Navigation Guard for Authentication
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    // SECURITY: Redirect unauthenticated users away from protected routes.
    console.warn(`Access denied for route: ${to.path}. Redirecting to login.`);
    next({ name: 'Login' });
  } else {
    next();
  }
});

export default router;
