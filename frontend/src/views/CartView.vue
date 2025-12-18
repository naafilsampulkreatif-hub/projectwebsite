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

const decrease = async (item: any) => {
        const newQty = (item.quantity || 0) - 1;
        try {
            await cartStore.updateQuantity(item.product_id || item.product.id, newQty);
        } catch (e) {
            console.error('Failed to decrease quantity', e);
        }
};

const increase = async (item: any) => {
        const newQty = (item.quantity || 0) + 1;
        try {
            await cartStore.updateQuantity(item.product_id || item.product.id, newQty);
        } catch (e) {
            console.error('Failed to increase quantity', e);
        }
};

const remove = async (item: any) => {
        try {
            await cartStore.removeItem(item.id);
        } catch (e) {
            console.error('Failed to remove item', e);
        }
};

const fmtPrice = (p: any) => {
    const n = (p === undefined || p === null) ? 0 : Number(p);
    if (isNaN(n)) return '0';
    return n.toLocaleString();
};
</script>

<template>
  <div class="min-h-screen bg-[#F3F3F3] text-black">
    <div class="container mx-auto px-4 py-12">
        <h1 class="text-4xl font-extrabold mb-10 border-b-4 border-neon-green inline-block pb-2">Keranjang Belanja</h1>

        <div v-if="cartStore.loading" class="text-center py-20 bg-white rounded-[2rem] shadow-lg">
            <p class="text-gray-500 mb-6 text-xl">Memuat keranjang...</p>
        </div>

        <div v-else-if="cartStore.items.length === 0" class="text-center py-20 bg-white rounded-[2rem] shadow-lg">
            <p class="text-gray-500 mb-6 text-xl">Keranjang anda kosong.</p>
            <router-link to="/shop" class="text-neon-green font-bold text-lg hover:underline uppercase tracking-widest">Lanjut Belanja ></router-link>
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
                                                                     <img :src="(item?.product?.image_url) || 'https://via.placeholder.com/150'" class="w-full h-full object-cover">
                                                                </div>
                                                                <div>
                                                                    <div class="font-bold text-lg text-gray-800">{{ item?.product?.name || 'Produk' }}</div>
                                                                    <div class="text-sm text-gray-500">SKU: {{ item.product.id }}</div>
                                                                </div>
                                                        </td>
                                                        <td class="py-6 font-medium text-gray-600">RP {{ fmtPrice(item?.product?.price ?? item?.price) }}</td>
                                                        <td class="py-6 font-medium text-gray-600">
                                                            <div class="inline-flex items-center border rounded-full overflow-hidden">
                                                                <button @click="decrease(item)" class="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-400 transition-colors">-</button>
                                                                <div class="px-4">{{ item.quantity ?? 0 }}</div>
                                                                <button @click="increase(item)" class="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-400 transition-colors">+</button>
                                                            </div>
                                                            <button @click="remove(item)" class="text-sm bg-red-600 text-white mt-2 px-3 py-1 rounded hover:bg-red-400 transition-colors">Hapus</button>
                                                        </td>
                                                        <td class="py-6 font-extrabold text-neon-green text-lg">RP {{ fmtPrice((item?.product?.price ?? item?.price) * (item.quantity ?? 0)) }}</td>
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

                <button @click="proceedToCheckout" class="w-full bg-red-600 text-white font-bold py-3 rounded-full hover:bg-red-400 transition-colors">
                    Checkout Sekarang
                </button>
            </div>
        </div>
    </div>
  </div>
</template>
