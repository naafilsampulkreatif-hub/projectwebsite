<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '../services/api';

const orders = ref<any[]>([]);
const loading = ref(true);

onMounted(async () => {
    try {
        const res = await api.get('/orders');
        orders.value = res.data;
    } catch (e) {
        console.error("Failed to load orders", e);
    } finally {
        loading.value = false;
    }
});
</script>

<template>
  <div class="min-h-screen text-black" style="background-color: #F9DFDF;">
    <div class="container mx-auto px-4 py-6 md:py-12">
        <h1 class="text-3xl md:text-4xl font-extrabold mb-6 md:mb-10 border-b-4 border-blue-500 inline-block pb-2">Riwayat Pesanan</h1>

        <div v-if="loading" class="text-center py-8">
            <p class="text-sm">Loading...</p>
        </div>

        <div v-else-if="orders.length === 0" class="text-center py-12 md:py-20 rounded-[2rem] shadow-lg" style="background-color: #FBEFEF;">
            <p class="text-gray-500 mb-4 md:mb-6 text-lg">Belum ada pesanan.</p>
            <router-link to="/" class="text-blue-500 font-bold text-lg hover:underline uppercase tracking-widest">Mulai Belanja ></router-link>
        </div>

        <div v-else class="space-y-6">
            <div v-for="order in orders" :key="order.id" class="p-8 rounded-[2rem] shadow-lg" style="background-color: #FBEFEF;\">
                <div class="flex justify-between items-center mb-4">
                    <span class="font-bold text-lg text-gray-800">Order #{{ order.id }}</span>
                    <span class="text-sm text-gray-500">{{ new Date(order.created_at).toLocaleDateString() }}</span>
                </div>
                <div class="flex justify-between items-center">
                    <span class="inline-block px-4 py-1 rounded-full text-sm font-bold bg-gray-100 text-gray-600 uppercase">{{ order.status }}</span>
                    <span class="font-extrabold text-black text-xl">RP {{ order.total_amount.toLocaleString() }}</span>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>
