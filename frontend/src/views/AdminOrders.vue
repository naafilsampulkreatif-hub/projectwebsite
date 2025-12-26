<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
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

const printInvoice = () => {
  window.print();
};

const getCustomerName = (order: any) => {
  // If it's a registered user order, get name from user object
  if (order.user?.name) {
    return order.user.name;
  }
  // Otherwise get from guest_info
  if (order.guest_info?.full_name) {
    return order.guest_info.full_name;
  }
  return 'Guest';
};

const getCustomerEmail = (order: any) => {
  // If it's a registered user order, get email from user object
  if (order.user?.email) {
    return order.user.email;
  }
  // Otherwise get from guest_info
  if (order.guest_info?.email) {
    return order.guest_info.email;
  }
  return '-';
};

const getShippingMethodName = (order: any) => {
  if (order.shipping_method?.name) {
    return order.shipping_method.name;
  }
  return 'COD (Bayar di Tempat)';
};

const getShippingCost = (order: any) => {
  if (order.shipping_method?.cost) {
    return order.shipping_method.cost;
  }
  return 0;
};
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold mb-2 text-slate-800">Pesanan</h1>
      <p class="text-slate-600">Kelola semua pesanan pelanggan Anda</p>
    </div>

    <!-- Filter & Search -->
    <div class="rounded-xl shadow-sm border border-slate-200 p-6 space-y-4" style="background-color: #FCF8F8;">
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
    <div class="rounded-xl shadow-sm border border-slate-200 overflow-hidden" style="background-color: #FCF8F8;">
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
                  <p class="font-bold text-slate-800">{{ getCustomerName(order) }}</p>
                  <p class="text-xs text-slate-500">{{ getCustomerEmail(order) }}</p>
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
      <div class="rounded-xl shadow-sm border border-slate-200 p-4" style="background-color: #FCF8F8;">
        <p class="text-slate-600 text-sm font-bold mb-1">Total Pesanan</p>
        <p class="text-2xl font-bold text-slate-800">{{ orders.length }}</p>
      </div>
      <div class="rounded-xl shadow-sm border border-slate-200 p-4" style="background-color: #FCF8F8;">
        <p class="text-slate-600 text-sm font-bold mb-1">Menunggu</p>
        <p class="text-2xl font-bold text-amber-600">{{ orders.filter(o => o.status === 'pending').length }}</p>
      </div>
      <div class="rounded-xl shadow-sm border border-slate-200 p-4" style="background-color: #FCF8F8;">
        <p class="text-slate-600 text-sm font-bold mb-1">Dikirim</p>
        <p class="text-2xl font-bold text-emerald-600">{{ orders.filter(o => o.status === 'shipped').length }}</p>
      </div>
      <div class="rounded-xl shadow-sm border border-slate-200 p-4" style="background-color: #FCF8F8;">
        <p class="text-slate-600 text-sm font-bold mb-1">Total Penjualan</p>
        <p class="text-2xl font-bold text-black">RP {{ orders.reduce((sum, o) => sum + (o.total_amount || 0), 0).toLocaleString() }}</p>
      </div>
    </div>

    <!-- Detail Modal - Invoice View -->
    <div v-if="showDetailModal && selectedOrder" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4 overflow-y-auto">
      <div class="rounded-xl shadow-2xl w-full max-w-4xl my-8" style="background-color: #FCF8F8;">
        <!-- Header -->
        <div class="bg-gradient-to-r from-sky-500 to-sky-600 text-white p-6 flex justify-between items-center sticky top-0">
          <div>
            <h2 class="text-2xl font-bold">Invoice #{{ selectedOrder.id }}</h2>
            <p class="text-sky-100 text-sm mt-1">{{ formatDate(selectedOrder.created_at) }}</p>
          </div>
          <button @click="closeDetail" class="text-white hover:bg-sky-700 p-2 rounded transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <!-- Invoice Content -->
        <div class="p-8 print:bg-white print:p-0" style="background-color: #FBEFEF;">
          <div class="p-8 print:p-0 print:bg-white" style="background-color: #FCF8F8;">
            <div class="max-w-4xl mx-auto">
              <!-- Header -->
              <div class="mb-6 pb-6 border-b border-gray-300">
                <div class="flex justify-between items-center">
                  <div>
                    <h1 class="text-2xl font-bold text-red-600">NIGHT STALKERS</h1>
                    <p class="text-xs text-gray-600 mt-1">PT Sampulkreativ</p>
                  </div>
                  <div class="text-right">
                    <h2 class="text-3xl font-bold text-gray-800">INVOICE</h2>
                    <p class="text-xs text-gray-500 mt-1">#{{ selectedOrder.id }}</p>
                  </div>
                </div>
              </div>

              <!-- Invoice Info -->
              <div class="grid grid-cols-2 gap-6 mb-6">
                <div>
                  <h3 class="font-bold text-xs text-gray-600 mb-2 uppercase">Tagihan Kepada:</h3>
                  <div v-if="selectedOrder.guest_info" class="text-sm space-y-0.5">
                    <p class="font-bold">{{ selectedOrder.guest_info.full_name || '-' }}</p>
                    <p class="text-xs">{{ selectedOrder.guest_info.address || '-' }}</p>
                    <p class="text-xs">{{ selectedOrder.guest_info.city || '-' }}, {{ selectedOrder.guest_info.province || '-' }} {{ selectedOrder.guest_info.postal_code || '-' }}</p>
                    <p class="text-xs text-gray-600">{{ selectedOrder.guest_info.email || '-' }}</p>
                    <p class="text-xs text-gray-600">{{ selectedOrder.guest_info.phone || '-' }}</p>
                  </div>
                  <div v-else class="text-xs text-gray-500">Informasi tidak tersedia</div>
                </div>
                <div class="text-right">
                  <div class="mb-2">
                    <span class="text-xs text-gray-600">TANGGAL:</span>
                    <p class="font-bold text-sm">{{ formatDate(selectedOrder.created_at).split(' ')[0] }}</p>
                  </div>
                </div>
              </div>

              <!-- Items Table -->
              <div class="mb-4">
                <table class="w-full text-xs border-collapse">
                  <thead>
                    <tr class="border-b border-gray-800">
                      <th class="text-left py-1 px-2 font-bold">PRODUK</th>
                      <th class="text-center py-1 px-2 font-bold">QTY</th>
                      <th class="text-right py-1 px-2 font-bold">HARGA</th>
                      <th class="text-right py-1 px-2 font-bold">JUMLAH</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="item in selectedOrder.order_items" :key="item.id" class="border-b border-gray-200">
                      <td class="py-2 px-2">{{ item.product_name }}</td>
                      <td class="text-center py-2 px-2">{{ item.quantity }}</td>
                      <td class="text-right py-2 px-2">Rp {{ item.price.toLocaleString() }}</td>
                      <td class="text-right py-2 px-2 font-bold">Rp {{ (item.price * item.quantity).toLocaleString() }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <!-- Shipping & Payment Info -->
              <div class="mb-4 pb-4 border-b border-gray-300">
                <div class="grid grid-cols-2 gap-4 text-xs">
                  <div>
                    <span class="text-gray-600 font-semibold">METODE PENGIRIMAN:</span>
                    <p class="font-bold text-sm">{{ getShippingMethodName(selectedOrder) }}</p>
                  </div>
                  <div class="text-right">
                    <span class="text-gray-600 font-semibold">BIAYA PENGIRIMAN:</span>
                    <p class="font-bold text-sm">Rp {{ getShippingCost(selectedOrder).toLocaleString() }}</p>
                  </div>
                </div>
              </div>

              <!-- Summary -->
              <div class="flex justify-end mb-4">
                <div class="w-56 text-xs">
                  <div class="flex justify-between py-1 px-2 bg-gray-50 border-b border-gray-300">
                    <span class="font-semibold">SUBTOTAL PRODUK</span>
                    <span class="font-semibold">Rp {{ selectedOrder.total_amount.toLocaleString() }}</span>
                  </div>
                  <div class="flex justify-between py-1 px-2 bg-gray-50 border-b border-gray-300">
                    <span class="font-semibold">BIAYA PENGIRIMAN</span>
                    <span class="font-semibold">Rp {{ getShippingCost(selectedOrder).toLocaleString() }}</span>
                  </div>
                  <div class="flex justify-between py-2 px-2 bg-orange-500 text-white font-bold">
                    <span>TOTAL KESELURUHAN</span>
                    <span>Rp {{ (selectedOrder.total_amount + getShippingCost(selectedOrder)).toLocaleString() }}</span>
                  </div>
                </div>
              </div>

              <!-- Notes -->
              <div class="mb-2 pb-2 border-b border-gray-300">
                <h3 class="font-bold text-xs text-gray-600 mb-1">CATATAN:</h3>
                <p class="text-xs text-gray-700">
                  Silakan transfer ke rekening kami atau bayar di tempat (COD). Terima kasih!
                </p>
              </div>

              <!-- Signature -->
              <div class="flex justify-end">
                <div class="text-center text-xs">
                  <p class="text-gray-600 mb-4">Hormat,</p>
                  <p class="font-bold text-sm">PT Sampulkreativ</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer with Print Button -->
        <div class="print:hidden bg-slate-50 p-4 border-t border-slate-200 flex justify-end gap-2">
          <button @click="closeDetail" class="px-6 py-2 bg-slate-300 hover:bg-slate-400 text-slate-800 font-bold rounded-lg transition-colors">
            Tutup
          </button>
          <button @click="printInvoice" class="px-6 py-2 bg-sky-600 hover:bg-sky-700 text-white font-bold rounded-lg transition-colors">
            Cetak Invoice
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<!-- Print Styles -->
<style>
  @media print {
    body {
      background: white;
      margin: 0;
      padding: 0;
    }
    .container {
      max-width: 100% !important;
      padding: 0 !important;
    }
    * {
      page-break-inside: avoid !important;
      break-inside: avoid !important;
    }
    html, body {
      height: 100%;
      margin: 0;
      padding: 0;
    }
  }
  
  @page {
    size: A4;
    margin: 10mm;
  }
</style>