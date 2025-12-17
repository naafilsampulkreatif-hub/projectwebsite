import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import CartView from '../views/CartView.vue'
import CheckoutView from '../views/CheckoutView.vue'
import InvoiceView from '../views/InvoiceView.vue'
import AdminView from '../views/AdminView.vue'
import AdminProducts from '../views/AdminProducts.vue'
import AdminOrders from '../views/AdminOrders.vue'
import AdminUsers from '../views/AdminUsers.vue'
import AdminSettings from '../views/AdminSettings.vue'
import OrdersView from '../views/OrdersView.vue'
import ShopView from '../views/ShopView.vue'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/cart',
      name: 'cart',
      component: CartView,
    },
    {
      path: '/checkout',
      name: 'checkout',
      component: CheckoutView,
    },
    {
      path: '/invoice/:id',
      name: 'invoice',
      component: InvoiceView,
    },
    {
      path: '/orders',
      name: 'orders',
      component: OrdersView
    },
    {
      path: '/shop',
      name: 'shop',
      component: ShopView
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminView,
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        { path: '', redirect: { name: 'admin-produk' } },
        { path: 'produk', name: 'admin-produk', component: AdminProducts },
        { path: 'pesanan', name: 'admin-pesanan', component: AdminOrders },
        { path: 'pengguna', name: 'admin-pengguna', component: AdminUsers },
        { path: 'pengaturan', name: 'admin-pengaturan', component: AdminSettings }
      ]
    }
  ]
})

// Navigation Guard
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login');
  } else if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next('/');
  } else {
    next();
  }
});

export default router
