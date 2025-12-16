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
    async checkout(guestInfo: any = null) {
      try {
        // Send guest_info if available
        const response = await api.post('/checkout', { guest_info: guestInfo });
        this.items = []; // Clear local cart
        return response.data; // Return full data (including order_id)
      } catch (error) {
        console.error("Checkout failed", error);
        throw error; // Throw to handle in view
      }
    }
  }
});
