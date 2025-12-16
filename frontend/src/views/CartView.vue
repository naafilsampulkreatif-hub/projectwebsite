<script setup lang="ts">
import { onMounted } from 'vue';
import { useCartStore } from '../stores/cart';
import { useRouter } from 'vue-router';
import Header from '../components/Header.vue';

const cartStore = useCartStore();
const router = useRouter();

onMounted(() => {
  cartStore.fetchCart();
});

const proceedToCheckout = () => {
    router.push('/checkout');
};
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <Header />
    <div class="container mx-auto px-4 py-12">
        <h1 class="text-3xl font-bold mb-8">Shopping Cart</h1>

        <div v-if="cartStore.items.length === 0" class="text-center py-12">
            <p class="text-gray-600 mb-4">Your cart is empty.</p>
            <router-link to="/" class="text-primary hover:underline">Continue Shopping</router-link>
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <!-- Cart Items -->
            <div class="lg:col-span-2 bg-white p-6 rounded shadow-sm">
                <table class="w-full text-left">
                    <thead>
                        <tr class="border-b">
                            <th class="pb-4">Product</th>
                            <th class="pb-4">Price</th>
                            <th class="pb-4">Quantity</th>
                            <th class="pb-4">Total</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in cartStore.items" :key="item.id" class="border-b last:border-0">
                            <td class="py-4 flex items-center space-x-4">
                                <img :src="item.product.image_url || 'https://via.placeholder.com/50'" class="w-16 h-16 object-cover rounded">
                                <span>{{ item.product.name }}</span>
                            </td>
                            <td class="py-4">${{ item.product.price }}</td>
                            <td class="py-4">{{ item.quantity }}</td>
                            <td class="py-4 font-bold text-primary">${{ item.product.price * item.quantity }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Summary -->
            <div class="bg-white p-6 rounded shadow-sm h-fit">
                <h2 class="text-xl font-bold mb-4">Cart Total</h2>
                <div class="flex justify-between mb-4 border-b pb-4">
                    <span>Subtotal:</span>
                    <span class="font-bold">${{ cartStore.totalPrice }}</span>
                </div>
                <div class="flex justify-between mb-6 text-lg">
                    <span>Total:</span>
                    <span class="font-bold text-primary">${{ cartStore.totalPrice }}</span>
                </div>

                <button @click="proceedToCheckout" class="w-full bg-primary text-white py-3 rounded hover:bg-orange-600 transition-colors">
                    Proceed to Checkout
                </button>
            </div>
        </div>
    </div>
  </div>
</template>
