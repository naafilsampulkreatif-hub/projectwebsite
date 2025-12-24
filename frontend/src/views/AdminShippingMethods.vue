<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../services/api'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()

interface ShippingMethod {
  id: number
  name: string
  description: string
  cost: number
  is_active: boolean
  created_at?: string
}

const shippingMethods = ref<ShippingMethod[]>([])
const loading = ref(false)
const showForm = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  description: '',
  cost: 0,
  is_active: true
})

onMounted(() => {
  loadShippingMethods()
})

const loadShippingMethods = async () => {
  loading.value = true
  try {
    const res = await api.get('/admin/shipping-methods')
    shippingMethods.value = res.data || []
  } catch (e) {
    console.error('Failed to load shipping methods', e)
    toast.show('Gagal memuat metode pengiriman', 'error')
  } finally {
    loading.value = false
  }
}

const openForm = () => {
  isEditing.value = false
  editingId.value = null
  form.value = {
    name: '',
    description: '',
    cost: 0,
    is_active: true
  }
  showForm.value = true
}

const editMethod = (method: ShippingMethod) => {
  isEditing.value = true
  editingId.value = method.id
  form.value = {
    name: method.name,
    description: method.description,
    cost: method.cost,
    is_active: method.is_active
  }
  showForm.value = true
}

const saveMethod = async () => {
  if (!form.value.name.trim()) {
    toast.show('Nama metode pengiriman harus diisi', 'error')
    return
  }
  if (form.value.cost < 0) {
    toast.show('Biaya pengiriman tidak boleh negatif', 'error')
    return
  }

  try {
    if (isEditing.value && editingId.value) {
      // Update
      await api.put(`/admin/shipping-methods/${editingId.value}`, form.value)
      toast.show('Metode pengiriman berhasil diperbarui', 'success')
    } else {
      // Create
      await api.post('/admin/shipping-methods', form.value)
      toast.show('Metode pengiriman berhasil ditambahkan', 'success')
    }
    showForm.value = false
    await loadShippingMethods()
  } catch (e: any) {
    console.error('Failed to save shipping method', e.response?.data || e.message)
    const errorMsg = e.response?.data?.error || 'Gagal menyimpan metode pengiriman'
    toast.show(errorMsg, 'error')
  }
}

const deleteMethod = async (id: number, name: string) => {
  if (!confirm(`Apakah Anda yakin ingin menghapus "${name}"?`)) {
    return
  }

  try {
    await api.delete(`/admin/shipping-methods/${id}`)
    toast.show('Metode pengiriman berhasil dihapus', 'success')
    await loadShippingMethods()
  } catch (e) {
    console.error('Failed to delete shipping method', e)
    if ((e as any).response?.status === 400) {
      toast.show('Tidak dapat menghapus metode pengiriman default (COD)', 'error')
    } else {
      toast.show('Gagal menghapus metode pengiriman', 'error')
    }
  }
}

const cancelForm = () => {
  showForm.value = false
}

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}
</script>

<template>
  <div class="p-6 max-w-5xl mx-auto">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-800">Metode Pengiriman</h1>
      <button
        @click="openForm"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition"
      >
        + Tambah Metode
      </button>
    </div>

    <!-- Form Modal -->
    <div v-if="showForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">
          {{ isEditing ? 'Edit Metode Pengiriman' : 'Tambah Metode Pengiriman Baru' }}
        </h2>

        <div class="space-y-4">
          <div>
            <label class="block text-gray-700 font-semibold mb-2">Nama Metode</label>
            <input
              v-model="form.name"
              type="text"
              placeholder="Contoh: COD, JNE, Tiki"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label class="block text-gray-700 font-semibold mb-2">Deskripsi</label>
            <textarea
              v-model="form.description"
              placeholder="Deskripsi metode pengiriman..."
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 h-20 resize-none"
            ></textarea>
          </div>

          <div>
            <label class="block text-gray-700 font-semibold mb-2">Biaya Pengiriman (Rp)</label>
            <input
              v-model.number="form.cost"
              type="number"
              placeholder="Contoh: 50000"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div class="flex items-center">
            <input
              v-model="form.is_active"
              type="checkbox"
              id="is_active"
              class="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500 cursor-pointer"
            />
            <label for="is_active" class="ml-2 text-gray-700 cursor-pointer">Aktif</label>
          </div>
        </div>

        <div class="flex gap-3 mt-6">
          <button
            @click="saveMethod"
            class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 rounded-lg transition"
          >
            {{ isEditing ? 'Perbarui' : 'Tambah' }}
          </button>
          <button
            @click="cancelForm"
            class="flex-1 bg-gray-400 hover:bg-gray-500 text-white font-semibold py-2 rounded-lg transition"
          >
            Batal
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8">
      <p class="text-gray-600">Memuat data...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="shippingMethods.length === 0" class="text-center py-8 bg-gray-100 rounded-lg">
      <p class="text-gray-600">Tidak ada metode pengiriman</p>
    </div>

    <!-- Table -->
    <div v-else class="bg-white rounded-lg shadow overflow-x-auto">
      <table class="w-full">
        <thead class="bg-gray-100 border-b">
          <tr>
            <th class="px-6 py-3 text-left text-sm font-semibold text-gray-700">Nama</th>
            <th class="px-6 py-3 text-left text-sm font-semibold text-gray-700">Deskripsi</th>
            <th class="px-6 py-3 text-left text-sm font-semibold text-gray-700">Biaya</th>
            <th class="px-6 py-3 text-left text-sm font-semibold text-gray-700">Status</th>
            <th class="px-6 py-3 text-left text-sm font-semibold text-gray-700">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="method in shippingMethods"
            :key="method.id"
            class="border-b hover:bg-gray-50"
          >
            <td class="px-6 py-4 font-semibold text-gray-800">{{ method.name }}</td>
            <td class="px-6 py-4 text-gray-600">{{ method.description }}</td>
            <td class="px-6 py-4 font-semibold text-gray-800">{{ formatCurrency(method.cost) }}</td>
            <td class="px-6 py-4">
              <span
                :class="method.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                class="px-3 py-1 rounded-full text-sm font-semibold"
              >
                {{ method.is_active ? 'Aktif' : 'Nonaktif' }}
              </span>
            </td>
            <td class="px-6 py-4 space-x-2">
              <button
                @click="editMethod(method)"
                class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded-lg text-sm transition"
              >
                Edit
              </button>
              <button
                @click="deleteMethod(method.id, method.name)"
                class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg text-sm transition"
              >
                Hapus
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
/* Smooth transitions */
button {
  transition: all 0.3s ease;
}
</style>
