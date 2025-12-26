<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';
import api from '../services/api';

import { getProvincesArray, getCitiesForProvince } from '../data/indonesianRegions';

const router = useRouter();
const authStore = useAuthStore();
const toast = useToastStore();

const user = ref<any>(null);
const orders = ref<any[]>([]);
const loading = ref(true);
const isEditing = ref(false);

const editForm = ref({
  name: '',
  email: '',
  phone: '',
  address: '',
  province: '',
  city: '',
  postal_code: '',
  profile_image: ''
});

const provinces = ref<string[]>(getProvincesArray());
const citiesForProvince = computed(() => {
  return editForm.value.province ? getCitiesForProvince(editForm.value.province) : [];
});

onMounted(async () => {
  if (!authStore.isAuthenticated) {
    router.push('/login');
    return;
  }

  user.value = authStore.user;
  editForm.value = {
    name: authStore.user?.name || '',
    email: authStore.user?.email || '',
    phone: authStore.user?.phone || '',
    address: authStore.user?.address || '',
    province: authStore.user?.province || '',
    city: authStore.user?.city || '',
    postal_code: authStore.user?.postal_code || '',
    profile_image: authStore.user?.profile_image || ''
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
    user.value = { ...user.value, ...editForm.value };
    authStore.user = { ...authStore.user, ...editForm.value };
    localStorage.setItem('user', JSON.stringify(authStore.user));
    isEditing.value = false;
    toast.show('Profil berhasil diperbarui', 'success');
  } catch (error: any) {
    toast.show(error.response?.data?.message || 'Gagal memperbarui profil', 'error');
  }
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
  <div class="min-h-screen py-6 md:py-12" style="background-color: #F9DFDF;">
    <div class="container mx-auto px-3 sm:px-4">
      <h1 class="text-2xl sm:text-3xl md:text-4xl font-extrabold mb-4 sm:mb-6 md:mb-10 border-b-4 border-emerald-400 inline-block pb-1 sm:pb-2 text-slate-800">Profil Saya</h1>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-3 md:gap-8">
        <!-- Profile Card -->
        <div class="rounded-lg sm:rounded-2xl shadow-sm p-3 sm:p-5 md:p-8 border border-slate-100" style="background-color: #FCF8F8;">
          <div class="text-center mb-3 sm:mb-5 md:mb-8">
            <div class="relative inline-block mb-3 sm:mb-4">
              <div class="w-20 sm:w-24 h-20 sm:h-24 rounded-full bg-gradient-to-br from-emerald-300 to-emerald-500 mx-auto flex items-center justify-center">
                <span class="text-white text-2xl sm:text-3xl font-bold">{{ user?.name?.charAt(0)?.toUpperCase() }}</span>
              </div>
            </div>
            <h2 class="text-lg sm:text-xl md:text-2xl font-bold text-slate-800">{{ user?.name }}</h2>
            <p class="text-xs sm:text-sm md:text-base text-slate-500">{{ user?.email }}</p>
          </div>

          <div class="space-y-1 sm:space-y-2 md:space-y-3 mb-3 sm:mb-5 md:mb-8">
            <div class="flex items-center justify-between p-2 sm:p-3 bg-slate-50 rounded-lg text-xs sm:text-sm md:text-base">
              <span class="text-slate-600">Role:</span>
              <span class="font-bold text-emerald-600">{{ user?.role === 'admin' ? 'Admin' : 'Pelanggan' }}</span>
            </div>
            <div class="flex items-center justify-between p-2 sm:p-3 bg-slate-50 rounded-lg text-xs sm:text-sm">
              <span class="text-slate-600">Anggota:</span>
              <span class="text-slate-700">{{ user?.created_at ? new Date(user.created_at).toLocaleDateString('id-ID') : '-' }}</span>
            </div>
          </div>

          <div v-if="!isEditing" class="space-y-2">
            <button @click="isEditing = true" class="w-full text-white font-bold py-2 text-xs sm:text-sm md:text-base rounded-lg transition-opacity hover:opacity-80" style="background-color: #8BAE66;">
              Edit Profil
            </button>
            <button @click="handleLogout" class="w-full text-white font-bold py-2 text-xs sm:text-sm md:text-base rounded-lg transition-opacity hover:opacity-80" style="background-color: #FF5555;">
              Logout
            </button>
          </div>

          <div v-else class="space-y-2 sm:space-y-2.5 md:space-y-3">
            <div>
              <label class="block text-xs sm:text-sm font-bold text-slate-600 mb-1">Nama Lengkap</label>
              <input v-model="editForm.name" type="text" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 sm:p-2.5 md:p-3 text-xs sm:text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div>
              <label class="block text-xs sm:text-sm font-bold text-slate-600 mb-1">Email</label>
              <input v-model="editForm.email" type="email" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 sm:p-2.5 md:p-3 text-xs sm:text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div>
              <label class="block text-xs sm:text-sm font-bold text-slate-600 mb-1">No. Telepon</label>
              <input v-model="editForm.phone" type="tel" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 sm:p-2.5 md:p-3 text-xs sm:text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div>
              <label class="block text-xs sm:text-sm font-bold text-slate-600 mb-1">Alamat</label>
              <textarea v-model="editForm.address" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 sm:p-2.5 md:p-3 text-xs sm:text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black" rows="2"></textarea>
            </div>
            <div>
              <label class="block text-xs sm:text-sm font-bold text-slate-600 mb-1">Provinsi</label>
              <select v-model="editForm.province" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 sm:p-2.5 md:p-3 text-xs sm:text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
                <option value="">Pilih Provinsi</option>
                <option v-for="prov in provinces" :key="prov" :value="prov">{{ prov }}</option>
              </select>
            </div>
            <div v-if="editForm.province">
              <label class="block text-xs sm:text-sm font-bold text-slate-600 mb-1">Kota/Kabupaten</label>
              <select v-model="editForm.city" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 sm:p-2.5 md:p-3 text-xs sm:text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
                <option value="">Pilih Kota</option>
                <option v-for="city in citiesForProvince" :key="city" :value="city">{{ city }}</option>
              </select>
            </div>
            <div>
              <label class="block text-xs sm:text-sm font-bold text-slate-600 mb-1">Kode Pos</label>
              <input v-model="editForm.postal_code" type="text" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 sm:p-2.5 md:p-3 text-xs sm:text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div class="space-y-2 pt-2 sm:pt-2.5 md:pt-3">
              <button @click="handleUpdateProfile" class="w-full text-white font-bold py-2 text-xs sm:text-sm md:text-base rounded-lg transition-opacity hover:opacity-80" style="background-color: #8BAE66;">
                Simpan Perubahan
              </button>
              <button @click="isEditing = false" class="w-full text-slate-800 font-bold py-2 text-xs sm:text-sm md:text-base rounded-lg transition-colors bg-slate-200 hover:bg-slate-300">
                Batal
              </button>
            </div>
          </div>
        </div>        <!-- Orders Section -->
        <div class="lg:col-span-2">
          <div class="rounded-lg sm:rounded-2xl shadow-sm p-3 sm:p-5 md:p-8 border border-slate-100" style="background-color: #FCF8F8;">
            <h3 class="text-lg sm:text-xl md:text-2xl font-bold text-slate-800 mb-3 sm:mb-4 md:mb-6">Riwayat Pesanan</h3>

            <div v-if="loading" class="text-center py-4 sm:py-6 md:py-8">
              <p class="text-slate-600 animate-pulse text-xs sm:text-sm">Memuat pesanan...</p>
            </div>

            <div v-else-if="orders.length === 0" class="text-center py-4 sm:py-6 md:py-8">
              <p class="text-slate-500 mb-2 sm:mb-3 text-xs sm:text-sm">Belum ada pesanan</p>
              <router-link to="/shop" class="text-emerald-500 hover:text-emerald-600 font-bold text-xs sm:text-sm">Mulai berbelanja →</router-link>
            </div>

            <div v-else class="space-y-2 sm:space-y-3">
              <div v-for="order in orders" :key="order.id" class="border border-slate-200 rounded-lg p-3 md:p-4 hover:shadow-md transition-shadow">
                <div class="flex justify-between items-start gap-2 mb-2">
                  <div class="min-w-0">
                    <p class="text-xs text-slate-500">Nomor Pesanan</p>
                    <p class="text-base md:text-lg font-bold text-slate-800">#{{ order.id }}</p>
                  </div>
                  <div :class="['px-2 md:px-3 py-1 rounded-full text-xs font-bold whitespace-nowrap', getStatusColor(order.status)]">
                    {{ getStatusLabel(order.status) }}
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-2 mb-2 text-xs md:text-sm">
                  <div>
                    <p class="text-slate-500">Tanggal</p>
                    <p class="text-slate-800 font-medium">{{ new Date(order.created_at).toLocaleDateString('id-ID') }}</p>
                  </div>
                  <div class="text-right">
                    <p class="text-slate-500">Total</p>
                    <p class="font-bold text-black">RP {{ order.total_amount?.toLocaleString() }}</p>
                  </div>
                </div>

                <router-link :to="`/invoice/${order.id}`" class="text-emerald-500 hover:text-emerald-600 font-bold text-xs md:text-sm">
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
