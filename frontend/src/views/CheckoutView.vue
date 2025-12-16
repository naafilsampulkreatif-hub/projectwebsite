<script setup lang="ts">
import { useCartStore } from '../stores/cart';
import { useRouter } from 'vue-router';
import Header from '../components/Header.vue';

const cartStore = useCartStore();
const router = useRouter();

const handleCheckout = async () => {
    const success = await cartStore.checkout();
    if (success) {
        alert("Order placed successfully!");
        router.push('/');
    } else {
        alert("Checkout failed.");
    }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <Header />
    <div class="container mx-auto px-4 py-12">
        <h1 class="text-3xl font-bold mb-8">Checkout</h1>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12">
             <!-- Billing Details Form (Mock) -->
             <div class="bg-white p-6 rounded shadow-sm">
                 <h2 class="text-xl font-bold mb-4">Billing Details</h2>
                 <form class="space-y-4">
                     <div class="grid grid-cols-2 gap-4">
                         <input type="text" placeholder="First Name" class="border p-2 rounded w-full">
                         <input type="text" placeholder="Last Name" class="border p-2 rounded w-full">
                     </div>
                     <input type="text" placeholder="Company Name (Optional)" class="border p-2 rounded w-full">
                     <input type="text" placeholder="Country" class="border p-2 rounded w-full">
                     <input type="text" placeholder="Street Address" class="border p-2 rounded w-full">
                     <input type="text" placeholder="City" class="border p-2 rounded w-full">
                     <input type="text" placeholder="Postcode / ZIP" class="border p-2 rounded w-full">
                     <input type="tel" placeholder="Phone" class="border p-2 rounded w-full">
                     <input type="email" placeholder="Email" class="border p-2 rounded w-full">
                 </form>
             </div>

             <!-- Order Review -->
             <div class="bg-white p-6 rounded shadow-sm h-fit">
                 <h2 class="text-xl font-bold mb-4">Your Order</h2>
                 <div class="space-y-4 mb-6">
                     <div v-for="item in cartStore.items" :key="item.id" class="flex justify-between">
                         <span>{{ item.product.name }} x {{ item.quantity }}</span>
                         <span>${{ item.product.price * item.quantity }}</span>
                     </div>
                 </div>
                 <div class="flex justify-between border-t pt-4 text-lg font-bold">
                     <span>Total</span>
                     <span class="text-primary">${{ cartStore.totalPrice }}</span>
                 </div>

                 <div class="mt-6">
                     <button @click="handleCheckout" class="w-full bg-primary text-white py-3 rounded hover:bg-orange-600 transition-colors">
                         Place Order
                     </button>
                 </div>
             </div>
        </div>
    </div>
  </div>
</template>
