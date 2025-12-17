<script setup lang="ts">
import { ref, onMounted } from 'vue'
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

// Tab switching
const activeTab = ref<'profile' | 'password' | 'activity' | 'stats'>('profile')

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

const getActionBadge = (action: string) => {
  const colors: Record<string, string> = {
    'create_product': 'bg-blue-100 text-blue-700',
    'update_product': 'bg-yellow-100 text-yellow-700',
    'delete_product': 'bg-red-100 text-red-700',
    'manage_order': 'bg-purple-100 text-purple-700',
    'moderate_comment': 'bg-green-100 text-green-700',
    'delete_comment': 'bg-red-100 text-red-700'
  }
  return colors[action] || 'bg-gray-100 text-gray-700'
}
</script>

<template>
  <div class="space-y-8">
    <!-- Tab Navigation -->
    <div class="flex gap-4 border-b overflow-x-auto">
      <button
        @click="activeTab = 'profile'"
        :class="[
          'px-6 py-3 font-bold uppercase border-b-2 transition whitespace-nowrap',
          activeTab === 'profile'
            ? 'border-neon-green text-neon-green'
            : 'border-transparent text-gray-400 hover:text-gray-600'
        ]"
      >
        Profil
      </button>
      <button
        @click="activeTab = 'password'"
        :class="[
          'px-6 py-3 font-bold uppercase border-b-2 transition whitespace-nowrap',
          activeTab === 'password'
            ? 'border-neon-green text-neon-green'
            : 'border-transparent text-gray-400 hover:text-gray-600'
        ]"
      >
        Ganti Password
      </button>
      <button
        @click="activeTab = 'stats'"
        :class="[
          'px-6 py-3 font-bold uppercase border-b-2 transition whitespace-nowrap',
          activeTab === 'stats'
            ? 'border-neon-green text-neon-green'
            : 'border-transparent text-gray-400 hover:text-gray-600'
        ]"
      >
        Statistik
      </button>
      <button
        @click="activeTab = 'activity'"
        :class="[
          'px-6 py-3 font-bold uppercase border-b-2 transition whitespace-nowrap',
          activeTab === 'activity'
            ? 'border-neon-green text-neon-green'
            : 'border-transparent text-gray-400 hover:text-gray-600'
        ]"
      >
        Riwayat Aktivitas
      </button>
    </div>

    <!-- Profile Tab -->
    <div v-if="activeTab === 'profile'" class="max-w-2xl">
      <h2 class="text-2xl font-bold mb-6">Profil Admin</h2>
      <div class="bg-white p-8 rounded-lg shadow space-y-6">
        <div v-if="!isEditingProfile" class="space-y-4">
          <div>
            <label class="text-sm text-gray-500 uppercase">Nama</label>
            <p class="text-2xl font-bold">{{ adminName }}</p>
          </div>
          <div>
            <label class="text-sm text-gray-500 uppercase">Email</label>
            <p class="text-2xl font-bold">{{ adminEmail }}</p>
          </div>
          <button
            @click="isEditingProfile = true"
            class="bg-neon-green text-black font-bold px-6 py-2 rounded-full hover:bg-[#00cc00]"
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
              class="flex-1 bg-neon-green text-black font-bold py-2 rounded hover:bg-[#00cc00]"
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
            class="bg-neon-green text-black font-bold px-6 py-3 rounded-full hover:bg-[#00cc00]"
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
              class="flex-1 bg-neon-green text-black font-bold py-2 rounded hover:bg-[#00cc00]"
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