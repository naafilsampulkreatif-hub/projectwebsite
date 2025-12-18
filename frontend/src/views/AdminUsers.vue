<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../services/api'
import { useToastStore } from '../stores/toast'
import ComingSoon from '../components/ComingSoon.vue'

const toast = useToastStore()

// User management
const users = ref<any[]>([])
const customers = ref<any[]>([])
const selectedUser = ref<any>(null)
const newRole = ref('customer')

// Comment moderation
const comments = ref<any[]>([])
const selectedComment = ref<any>(null)
const moderationStatus = ref('')
const moderationNotes = ref('')

// Tab switching
const activeTab = ref<'users' | 'customers' | 'comments'>('users')

onMounted(async () => {
  await loadUsers()
  await loadCustomers()
  await loadComments()
})

const loadUsers = async () => {
  try {
    const res = await api.get('/admin/users')
    users.value = res.data || []
  } catch (e) {
    console.error('Failed to load users', e)
    toast.show('Gagal memuat data pengguna', 'error')
  }
}

const loadCustomers = async () => {
  try {
    const res = await api.get('/admin/customer-info')
    console.log('Customer info response:', res.data);
    customers.value = Array.isArray(res.data) ? res.data : []
  } catch (e: any) {
    console.error('Failed to load customers', e.response?.data || e.message)
    customers.value = []
  }
}

const loadComments = async () => {
  try {
    const res = await api.get('/admin/comments')
    comments.value = res.data || []
  } catch (e) {
    console.error('Failed to load comments', e)
  }
}

const updateUserRole = async (userId: number, role: string) => {
  try {
    await api.put(`/admin/users/${userId}/role`, { role })
    toast.show('Role pengguna diperbarui', 'success')
    await loadUsers()
  } catch (e) {
    console.error('Failed to update user role', e)
    toast.show('Gagal memperbarui role', 'error')
  }
}

const deleteUser = async (userId: number) => {
  if (!confirm('Yakin ingin menghapus pengguna ini?')) return

  try {
    await api.delete(`/admin/users/${userId}`)
    toast.show('Pengguna dihapus', 'success')
    await loadUsers()
  } catch (e) {
    console.error('Failed to delete user', e)
    toast.show('Gagal menghapus pengguna', 'error')
  }
}

const moderateComment = async (commentId: number, status: string, notes: string) => {
  try {
    await api.post(`/admin/comments/${commentId}/moderate`, {
      status,
      admin_notes: notes
    })
    toast.show(`Komentar ${status === 'approved' ? 'disetujui' : 'ditolak'}`, 'success')
    selectedComment.value = null
    moderationStatus.value = ''
    moderationNotes.value = ''
    await loadComments()
  } catch (e) {
    console.error('Failed to moderate comment', e)
    toast.show('Gagal moderasi komentar', 'error')
  }
}

const deleteComment = async (commentId: number) => {
  if (!confirm('Yakin ingin menghapus komentar ini?')) return

  try {
    await api.delete(`/admin/comments/${commentId}`)
    toast.show('Komentar dihapus', 'success')
    await loadComments()
  } catch (e) {
    console.error('Failed to delete comment', e)
    toast.show('Gagal menghapus komentar', 'error')
  }
}

const selectCommentForModeration = (comment: any) => {
  selectedComment.value = comment
  moderationStatus.value = comment.status
  moderationNotes.value = comment.admin_notes || ''
}
</script>

<template>
  <div class="space-y-8">
    <!-- Tab Navigation -->
    <div class="flex gap-4 border-b">
      <button
        @click="activeTab = 'users'"
        :class="[
          'px-6 py-3 font-bold uppercase border-b-2 transition',
          activeTab === 'users'
            ? 'border-neon-green text-neon-green'
            : 'border-transparent text-gray-400 hover:text-gray-600'
        ]"
      >
        Pengguna
      </button>
      <button
        @click="activeTab = 'customers'"
        :class="[
          'px-6 py-3 font-bold uppercase border-b-2 transition',
          activeTab === 'customers'
            ? 'border-neon-green text-neon-green'
            : 'border-transparent text-gray-400 hover:text-gray-600'
        ]"
      >
        Data Pelanggan
      </button>
      <button
        @click="activeTab = 'comments'"
        :class="[
          'px-6 py-3 font-bold uppercase border-b-2 transition',
          activeTab === 'comments'
            ? 'border-neon-green text-neon-green'
            : 'border-transparent text-gray-400 hover:text-gray-600'
        ]"
      >
        Moderasi Komentar
      </button>
    </div>

    <!-- Users Tab -->
    <div v-if="activeTab === 'users'" class="space-y-4">
      <h2 class="text-2xl font-bold">Daftar Pengguna</h2>
      <div class="overflow-x-auto bg-white rounded-lg shadow">
        <table class="w-full">
          <thead class="bg-gray-100 border-b">
            <tr>
              <th class="px-4 py-3 text-left font-bold">ID</th>
              <th class="px-4 py-3 text-left font-bold">Nama</th>
              <th class="px-4 py-3 text-left font-bold">Email</th>
              <th class="px-4 py-3 text-left font-bold">Role</th>
              <th class="px-4 py-3 text-left font-bold">Terdaftar</th>
              <th class="px-4 py-3 text-left font-bold">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id" class="border-b hover:bg-gray-50">
              <td class="px-4 py-3">{{ user.id }}</td>
              <td class="px-4 py-3">{{ user.name }}</td>
              <td class="px-4 py-3">{{ user.email }}</td>
              <td class="px-4 py-3">
                <select
                  :value="user.role"
                  @change="updateUserRole(user.id, $event.target.value)"
                  class="border rounded px-2 py-1"
                >
                  <option value="customer">Customer</option>
                  <option value="admin">Admin</option>
                </select>
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ new Date(user.created_at).toLocaleDateString() }}</td>
              <td class="px-4 py-3">
                <button
                  @click="deleteUser(user.id)"
                  class="bg-red-600 text-white px-3 py-1 rounded text-sm hover:bg-red-400 transition-colors"
                >
                  Hapus
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="users.length === 0" class="p-4 text-center text-gray-500">Tidak ada pengguna</div>
      </div>
    </div>

    <!-- Customers Tab -->
    <div v-if="activeTab === 'customers'" class="space-y-4">
      <h2 class="text-2xl font-bold">Data Pelanggan yang Checkout</h2>
      <ComingSoon />
    </div>

    <!-- Comments Tab -->
    <div v-if="activeTab === 'comments'" class="space-y-4">
      <h2 class="text-2xl font-bold">Moderasi Komentar</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Comments List -->
        <div class="space-y-3">
          <div
            v-for="comment in comments"
            :key="comment.id"
            @click="selectCommentForModeration(comment)"
            :class="[
              'p-4 rounded-lg cursor-pointer border-2 transition',
              selectedComment?.id === comment.id
                ? 'border-neon-green bg-green-50'
                : 'border-gray-200 bg-white hover:border-neon-green'
            ]"
          >
            <div class="flex justify-between items-start mb-2">
              <div class="font-bold">{{ comment.name }}</div>
              <span
                :class="[
                  'text-xs font-bold px-2 py-1 rounded',
                  comment.status === 'approved'
                    ? 'bg-green-100 text-green-700'
                    : comment.status === 'rejected'
                      ? 'bg-red-100 text-red-700'
                      : 'bg-yellow-100 text-yellow-700'
                ]"
              >
                {{ comment.status }}
              </span>
            </div>
            <div class="text-sm text-gray-600 mb-2">{{ comment.email }}</div>
            <p class="text-sm text-gray-700 line-clamp-2">{{ comment.text }}</p>
          </div>
          <div v-if="comments.length === 0" class="text-center text-gray-500 py-8">Tidak ada komentar</div>
        </div>

        <!-- Moderation Form -->
        <div v-if="selectedComment" class="bg-white p-6 rounded-lg shadow space-y-4">
          <h3 class="text-xl font-bold">Detail Komentar</h3>

          <div class="bg-gray-50 p-4 rounded">
            <p class="font-bold mb-2">Komentar:</p>
            <p class="text-gray-700">{{ selectedComment.text }}</p>
          </div>

          <div class="space-y-2">
            <label class="block font-bold">Status</label>
            <select
              v-model="moderationStatus"
              class="w-full border rounded px-3 py-2"
            >
              <option value="pending">Pending</option>
              <option value="approved">Disetujui</option>
              <option value="rejected">Ditolak</option>
            </select>
          </div>

          <div class="space-y-2">
            <label class="block font-bold">Catatan Admin</label>
            <textarea
              v-model="moderationNotes"
              rows="4"
              placeholder="Tuliskan catatan (opsional)"
              class="w-full border rounded px-3 py-2"
            ></textarea>
          </div>

          <div class="flex gap-2">
            <button
              @click="moderateComment(selectedComment.id, moderationStatus, moderationNotes)"
              class="flex-1 bg-red-600 text-white font-bold py-2 rounded hover:bg-red-400 transition-colors"
            >
              Simpan
            </button>
            <button
              @click="deleteComment(selectedComment.id)"
              class="flex-1 bg-red-600 text-white font-bold py-2 rounded hover:bg-red-400 transition-colors"
            >
              Hapus
            </button>
          </div>
        </div>

        <div v-else class="bg-gray-50 p-6 rounded-lg flex items-center justify-center text-gray-500">
          Pilih komentar untuk dimoderasi
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
}
</style>