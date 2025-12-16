<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const route = useRoute();
const authStore = useAuthStore();
const orderId = route.params.id;
const order = ref<any>(null);
const items = ref<any[]>([]);
const loading = ref(true);
const error = ref('');

const fetchInvoice = async () => {
    try {
        const headers: any = {};
        if (authStore.token) {
            headers['Authorization'] = `Bearer ${authStore.token}`;
        }
        const sessionId = localStorage.getItem('session_id');
        if (sessionId) {
            headers['X-Session-ID'] = sessionId;
        }

        const res = await fetch(`${import.meta.env.VITE_API_URL}/orders/${orderId}/invoice`, {
            headers
        });

        if (!res.ok) throw new Error('Failed to load invoice');

        const data = await res.json();
        order.value = data.order;
        items.value = data.items;
    } catch (e: any) {
        error.value = e.message;
    } finally {
        loading.value = false;
    }
};

onMounted(() => {
    fetchInvoice();
});
</script>

<template>
    <div class="min-h-screen bg-[#F3F3F3] text-black py-12 px-4">
        <div class="container mx-auto max-w-3xl">
            <div v-if="loading" class="text-center">Loading...</div>
            <div v-else-if="error" class="text-center text-red-500">{{ error }}</div>
            <div v-else class="bg-white p-8 rounded-[2rem] shadow-xl border-t-8 border-neon-green">

                <div class="text-center mb-10">
                    <h1 class="text-3xl font-extrabold uppercase tracking-widest mb-2">Nota Pembelian</h1>
                    <p class="text-gray-500">Night Stalkers Inc.</p>
                </div>

                <div class="flex justify-between items-start mb-8 pb-8 border-b border-gray-100">
                    <div>
                        <p class="text-sm text-gray-500 uppercase font-bold">Order ID</p>
                        <p class="text-xl font-bold">#{{ order.id }}</p>
                    </div>
                    <div class="text-right">
                        <p class="text-sm text-gray-500 uppercase font-bold">Tanggal</p>
                        <p class="font-medium">{{ new Date(order.created_at).toLocaleDateString() }}</p>
                    </div>
                </div>

                <div class="mb-8">
                    <p class="text-sm text-gray-500 uppercase font-bold mb-2">Dikirim Kepada</p>
                    <div class="bg-gray-50 p-4 rounded-xl">
                        <p class="font-bold text-lg">{{ order.guest_info?.name || 'Customer' }}</p>
                        <p>{{ order.guest_info?.address }}</p>
                        <p>{{ order.guest_info?.email }}</p>
                        <p>{{ order.guest_info?.phone }}</p>
                    </div>
                </div>

                <div class="mb-8">
                    <p class="text-sm text-gray-500 uppercase font-bold mb-4">Rincian Produk</p>
                    <table class="w-full">
                        <thead>
                            <tr class="text-left text-sm text-gray-400 border-b border-gray-100">
                                <th class="pb-2">Produk</th>
                                <th class="pb-2 text-center">Qty</th>
                                <th class="pb-2 text-right">Harga</th>
                                <th class="pb-2 text-right">Total</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-50">
                            <tr v-for="(item, index) in items" :key="index">
                                <td class="py-3 font-medium">{{ item.name }}</td>
                                <td class="py-3 text-center">{{ item.quantity }}</td>
                                <td class="py-3 text-right">RP {{ item.price.toLocaleString() }}</td>
                                <td class="py-3 text-right font-bold">RP {{ item.total.toLocaleString() }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="flex justify-end mb-8">
                    <div class="w-full md:w-1/2 space-y-3">
                        <div class="flex justify-between items-center text-xl font-extrabold">
                            <span>Total Tagihan</span>
                            <span class="text-neon-green">RP {{ order.total_amount.toLocaleString() }}</span>
                        </div>
                        <div class="flex justify-between items-center text-sm text-gray-500 bg-gray-100 p-2 rounded">
                            <span>Metode Pembayaran</span>
                            <span class="font-bold text-black">COD (Cash on Delivery)</span>
                        </div>
                        <div class="text-right text-xs text-orange-500 font-bold mt-2">
                            * Menunggu Konfirmasi Admin
                        </div>
                    </div>
                </div>

                <div class="text-center pt-8 border-t border-gray-100">
                    <router-link to="/" class="inline-block bg-black text-white px-8 py-3 rounded-full font-bold hover:bg-neon-green hover:text-black transition-colors">
                        Kembali ke Beranda
                    </router-link>
                </div>

            </div>
        </div>
    </div>
</template>
