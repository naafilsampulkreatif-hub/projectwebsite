import { defineStore } from 'pinia';
import api from '../services/api';

export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [] as any[], // Cart items
  }),
  getters: {
    totalItems: (state) => state.items.reduce((acc, item) => acc + item.quantity, 0),
    totalPrice: (state) => state.items.reduce((acc, item) => acc + (item.product.price * item.quantity), 0),
  },
  actions: {
    // Fetch cart from backend
    async fetchCart() {
      try {
        const response = await api.get('/cart');
        this.items = response.data || [];
      } catch (error) {
        console.error("Fetch cart failed", error);
      }
    },
    // Add item to cart
    async addToCart(productId: number, quantity: number = 1) {
      try {
        await api.post('/cart', { product_id: productId, quantity });
        await this.fetchCart(); // Refresh cart
      } catch (error) {
        console.error("Add to cart failed", error);
      }
    },
    // Checkout
    async checkout() {
      try {
        await api.post('/checkout', {});
        this.items = []; // Clear local cart
        return true;
      } catch (error) {
        console.error("Checkout failed", error);
        return false;
      }
    }
  }
});
