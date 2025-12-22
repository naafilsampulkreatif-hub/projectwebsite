<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import api from '../services/api'
import { useToastStore } from '../stores/toast'
import { useAuthStore } from '../stores/auth'

const toast = useToastStore()
const authStore = useAuthStore()

// Profile
const adminName = ref('')
const adminEmail = ref('')
const isEditingProfile = ref(false)

// Password change
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const isChangingPassword = ref(false)

// Activity Log
const activityLog = ref<any[]>([])
const userStats = ref<any>(null)

// Customer Info
const customerInfo = ref<any[]>([])
const loadingCustomerInfo = ref(false)
const searchCustomer = ref('')

// Tab switching
const activeTab = ref<'profile' | 'password' | 'activity' | 'stats' | 'customers'>('profile')

onMounted(async () => {
  await loadAdminProfile()
  await loadActivityLog()

  await loadStats()
})

const loadAdminProfile = async () => {
  try {
    const res = await api.get('/admin/profile')
    adminName.value = res.data.name
    adminEmail.value = res.data.email
  } catch (e) {
    console.error('Failed to load profile', e)
    toast.show('Gagal memuat profil admin', 'error')
  }
}

const loadActivityLog = async () => {
  try {
    const res = await api.get('/admin/activity-log?limit=100')
    activityLog.value = res.data || []
  } catch (e) {
    console.error('Failed to load activity log', e)
  }
}

const loadStats = async () => {
  try {
    const res = await api.get('/admin/stats')
    userStats.value = res.data
  } catch (e) {
    console.error('Failed to load stats', e)
  }
}

const loadCustomerInfo = async () => {
  loadingCustomerInfo.value = true
  try {
    const res = await api.get('/admin/customer-info')
    customerInfo.value = res.data || []
  } catch (e) {
    console.error('Failed to load customer info', e)
    toast.show('Gagal memuat data pelanggan', 'error')
  } finally {
    loadingCustomerInfo.value = false
  }
}

const filteredCustomers = computed(() => {
  if (!searchCustomer.value) {
    return customerInfo.value
  }
  const query = searchCustomer.value.toLowerCase()
  return customerInfo.value.filter(c =>
    c.full_name?.toLowerCase().includes(query) ||
    c.email?.toLowerCase().includes(query) ||
    c.phone?.includes(query)
  )
})

const updateProfile = async () => {
  try {
    await api.put('/admin/profile', {
      name: adminName.value,
      email: adminEmail.value
    })
    toast.show('Profil diperbarui', 'success')
    isEditingProfile.value = false
    authStore.user = { ...authStore.user, name: adminName.value, email: adminEmail.value }
    localStorage.setItem('user', JSON.stringify(authStore.user))
  } catch (e) {
    console.error('Failed to update profile', e)
    toast.show('Gagal memperbarui profil', 'error')
  }
}

const changePassword = async () => {
  if (newPassword.value !== confirmPassword.value) {
    toast.show('Password tidak cocok', 'error')
    return
  }

  if (newPassword.value.length < 6) {
    toast.show('Password minimal 6 karakter', 'error')
    return
  }

  try {
    await api.post('/admin/password', {
      old_password: oldPassword.value,
      new_password: newPassword.value
    })
    toast.show('Password berhasil diubah', 'success')
    oldPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
    isChangingPassword.value = false
  } catch (e: any) {
    console.error('Failed to change password', e)
    toast.show(e.response?.data || 'Gagal mengubah password', 'error')
  }
}
</script>

<template>
  <div class="space-y-8">
    <!-- Tab Navigation -->
    <div class="flex gap-2 md:gap-4 border-b overflow-x-auto pb-0">
      <button
        @click="activeTab = 'profile'"
        :class="[
          'px-4 md:px-6 py-3 font-bold uppercase border-b-2 transition whitespace-nowrap text-sm md:text-base',
          activeTab === 'profile'
            ? 'border-emerald-500 text-emerald-600'
            : 'border-transparent text-slate-400 hover:text-slate-600'
        ]"
      >
        Profil
      </button>
      <button
        @click="activeTab = 'password'"
        :class="[
          'px-4 md:px-6 py-3 font-bold uppercase border-b-2 transition whitespace-nowrap text-sm md:text-base',
          activeTab === 'password'
            ? 'border-emerald-500 text-emerald-600'
            : 'border-transparent text-slate-400 hover:text-slate-600'
        ]"
      >
        Password
      </button>
      <button
        @click="activeTab = 'customers'; loadCustomerInfo()"
        :class="[
          'px-4 md:px-6 py-3 font-bold uppercase border-b-2 transition whitespace-nowrap text-sm md:text-base',
          activeTab === 'customers'
            ? 'border-emerald-500 text-emerald-600'
            : 'border-transparent text-slate-400 hover:text-slate-600'
        ]"
      >
        Data Pelanggan
      </button>
      <button
        @click="activeTab = 'stats'"
        :class="[
          'px-4 md:px-6 py-3 font-bold uppercase border-b-2 transition whitespace-nowrap text-sm md:text-base',
          activeTab === 'stats'
            ? 'border-emerald-500 text-emerald-600'
            : 'border-transparent text-slate-400 hover:text-slate-600'
        ]"
      >
        Statistik
      </button>
      <button
        @click="activeTab = 'activity'"
        :class="[
          'px-4 md:px-6 py-3 font-bold uppercase border-b-2 transition whitespace-nowrap text-sm md:text-base',
          activeTab === 'activity'
            ? 'border-emerald-500 text-emerald-600'
            : 'border-transparent text-slate-400 hover:text-slate-600'
        ]"
      >
        Aktivitas
      </button>
    </div>

    <!-- Profile Tab -->
    <div v-if="activeTab === 'profile'" class="max-w-2xl">
      <h2 class="text-2xl font-bold mb-6 text-slate-800">Profil Admin</h2>
      <div class="bg-white p-8 rounded-lg shadow space-y-6 border border-slate-200">
        <div v-if="!isEditingProfile" class="space-y-4">
          <div>
            <label class="text-sm text-slate-500 uppercase font-bold">Nama</label>
            <p class="text-2xl font-bold text-slate-800">{{ adminName }}</p>
          </div>
          <div>
            <label class="text-sm text-slate-500 uppercase font-bold">Email</label>
            <p class="text-2xl font-bold text-slate-800">{{ adminEmail }}</p>
          </div>
          <button
            @click="isEditingProfile = true"
            class="text-white font-bold px-6 py-2 rounded-full hover:opacity-80"
            style="background-color: #8BAE66;"
          >
            Edit Profil
          </button>
        </div>

        <div v-else class="space-y-4">
          <div>
            <label class="block font-bold mb-2">Nama</label>
            <input
              v-model="adminName"
              type="text"
              class="w-full border rounded-lg px-4 py-2"
            />
          </div>
          <div>
            <label class="block font-bold mb-2">Email</label>
            <input
              v-model="adminEmail"
              type="email"
              class="w-full border rounded-lg px-4 py-2"
            />
          </div>
          <div class="flex gap-2">
            <button
              @click="updateProfile"
              class="flex-1 text-white font-bold py-2 rounded hover:opacity-80"
              style="background-color: #8BAE66;"
            >
              Simpan
            </button>
            <button
              @click="isEditingProfile = false"
              class="flex-1 bg-gray-300 text-gray-700 font-bold py-2 rounded hover:bg-gray-400"
            >
              Batal
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Password Tab -->
    <div v-if="activeTab === 'password'" class="max-w-2xl">
      <h2 class="text-2xl font-bold mb-6">Ganti Password</h2>
      <div class="bg-white p-8 rounded-lg shadow">
        <div v-if="!isChangingPassword" class="text-center py-8">
          <button
            @click="isChangingPassword = true"
            class="text-white font-bold px-6 py-3 rounded-full hover:opacity-80"
            style="background-color: #FF5555;"
          >
            Mulai Ganti Password
          </button>
        </div>

        <div v-else class="space-y-4">
          <div>
            <label class="block font-bold mb-2">Password Lama</label>
            <input
              v-model="oldPassword"
              type="password"
              class="w-full border rounded-lg px-4 py-2"
              placeholder="Masukkan password lama"
            />
          </div>
          <div>
            <label class="block font-bold mb-2">Password Baru</label>
            <input
              v-model="newPassword"
              type="password"
              class="w-full border rounded-lg px-4 py-2"
              placeholder="Masukkan password baru"
            />
          </div>
          <div>
            <label class="block font-bold mb-2">Konfirmasi Password</label>
            <input
              v-model="confirmPassword"
              type="password"
              class="w-full border rounded-lg px-4 py-2"
              placeholder="Konfirmasi password"
            />
          </div>
          <div class="flex gap-2">
            <button
              @click="changePassword"
              class="flex-1 text-white font-bold py-2 rounded hover:opacity-80"
              style="background-color: #FF5555;"
            >
              Ubah Password
            </button>
            <button
              @click="isChangingPassword = false"
              class="flex-1 bg-gray-300 text-gray-700 font-bold py-2 rounded hover:bg-gray-400"
            >
              Batal
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Tab -->
    <div v-if="activeTab === 'stats'" class="space-y-6">
      <h2 class="text-2xl font-bold">Statistik Sistem</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6" v-if="userStats">
        <div class="bg-white p-6 rounded-lg shadow text-center">
          <p class="text-gray-500 text-sm uppercase font-bold mb-2">Total Pengguna</p>
          <p class="text-4xl font-bold text-neon-green">{{ userStats.total_users }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow text-center">
          <p class="text-gray-500 text-sm uppercase font-bold mb-2">Admin</p>
          <p class="text-4xl font-bold text-blue-600">{{ userStats.admin_users }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow text-center">
          <p class="text-gray-500 text-sm uppercase font-bold mb-2">Customer</p>
          <p class="text-4xl font-bold text-purple-600">{{ userStats.customer_users }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow text-center">
          <p class="text-gray-500 text-sm uppercase font-bold mb-2">Total Pesanan</p>
          <p class="text-4xl font-bold text-orange-600">{{ userStats.total_orders }}</p>
        </div>
      </div>
    </div>

    <!-- Activity Log Tab -->
    <div v-if="activeTab === 'activity'" class="space-y-4">
      <h2 class="text-2xl font-bold mb-6">Riwayat Aktivitas Admin</h2>
      <div class="space-y-3">
        <div
          v-for="log in activityLog"
          :key="log.id"
          class="bg-white p-4 rounded-lg shadow border-l-4 border-neon-green"
        >
          <div class="flex justify-between items-start mb-2">
            <div class="flex gap-2 items-center">
              <span :class="['px-3 py-1 text-xs font-bold rounded text-white bg-neon-green']">
                {{ log.action }}
              </span>
              <span class="text-gray-600 text-sm">Entity: {{ log.entity_type }}</span>
            </div>
            <span class="text-xs text-gray-400">{{ new Date(log.created_at).toLocaleString() }}</span>
          </div>
          <p class="text-gray-700 text-sm" v-if="log.details">
            {{ JSON.stringify(JSON.parse(log.details), null, 2).substring(0, 100) }}...
          </p>
        </div>
        <div v-if="activityLog.length === 0" class="text-center py-8 text-gray-500">
          Belum ada aktivitas
        </div>
      </div>
    </div>

    <!-- Customers Tab -->
    <div v-if="activeTab === 'customers'">
      <h2 class="text-2xl font-bold mb-6 text-slate-800">Data Pelanggan</h2>
      
      <div class="bg-white p-6 rounded-lg shadow border border-slate-200 mb-6">
        <input 
          v-model="searchCustomer"
          type="text"
          placeholder="Cari berdasarkan nama, email, atau no. HP..."
          class="w-full bg-slate-50 border border-slate-200 rounded-lg p-3 focus:outline-none focus:border-emerald-400"
        />
      </div>

      <div class="bg-white rounded-lg shadow border border-slate-200 overflow-hidden">
        <div v-if="loadingCustomerInfo" class="p-8 text-center">
          <p class="text-slate-600 animate-pulse">Memuat data pelanggan...</p>
        </div>

        <div v-else-if="filteredCustomers.length === 0" class="p-8 text-center">
          <p class="text-slate-500">Tidak ada data pelanggan</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-slate-50 border-b border-slate-200">
              <tr>
                <th class="px-6 py-3 text-left font-bold text-slate-700">Nama Lengkap</th>
                <th class="px-6 py-3 text-left font-bold text-slate-700">Email</th>
                <th class="px-6 py-3 text-left font-bold text-slate-700">No. HP</th>
                <th class="px-6 py-3 text-left font-bold text-slate-700">Alamat</th>
                <th class="px-6 py-3 text-left font-bold text-slate-700">Kota</th>
                <th class="px-6 py-3 text-left font-bold text-slate-700">Tanggal</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200">
              <tr v-for="customer in filteredCustomers" :key="customer.id" class="hover:bg-slate-50 transition-colors">
                <td class="px-6 py-4 font-bold text-slate-800">{{ customer.full_name }}</td>
                <td class="px-6 py-4 text-slate-700">{{ customer.email }}</td>
                <td class="px-6 py-4 text-slate-700">{{ customer.phone }}</td>
                <td class="px-6 py-4 text-slate-700 text-sm">{{ customer.address }}</td>
                <td class="px-6 py-4 text-slate-700">{{ customer.city }}, {{ customer.province }}</td>
                <td class="px-6 py-4 text-xs text-slate-600">{{ new Date(customer.created_at).toLocaleDateString('id-ID') }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="filteredCustomers.length > 0" class="mt-4 text-slate-600 text-sm">
        Total: {{ filteredCustomers.length }} pelanggan
      </div>
    </div>
  </div>
</template>

<style scoped>
input[type='text'],
input[type='email'],
input[type='password'] {
  transition: border-color 0.2s;
}

input[type='text']:focus,
input[type='email']:focus,
input[type='password']:focus {
  outline: none;
  border-color: var(--neon-green, #00ff00);
}
</style>