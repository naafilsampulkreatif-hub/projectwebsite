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
const profilePhotoFile = ref<File | null>(null);
const profilePhotoPreview = ref<string>('');
const editForm = ref({
  name: '',
  email: '',
  phone: '',
  address: '',
  province: '',
  city: '',
  postal_code: ''
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
    postal_code: authStore.user?.postal_code || ''
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
    profilePhotoFile.value = null;
    profilePhotoPreview.value = '';
    toast.show('Profil berhasil diperbarui', 'success');
  } catch (error: any) {
    toast.show(error.response?.data?.message || 'Gagal memperbarui profil', 'error');
  }
};

const handlePhotoSelect = (e: Event) => {
  const input = e.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    if (!file.type.startsWith('image/')) {
      toast.show('File harus berupa gambar', 'error');
      return;
    }
    profilePhotoFile.value = file;
    const reader = new FileReader();
    reader.onload = (event) => {
      profilePhotoPreview.value = event.target?.result as string;
      user.value.profile_image = profilePhotoPreview.value;
    };
    reader.readAsDataURL(file);
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
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 py-6 md:py-12">
    <div class="container mx-auto px-4">
      <h1 class="text-3xl md:text-4xl font-extrabold mb-6 md:mb-10 border-b-4 border-emerald-400 inline-block pb-2 text-slate-800">Profil Saya</h1>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 md:gap-8">
        <!-- Profile Card -->
        <div class="bg-white rounded-2xl shadow-sm p-4 md:p-8 border border-slate-100">
          <div class="text-center mb-4 md:mb-8">
            <div class="relative inline-block mb-4">
              <div v-if="profilePhotoPreview" class="w-24 h-24 rounded-full overflow-hidden mx-auto flex items-center justify-center bg-gray-200">
                <img :src="profilePhotoPreview" alt="Profile" class="w-full h-full object-cover">
              </div>
              <div v-else-if="user?.profile_image" class="w-24 h-24 rounded-full overflow-hidden mx-auto flex items-center justify-center bg-gray-200">
                <img :src="user.profile_image" alt="Profile" class="w-full h-full object-cover">
              </div>
              <div v-else class="w-24 h-24 rounded-full bg-gradient-to-br from-emerald-300 to-emerald-500 mx-auto flex items-center justify-center">
                <span class="text-white text-3xl font-bold">{{ user?.name?.charAt(0)?.toUpperCase() }}</span>
              </div>
              <input v-if="isEditing" type="file" accept="image/*" @change="handlePhotoSelect" class="hidden" id="profilePhotoInput">
              <label v-if="isEditing" for="profilePhotoInput" class="absolute bottom-0 right-0 bg-emerald-500 text-white rounded-full p-2 cursor-pointer hover:bg-emerald-600 transition-colors">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"></path>
                </svg>
              </label>
            </div>
            <h2 class="text-xl md:text-2xl font-bold text-slate-800">{{ user?.name }}</h2>
            <p class="text-sm md:text-base text-slate-500">{{ user?.email }}</p>
          </div>

          <div class="space-y-2 md:space-y-3 mb-4 md:mb-8">
            <div class="flex items-center justify-between p-2 md:p-3 bg-slate-50 rounded-lg text-sm md:text-base">
              <span class="text-slate-600">Role:</span>
              <span class="font-bold text-emerald-600">{{ user?.role === 'admin' ? 'Admin' : 'Pelanggan' }}</span>
            </div>
            <div class="flex items-center justify-between p-2 md:p-3 bg-slate-50 rounded-lg text-sm">
              <span class="text-slate-600">Anggota:</span>
              <span class="text-slate-700">{{ user?.created_at ? new Date(user.created_at).toLocaleDateString('id-ID') : '-' }}</span>
            </div>
          </div>

          <div v-if="!isEditing" class="space-y-2">
            <button @click="isEditing = true" class="w-full text-white font-bold py-2 text-sm md:text-base rounded-lg transition-opacity hover:opacity-80" style="background-color: #8BAE66;">
              Edit Profil
            </button>
            <button @click="handleLogout" class="w-full text-white font-bold py-2 text-sm md:text-base rounded-lg transition-opacity hover:opacity-80" style="background-color: #FF5555;">
              Logout
            </button>
          </div>

          <div v-else class="space-y-2 md:space-y-3">
            <div>
              <label class="block text-xs md:text-sm font-bold text-slate-600 mb-1">Nama Lengkap</label>
              <input v-model="editForm.name" type="text" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 md:p-3 text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div>
              <label class="block text-xs md:text-sm font-bold text-slate-600 mb-1">Email</label>
              <input v-model="editForm.email" type="email" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 md:p-3 text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div>
              <label class="block text-xs md:text-sm font-bold text-slate-600 mb-1">No. Telepon</label>
              <input v-model="editForm.phone" type="tel" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 md:p-3 text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div>
              <label class="block text-xs md:text-sm font-bold text-slate-600 mb-1">Alamat</label>
              <textarea v-model="editForm.address" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 md:p-3 text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black" rows="2"></textarea>
            </div>
            <div>
              <label class="block text-xs md:text-sm font-bold text-slate-600 mb-1">Provinsi</label>
              <input v-model="editForm.province" type="text" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 md:p-3 text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div>
              <label class="block text-xs md:text-sm font-bold text-slate-600 mb-1">Kota</label>
              <input v-model="editForm.city" type="text" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 md:p-3 text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div>
              <label class="block text-xs md:text-sm font-bold text-slate-600 mb-1">Kode Pos</label>
              <input v-model="editForm.postal_code" type="text" class="w-full bg-slate-50 border border-slate-200 rounded-lg p-2 md:p-3 text-sm focus:outline-none focus:border-emerald-400 transition-colors text-black">
            </div>
            <div class="space-y-2 pt-2 md:pt-3">
              <button @click="handleUpdateProfile" class="w-full text-white font-bold py-2 text-sm rounded-lg transition-opacity hover:opacity-80" style="background-color: #8BAE66;">
                Simpan Perubahan
              </button>
              <button @click="isEditing = false" class="w-full text-slate-800 font-bold py-2 text-sm rounded-lg transition-colors bg-slate-200 hover:bg-slate-300">
                Batal
              </button>
            </div>
          </div>
        </div>        <!-- Orders Section -->
        <div class="lg:col-span-2">
          <div class="bg-white rounded-2xl shadow-sm p-4 md:p-8 border border-slate-100">
            <h3 class="text-xl md:text-2xl font-bold text-slate-800 mb-4 md:mb-6">Riwayat Pesanan</h3>

            <div v-if="loading" class="text-center py-6 md:py-8">
              <p class="text-slate-600 animate-pulse text-sm">Memuat pesanan...</p>
            </div>

            <div v-else-if="orders.length === 0" class="text-center py-6 md:py-8">
              <p class="text-slate-500 mb-3 text-sm">Belum ada pesanan</p>
              <router-link to="/shop" class="text-emerald-500 hover:text-emerald-600 font-bold text-sm">Mulai berbelanja →</router-link>
            </div>

            <div v-else class="space-y-3">
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
