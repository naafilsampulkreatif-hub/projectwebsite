<script setup lang="ts">
import { onMounted, computed } from 'vue';
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

// Filter unique items by product id
const uniqueCartItems = computed(() => {
    const seen = new Set();
    return cartStore.items.filter(item => {
        const pid = item.product?.id || item.product_id || item.id;
        if (seen.has(pid)) return false;
        seen.add(pid);
        return true;
    });
});
</script>

<template>
  <div class="min-h-screen text-black" style="background-color: #F9DFDF;">
    <div class="container mx-auto px-3 sm:px-4 py-6 md:py-12">
        <h1 class="text-2xl sm:text-3xl md:text-4xl font-extrabold mb-4 sm:mb-6 md:mb-10 border-b-4 border-blue-500 inline-block pb-1 sm:pb-2">Keranjang Belanja</h1>

        <div v-if="cartStore.loading" class="text-center py-6 sm:py-12 md:py-20 rounded-lg sm:rounded-2xl md:rounded-[2rem] shadow-lg" style="background-color: #FBEFEF;">
            <p class="text-gray-500 mb-3 sm:mb-4 md:mb-6 text-base sm:text-lg">Memuat keranjang...</p>
        </div>

        <div v-else-if="cartStore.items.length === 0" class="text-center py-6 sm:py-12 md:py-20 rounded-lg sm:rounded-2xl md:rounded-[2rem] shadow-lg" style="background-color: #FBEFEF;">
            <p class="text-gray-500 mb-3 sm:mb-4 md:mb-6 text-base sm:text-lg">Keranjang anda kosong.</p>
            <router-link to="/shop" class="text-blue-500 font-bold text-base sm:text-lg hover:underline uppercase tracking-widest">Lanjut Belanja ></router-link>
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6 md:gap-10">
            <!-- Cart Items -->
            <div class="lg:col-span-2 p-4 sm:p-6 md:p-8 rounded-lg sm:rounded-2xl md:rounded-[2rem] shadow-lg" style="background-color: #FBEFEF;">
                <div class="overflow-x-auto">
                <table class="w-full text-left text-xs sm:text-sm">
                    <thead>
                        <tr class="border-b-2 border-gray-100">
                            <th class="pb-3 sm:pb-4 md:pb-6 font-bold text-gray-400 uppercase tracking-wider text-xs md:text-sm">Produk</th>
                            <th class="pb-3 sm:pb-4 md:pb-6 font-bold text-gray-400 uppercase tracking-wider text-xs md:text-sm">Harga</th>
                            <th class="pb-3 sm:pb-4 md:pb-6 font-bold text-gray-400 uppercase tracking-wider text-xs md:text-sm">Jumlah</th>
                            <th class="pb-3 sm:pb-4 md:pb-6 font-bold text-gray-400 uppercase tracking-wider text-xs md:text-sm">Total</th>
                        </tr>
                    </thead>
                                        <tbody class="divide-y divide-gray-100">
                                                <tr v-for="item in uniqueCartItems" :key="item.id">
                                                        <td class="py-3 sm:py-4 md:py-6 flex items-center gap-3 sm:gap-6">
                                                                <div class="w-12 h-12 sm:w-16 sm:h-16 md:w-20 md:h-20 bg-gray-100 rounded-lg sm:rounded-xl overflow-hidden shadow-sm">
                                                                     <img :src="(item?.product?.image_url) || 'https://via.placeholder.com/150'" class="w-full h-full object-cover">
                                                                </div>
                                                                <div>
                                                                    <div class="font-bold text-sm sm:text-base md:text-lg text-gray-800">{{ item?.product?.name || 'Produk' }}</div>
                                                                    <div class="text-xs sm:text-sm text-gray-500">SKU: {{ item.product.id }}</div>
                                                                </div>
                                                        </td>
                                                        <td class="py-3 sm:py-4 md:py-6 font-medium text-black text-xs sm:text-base">RP {{ fmtPrice(item?.product?.price ?? item?.price) }}</td>
                                                        <td class="py-3 sm:py-4 md:py-6 font-medium text-gray-600">
                                                            <div class="inline-flex items-center border rounded-full overflow-hidden text-xs sm:text-sm">
                                                                <button @click="decrease(item)" class="px-2 sm:px-3 py-1 text-white rounded transition-colors hover:opacity-80" style="background-color: #547792;">-</button>
                                                                <div class="px-2 sm:px-4">{{ item.quantity ?? 0 }}</div>
                                                                <button @click="increase(item)" class="px-2 sm:px-3 py-1 text-white rounded transition-colors hover:opacity-80" style="background-color: #547792;">+</button>
                                                            </div>
                                                            <button @click="remove(item)" class="text-xs sm:text-sm text-white mt-1 sm:mt-2 px-2 sm:px-3 py-1 rounded transition-colors hover:opacity-80" style="background-color: #FF5555;">Hapus</button>
                                                        </td>
                                                        <td class="py-3 sm:py-4 md:py-6 font-extrabold text-black text-sm sm:text-base md:text-lg">RP {{ fmtPrice((item?.product?.price ?? item?.price) * (item.quantity ?? 0)) }}</td>
                                                </tr>
                                        </tbody>
                </table>
                </div>
            </div>

            <!-- Summary -->
            <div class="p-4 sm:p-6 md:p-8 rounded-lg sm:rounded-2xl md:rounded-[2rem] shadow-lg h-fit" style="background-color: #FBEFEF;">
                <h2 class="text-lg sm:text-xl md:text-2xl font-bold mb-4 sm:mb-6">Total Belanja</h2>
                <div class="flex justify-between mb-3 sm:mb-4 border-b border-gray-100 pb-3 sm:pb-4 text-xs sm:text-sm">
                    <span class="text-gray-500">Subtotal:</span>
                    <span class="font-bold text-black">RP {{ cartStore.totalPrice.toLocaleString() }}</span>
                </div>
                <div class="flex justify-between mb-6 sm:mb-8 text-sm sm:text-base md:text-xl">
                    <span class="font-bold">Total:</span>
                    <span class="font-extrabold text-black">RP {{ cartStore.totalPrice.toLocaleString() }}</span>
                </div>

                <button @click="proceedToCheckout" class="w-full text-white font-bold py-2 sm:py-3 md:py-4 rounded-full text-sm sm:text-base transition-colors hover:opacity-80" style="background-color: #547792;">
                    Checkout Sekarang
                </button>
            </div>
        </div>
    </div>
  </div>
</template>
