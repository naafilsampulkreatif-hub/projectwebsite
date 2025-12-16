<script setup lang="ts">
import { onMounted } from 'vue';
import { useCartStore } from '../stores/cart';
import { useRouter } from 'vue-router';

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
  <div class="min-h-screen bg-[#F3F3F3] text-black">
    <div class="container mx-auto px-4 py-12">
        <h1 class="text-4xl font-extrabold mb-10 border-b-4 border-neon-green inline-block pb-2">Keranjang Belanja</h1>

        <div v-if="cartStore.items.length === 0" class="text-center py-20 bg-white rounded-[2rem] shadow-lg">
            <p class="text-gray-500 mb-6 text-xl">Keranjang anda kosong.</p>
            <router-link to="/" class="text-neon-green font-bold text-lg hover:underline uppercase tracking-widest">Lanjut Belanja ></router-link>
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-10">
            <!-- Cart Items -->
            <div class="lg:col-span-2 bg-white p-8 rounded-[2rem] shadow-lg">
                <table class="w-full text-left">
                    <thead>
                        <tr class="border-b-2 border-gray-100">
                            <th class="pb-6 font-bold text-gray-400 uppercase tracking-wider text-sm">Produk</th>
                            <th class="pb-6 font-bold text-gray-400 uppercase tracking-wider text-sm">Harga</th>
                            <th class="pb-6 font-bold text-gray-400 uppercase tracking-wider text-sm">Jumlah</th>
                            <th class="pb-6 font-bold text-gray-400 uppercase tracking-wider text-sm">Total</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100">
                        <tr v-for="item in cartStore.items" :key="item.id">
                            <td class="py-6 flex items-center space-x-6">
                                <div class="w-20 h-20 bg-gray-100 rounded-xl overflow-hidden shadow-sm">
                                   <img :src="item.product.image_url || 'https://via.placeholder.com/150'" class="w-full h-full object-cover">
                                </div>
                                <span class="font-bold text-lg text-gray-800">{{ item.product.name }}</span>
                            </td>
                            <td class="py-6 font-medium text-gray-600">RP {{ item.product.price.toLocaleString() }}</td>
                            <td class="py-6 font-medium text-gray-600">{{ item.quantity }}</td>
                            <td class="py-6 font-extrabold text-neon-green text-lg">RP {{ (item.product.price * item.quantity).toLocaleString() }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Summary -->
            <div class="bg-white p-8 rounded-[2rem] shadow-lg h-fit">
                <h2 class="text-2xl font-bold mb-6">Total Belanja</h2>
                <div class="flex justify-between mb-4 border-b border-gray-100 pb-4">
                    <span class="text-gray-500">Subtotal:</span>
                    <span class="font-bold text-gray-800">RP {{ cartStore.totalPrice.toLocaleString() }}</span>
                </div>
                <div class="flex justify-between mb-8 text-xl">
                    <span class="font-bold">Total:</span>
                    <span class="font-extrabold text-neon-green">RP {{ cartStore.totalPrice.toLocaleString() }}</span>
                </div>

                <button @click="proceedToCheckout" class="w-full bg-neon-green text-black py-4 rounded-full font-bold uppercase tracking-wider hover:bg-[#00cc00] transition-colors shadow-lg hover:shadow-neon-green/50">
                    Checkout Sekarang
                </button>
            </div>
        </div>
    </div>
  </div>
</template>
