<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';
import api from '../services/api';

const router = useRouter();
const authStore = useAuthStore();
const toast = useToastStore();

const user = ref<any>(null);
const orders = ref<any[]>([]);
const loading = ref(true);
const isEditing = ref(false);
const editForm = ref({
  name: '',
  email: ''
});

onMounted(async () => {
  if (!authStore.isAuthenticated) {
    router.push('/login');
    return;
  }

  user.value = authStore.user;
  editForm.value = {
    name: authStore.user?.name || '',
    email: authStore.user?.email || ''
  };

  await fetchOrders();
});

const fetchOrders = async () => {
  try {
    const response = await api.get('/orders');
    orders.value = response.data || [];
  } catch (error) {
    console.error('Failed to fetch orders:', error);
    orders.value = [];
  } finally {
    loading.value = false;
  }
};

const handleLogout = () => {
  authStore.logout();
  toast.show('Berhasil logout', 'success');
  router.push('/');
};

const handleUpdateProfile = async () => {
  try {
    await api.put('/profile', editForm.value);
    user.value.name = editForm.value.name;
    user.value.email = editForm.value.email;
    isEditing.value = false;
    toast.show('Profil berhasil diperbarui', 'success');
  } catch (error: any) {
    toast.show(error.response?.data?.message || 'Gagal memperbarui profil', 'error');
  }
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'pending':
      return 'Menunggu Pembayaran';
    case 'paid':
      return 'Sudah Dibayar';
    case 'shipped':
      return 'Dikirim';
    case 'cancelled':
      return 'Dibatalkan';
    default:
      return status;
  }
};

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
      return 'bg-gray-50 text-gray-700';
  }
};
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 py-12">
    <div class="container mx-auto px-4">
      <h1 class="text-4xl font-extrabold mb-10 border-b-4 border-emerald-400 inline-block pb-2 text-slate-800">Profil Saya</h1>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Profile Card -->
        <div class="bg-white rounded-2xl shadow-sm p-8 border border-slate-100">
          <div class="text-center mb-8">
            <div class="w-24 h-24 rounded-full bg-gradient-to-br from-emerald-300 to-emerald-500 mx-auto mb-4 flex items-center justify-center">
              <span class="text-white text-4xl font-bold">{{ user?.name?.charAt(0)?.toUpperCase() }}</span>
            </div>
            <h2 class="text-2xl font-bold text-slate-800">{{ user?.name }}</h2>
            <p class="text-slate-500">{{ user?.email }}</p>
          </div>

          <div class="space-y-3 mb-8">
            <div class="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
              <span class="text-slate-600">Role:</span>
              <span class="font-bold text-emerald-600">{{ user?.role === 'admin' ? 'Admin' : 'Pelanggan' }}</span>
            </div>
            <div class="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
              <span class="text-slate-600">Anggota sejak:</span>
              <span class="text-sm text-slate-700">{{ user?.created_at ? formatDate(user.created_at) : '-' }}</span>
            </div>
          </div>

          <div v-if="!isEditing" class="space-y-3">
            <button @click="isEditing = true" class="w-full bg-emerald-500 hover:bg-emerald-600 text-white font-bold py-2 rounded-lg transition-colors">
              Edit Profil
            </button>
            <button @click="handleLogout" class="w-full bg-red-500 hover:bg-red-600 text-white font-bold py-2 rounded-lg transition-colors">
              Logout
            </button>
          </div>

          <div v-else class="space-y-4">
            <div>
              <label class="block text-sm font-bold text-slate-600 mb-2">Nama Lengkap</label>
              <input v-model="editForm.name" type="text" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-3 focus:outline-none focus:border-emerald-400 transition-colors">
            </div>
            <div>
              <label class="block text-sm font-bold text-slate-600 mb-2">Email</label>
              <input v-model="editForm.email" type="email" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-3 focus:outline-none focus:border-emerald-400 transition-colors">
            </div>
            <div class="space-y-3">
              <button @click="handleUpdateProfile" class="w-full bg-emerald-500 hover:bg-emerald-600 text-white font-bold py-2 rounded-lg transition-colors">
                Simpan Perubahan
              </button>
              <button @click="isEditing = false" class="w-full bg-slate-200 hover:bg-slate-300 text-slate-800 font-bold py-2 rounded-lg transition-colors">
                Batal
              </button>
            </div>
          </div>
        </div>

        <!-- Orders Section -->
        <div class="lg:col-span-2">
          <div class="bg-white rounded-2xl shadow-sm p-8 border border-slate-100">
            <h3 class="text-2xl font-bold text-slate-800 mb-6">Riwayat Pesanan</h3>

            <div v-if="loading" class="text-center py-8">
              <p class="text-slate-600 animate-pulse">Memuat pesanan...</p>
            </div>

            <div v-else-if="orders.length === 0" class="text-center py-8">
              <p class="text-slate-500 mb-4">Belum ada pesanan</p>
              <router-link to="/shop" class="text-emerald-500 hover:text-emerald-600 font-bold">Mulai berbelanja →</router-link>
            </div>

            <div v-else class="space-y-4">
              <div v-for="order in orders" :key="order.id" class="border border-slate-200 rounded-lg p-6 hover:shadow-md transition-shadow">
                <div class="flex justify-between items-start mb-4">
                  <div>
                    <p class="text-sm text-slate-500">Nomor Pesanan</p>
                    <p class="text-lg font-bold text-slate-800">#{{ order.id }}</p>
                  </div>
                  <div :class="['px-3 py-1 rounded-full text-sm font-bold', getStatusColor(order.status)]">
                    {{ getStatusLabel(order.status) }}
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-4 mb-4 text-sm">
                  <div>
                    <p class="text-slate-500">Tanggal</p>
                    <p class="text-slate-800 font-medium">{{ formatDate(order.created_at) }}</p>
                  </div>
                  <div class="text-right">
                    <p class="text-slate-500">Total</p>
                    <p class="text-lg font-bold text-emerald-600">RP {{ order.total_amount?.toLocaleString() }}</p>
                  </div>
                </div>

                <router-link :to="`/invoice/${order.id}`" class="text-emerald-500 hover:text-emerald-600 font-bold text-sm">
                  Lihat Invoice →
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
