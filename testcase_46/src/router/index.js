import { createRouter, createWebHistory } from 'vue-router';
import FormulaInput from '@/components/FormulaInput.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: FormulaInput,
    meta: {
      title: 'Secure Formula Calculator'
    }
  }
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
});

export default router;
