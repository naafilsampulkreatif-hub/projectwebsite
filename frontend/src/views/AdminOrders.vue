<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '../services/api';

const orders = ref<any[]>([]);
const loading = ref(true);

const fetchOrders = async () => {
  loading.value = true;
  try {
    const res = await api.get('/admin/orders');
    orders.value = res.data;
  } catch (e) {
    console.error('Failed to fetch admin orders', e);
  } finally {
    loading.value = false;
  }
}

onMounted(fetchOrders);
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-6">Pesanan</h1>

    <div v-if="loading" class="text-center py-12">Loading...</div>

    <div v-else>
      <div v-if="orders.length === 0" class="bg-white p-6 rounded-lg shadow">
        <p class="text-gray-600">Belum ada pesanan.</p>
      </div>

      <div v-else class="space-y-6">
        <div v-for="o in orders" :key="o.id" class="bg-white p-6 rounded-lg shadow">
          <div class="flex justify-between items-start mb-4">
            <div>
              <div class="text-sm text-gray-500">Order #{{ o.id }}</div>
              <div class="font-bold text-lg">RP {{ o.total_amount.toLocaleString() }}</div>
            </div>
            <div class="text-sm text-gray-500">{{ new Date(o.created_at).toLocaleString() }}</div>
          </div>

          <div class="mb-3 text-sm">
            <div v-if="o.user">Pemesan: <span class="font-bold">{{ o.user.name }}</span> ({{ o.user.email }})</div>
            <div v-else-if="o.guest_info">Pengirim: <span class="font-bold">{{ o.guest_info.name || 'Tamu' }}</span></div>
            <div v-else>Session: <span class="font-mono text-xs">{{ o.session_id }}</span></div>
          </div>

          <div>
            <table class="w-full text-left">
              <thead class="text-sm text-gray-500 border-b pb-2">
                <tr><th>Produk</th><th class="text-right">Qty</th><th class="text-right">Harga</th><th class="text-right">Total</th></tr>
              </thead>
              <tbody class="divide-y mt-2">
                <tr v-for="(item, idx) in o.items" :key="idx" class="text-sm">
                  <td class="py-2">{{ item.name }}</td>
                  <td class="py-2 text-right">{{ item.quantity }}</td>
                  <td class="py-2 text-right">RP {{ item.price.toLocaleString() }}</td>
                  <td class="py-2 text-right font-bold">RP {{ item.total.toLocaleString() }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>