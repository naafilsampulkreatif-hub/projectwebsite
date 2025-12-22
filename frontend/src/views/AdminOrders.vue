<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '../services/api';
import { useToastStore } from '../stores/toast';

const toast = useToastStore();
const orders = ref<any[]>([]);
const loading = ref(true);
const selectedStatus = ref('all');
const searchQuery = ref('');
const selectedOrder = ref<any>(null);
const showDetailModal = ref(false);

onMounted(async () => {
  await fetchOrders();
});

const fetchOrders = async () => {
  loading.value = true;
  try {
    const response = await api.get('/admin/orders');
    orders.value = response.data || [];
  } catch (error: any) {
    console.error('Failed to fetch orders:', error);
    toast.show('Gagal memuat pesanan', 'error');
  } finally {
    loading.value = false;
  }
};

const filteredOrders = computed(() => {
  let filtered = orders.value;

  if (selectedStatus.value !== 'all') {
    filtered = filtered.filter(o => o.status === selectedStatus.value);
  }

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(o => 
      o.id.toString().includes(query) ||
      (o.guest_info?.full_name?.toLowerCase().includes(query)) ||
      (o.guest_info?.email?.toLowerCase().includes(query))
    );
  }

  return filtered;
});

const getStatusColor = (status: string) => {
  switch (status) {
    case 'pending':
      return 'bg-amber-50 text-amber-700';
    case 'paid':
      return 'bg-blue-50 text-blue-700';
    case 'shipped':
      return 'bg-emerald-50 text-emerald-700';
    case 'cancelled':
      return 'bg-rose-50 text-rose-700';
    default:
      return 'bg-slate-50 text-slate-700';
  }
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const updateOrderStatus = async (orderId: number, newStatus: string) => {
  try {
    await api.put(`/admin/orders/${orderId}`, { status: newStatus });
    const order = orders.value.find(o => o.id === orderId);
    if (order) {
      order.status = newStatus;
      toast.show('Status pesanan berhasil diperbarui', 'success');
    }
  } catch (error: any) {
    toast.show('Gagal memperbarui status pesanan', 'error');
  }
};

const openDetail = (order: any) => {
  selectedOrder.value = order;
  showDetailModal.value = true;
};

const closeDetail = () => {
  showDetailModal.value = false;
  selectedOrder.value = null;
};

import { computed } from 'vue';
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold mb-2 text-slate-800">Pesanan</h1>
      <p class="text-slate-600">Kelola semua pesanan pelanggan Anda</p>
    </div>

    <!-- Filter & Search -->
    <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6 space-y-4">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-bold text-slate-600 mb-2">Cari Pesanan</label>
          <input 
            v-model="searchQuery"
            type="text"
            placeholder="ID, nama, atau email..."
            class="w-full bg-slate-50 border border-slate-200 rounded-lg p-3 text-sm focus:outline-none focus:border-emerald-400"
          />
        </div>
        <div>
          <label class="block text-sm font-bold text-slate-600 mb-2">Status</label>
          <select 
            v-model="selectedStatus"
            class="w-full bg-slate-50 border border-slate-200 rounded-lg p-3 text-sm focus:outline-none focus:border-emerald-400"
          >
            <option value="all">Semua Status</option>
            <option value="pending">Menunggu Pembayaran</option>
            <option value="paid">Sudah Dibayar</option>
            <option value="shipped">Dikirim</option>
            <option value="cancelled">Dibatalkan</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Orders Table -->
    <div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
      <div v-if="loading" class="p-8 text-center">
        <p class="text-slate-600 animate-pulse">Memuat pesanan...</p>
      </div>

      <div v-else-if="filteredOrders.length === 0" class="p-8 text-center">
        <p class="text-slate-500">Tidak ada pesanan ditemukan</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 border-b border-slate-200">
            <tr>
              <th class="px-6 py-3 text-left font-bold text-slate-700">ID</th>
              <th class="px-6 py-3 text-left font-bold text-slate-700">Pelanggan</th>
              <th class="px-6 py-3 text-left font-bold text-slate-700">Total</th>
              <th class="px-6 py-3 text-left font-bold text-slate-700">Status</th>
              <th class="px-6 py-3 text-left font-bold text-slate-700">Tanggal</th>
              <th class="px-6 py-3 text-left font-bold text-slate-700">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="order in filteredOrders" :key="order.id" class="hover:bg-slate-50 transition-colors">
              <td class="px-6 py-4 font-bold text-slate-800">#{{ order.id }}</td>
              <td class="px-6 py-4">
                <div>
                  <p class="font-bold text-slate-800">{{ order.guest_info?.full_name || 'Guest' }}</p>
                  <p class="text-xs text-slate-500">{{ order.guest_info?.email }}</p>
                </div>
              </td>
              <td class="px-6 py-4 font-bold text-black">RP {{ order.total_amount?.toLocaleString() }}</td>
              <td class="px-6 py-4">
                <select 
                  :value="order.status"
                  @change="updateOrderStatus(order.id, ($event.target as HTMLSelectElement).value)"
                  :class="['px-3 py-1 rounded-full text-xs font-bold border-none focus:outline-none cursor-pointer', getStatusColor(order.status)]"
                >
                  <option value="pending">Menunggu Pembayaran</option>
                  <option value="paid">Sudah Dibayar</option>
                  <option value="shipped">Dikirim</option>
                  <option value="cancelled">Dibatalkan</option>
                </select>
              </td>
              <td class="px-6 py-4 text-slate-600 text-xs">{{ formatDate(order.created_at) }}</td>
              <td class="px-6 py-4">
                <button @click="openDetail(order)" class="text-sky-500 hover:text-sky-600 font-bold text-sm cursor-pointer\">Lihat Detail</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Summary Stats -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-4">
        <p class="text-slate-600 text-sm font-bold mb-1">Total Pesanan</p>
        <p class="text-2xl font-bold text-slate-800">{{ orders.length }}</p>
      </div>
      <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-4">
        <p class="text-slate-600 text-sm font-bold mb-1">Menunggu</p>
        <p class="text-2xl font-bold text-amber-600">{{ orders.filter(o => o.status === 'pending').length }}</p>
      </div>
      <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-4">
        <p class="text-slate-600 text-sm font-bold mb-1">Dikirim</p>
        <p class="text-2xl font-bold text-emerald-600">{{ orders.filter(o => o.status === 'shipped').length }}</p>
      </div>
      <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-4">
        <p class="text-slate-600 text-sm font-bold mb-1">Total Penjualan</p>
        <p class="text-2xl font-bold text-black">RP {{ orders.reduce((sum, o) => sum + (o.total_amount || 0), 0).toLocaleString() }}</p>
      </div>
    </div>

    <!-- Detail Modal -->
    <div v-if="showDetailModal && selectedOrder" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-2xl max-w-2xl w-full max-h-96 overflow-y-auto">
        <!-- Header -->
        <div class="bg-gradient-to-r from-sky-500 to-sky-600 text-white p-6 flex justify-between items-center sticky top-0">
          <div>
            <h2 class="text-2xl font-bold">Detail Pesanan #{{ selectedOrder.id }}</h2>
            <p class="text-sky-100 text-sm mt-1">{{ formatDate(selectedOrder.created_at) }}</p>
          </div>
          <button @click="closeDetail" class="text-white hover:bg-sky-700 p-2 rounded transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <!-- Content -->
        <div class="p-6 space-y-6">
          <!-- Customer Info -->
          <div v-if="selectedOrder.guest_info">
            <h3 class="font-bold text-slate-800 mb-3">Informasi Pelanggan</h3>
            <div class="bg-slate-50 p-4 rounded-lg space-y-2 text-sm">
              <p><span class="font-semibold text-slate-700">Nama:</span> {{ selectedOrder.guest_info.full_name }}</p>
              <p><span class="font-semibold text-slate-700">Email:</span> {{ selectedOrder.guest_info.email }}</p>
              <p><span class="font-semibold text-slate-700">Telepon:</span> {{ selectedOrder.guest_info.phone }}</p>
              <p><span class="font-semibold text-slate-700">Alamat:</span> {{ selectedOrder.guest_info.address }}, {{ selectedOrder.guest_info.city }}, {{ selectedOrder.guest_info.province }}</p>
            </div>
          </div>

          <!-- Order Items -->
          <div>
            <h3 class="font-bold text-slate-800 mb-3">Produk yang Dibeli</h3>
            <div class="space-y-3">
              <div v-for="item in selectedOrder.order_items" :key="item.id" class="border border-slate-200 rounded-lg p-3 flex gap-4">
                <div class="flex-1">
                  <p class="font-semibold text-slate-800">{{ item.product_name }}</p>
                  <p class="text-sm text-slate-600">Kuantitas: {{ item.quantity }}</p>
                  <p class="text-sm text-black">Harga: RP {{ item.price.toLocaleString() }} x {{ item.quantity }}</p>
                </div>
                <div class="font-bold text-black">
                  RP {{ (item.price * item.quantity).toLocaleString() }}
                </div>
              </div>
            </div>
          </div>

          <!-- Summary -->
          <div class="bg-slate-50 p-4 rounded-lg">
            <div class="flex justify-between items-center font-bold text-lg">
              <span class="text-slate-800">Total:</span>
              <span class="text-black">RP {{ selectedOrder.total_amount?.toLocaleString() }}</span>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="bg-slate-50 p-4 border-t border-slate-200 flex justify-end gap-2">
          <button @click="closeDetail" class="px-6 py-2 bg-slate-300 hover:bg-slate-400 text-slate-800 font-bold rounded-lg transition-colors">
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>