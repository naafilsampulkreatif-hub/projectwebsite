import { defineStore } from 'pinia';
import api from '../services/api';

export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [] as any[], // Cart items
    loading: false as boolean,
  }),
  getters: {
    totalItems: (state) => state.items.reduce((acc, item) => acc + (item.quantity || 0), 0),
    totalPrice: (state) => state.items.reduce((acc, item) => {
      const price = item?.product?.price ?? item?.price ?? 0;
      const qty = item?.quantity ?? 0;
      return acc + (Number(price) * Number(qty));
    }, 0),
  },
  actions: {
    // Fetch cart from backend
    async fetchCart() {
      this.loading = true;
      try {
        // Clear local optimistic placeholders before fetching authoritative data
        this.items = [];
        const response = await api.get('/cart');
        this.items = response.data || [];
      } catch (error) {
        console.error("Fetch cart failed", error);
        // keep items empty instead of leaving optimistic placeholders
        this.items = [];
      } finally {
        this.loading = false;
      }
    },
    // Add item to cart
    async addToCart(productId: number, quantity: number = 1) {
        // Optimistic update: update local store immediately
        const existing = this.items.find(i => i.product && i.product.id === productId);
        if (existing) {
          existing.quantity += quantity;
        } else {
          // Add a minimal placeholder until server returns full product info
          this.items.push({ id: Date.now(), product: { id: productId, name: 'Loading...', price: 0, image_url: '' }, product_id: productId, productId, quantity });
        }

        try {
          await api.post('/cart', { product_id: productId, quantity });
          await this.fetchCart(); // Refresh cart with authoritative data from server
        } catch (error) {
          console.error("Add to cart failed", error);
          // Revert optimistic update by re-fetching cart (or removing placeholder)
          try { await this.fetchCart(); } catch (e) { console.error('Failed to refresh cart after add error', e); }
          throw error;
        }
    },
    // Update quantity for a product in cart
    async updateQuantity(productId: number, quantity: number) {
      try {
        await api.put('/cart', { product_id: productId, quantity });
        await this.fetchCart();
      } catch (error) {
        console.error('Update cart quantity failed', error);
        throw error;
      }
    },
    // Remove a cart item by cart_items.id
    async removeItem(cartItemId: number) {
      try {
        await api.delete(`/cart/${cartItemId}`);
        await this.fetchCart();
      } catch (error) {
        console.error('Remove cart item failed', error);
        throw error;
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
