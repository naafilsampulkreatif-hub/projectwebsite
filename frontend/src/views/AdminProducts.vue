<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '../services/api';

const products = ref<any[]>([]);
const isEditing = ref(false);
const showForm = ref(false);

const form = ref({
    id: 0,
    name: '',
    description: '',
    price: 0,
    stock: 0,
    image_url: '',
    category_id: 1
});

const generateSlug = (name: string) => {
    return name
        .toLowerCase()
        .trim()
        .replace(/[^\w\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .slice(0, 100);
};

const fetchProducts = async () => {
    try {
        const res = await api.get('/products');
        products.value = res.data;
    } catch (e) {
        console.error("Failed to fetch products", e);
    }
};

onMounted(fetchProducts);

const resetForm = () => {
    form.value = {
        id: 0,
        name: '',
        description: '',
        price: 0,
        stock: 0,
        image_url: '',
        category_id: 1
    };
    isEditing.value = false;
    showForm.value = false;
};

const openAdd = () => {
    resetForm();
    showForm.value = true;
};

const openEdit = (product: any) => {
    form.value = { ...product };
    isEditing.value = true;
    showForm.value = true;
};

const saveProduct = async () => {
    try {
        const productData = {
             ...form.value,
             price: Number(form.value.price),
             slug: generateSlug(form.value.name),
             stock: Number(form.value.stock),
             category_id: Number(form.value.category_id)
        };
        let res;
        if (isEditing.value) {
            res = await api.put(`/admin/products/${form.value.id}`, productData);
        } else {
            res = await api.post('/admin/products', productData);
        }

        // If backend returned created product, update list
        if (res && (res.status === 200 || res.status === 201)) {
            await fetchProducts();
            resetForm();
        } else {
            console.error('Unexpected response saving product', res);
            alert('Error saving product: unexpected response')
        }
    } catch (e: any) {
        console.error('Error saving product:', e);
        const msg = (e as any)?.response?.data || (e as any)?.message || 'Unknown error';
        alert(`Error saving product: ${typeof msg === 'string' ? msg : JSON.stringify(msg)}`);
    }
};

const deleteProduct = async (id: number) => {
    if (!confirm("Apakah anda yakin ingin menghapus produk ini?")) return;
    try {
        await api.delete(`/admin/products/${id}`);
        await fetchProducts();
    } catch (e) {
        alert("Error deleting product");
    }
};

const handleFileUpload = async (event: any) => {
    const file = event.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);

    try {
        const res = await api.post('/admin/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
        form.value.image_url = res.data.url;
    } catch (e) {
        console.error(e);
        alert("Upload gagal");
    }
};
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-10">
        <h1 class="text-3xl font-bold">Dashboard Produk</h1>
        <button @click="openAdd" class="bg-red-600 text-white px-6 py-3 rounded-full font-bold uppercase shadow-lg hover:bg-red-400 transition-all hover:-translate-y-1">
            + Tambah Produk
        </button>
    </div>

    <!-- Product Form -->
    <div v-if="showForm" class="bg-white p-8 rounded-[2rem] shadow-xl mb-10 border border-gray-100">
        <h2 class="text-2xl font-bold mb-6 border-b pb-2">{{ isEditing ? 'Edit' : 'Tambah' }} Produk</h2>
        <form @submit.prevent="saveProduct" class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <input v-model="form.name" placeholder="Nama Produk" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green" required>
            <div class="bg-gray-50 border border-gray-200 p-4 rounded-xl text-gray-500 flex items-center">
                <span class="text-sm">Slug: <strong>{{ generateSlug(form.name) || 'product-name' }}</strong></span>
            </div>
            <input v-model="form.price" type="number" step="100" placeholder="Harga (IDR)" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green" required>
            <input v-model="form.stock" type="number" placeholder="Stok" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green" required>

            <div class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus-within:border-neon-green">
                <label class="block text-sm text-gray-500 mb-1">Gambar Produk</label>
                <input type="file" @change="handleFileUpload" class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-neon-green file:text-black hover:file:bg-[#00cc00]"/>
                <input v-model="form.image_url" type="hidden">
                <div v-if="form.image_url" class="mt-2">
                    <img :src="form.image_url" class="h-20 w-20 object-cover rounded">
                </div>
            </div>

            <input v-model="form.category_id" type="number" placeholder="ID Kategori" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green">
            <textarea v-model="form.description" placeholder="Deskripsi" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green md:col-span-2 h-32"></textarea>

            <div class="md:col-span-2 flex space-x-4 justify-end">
                <button type="button" @click="resetForm" class="px-6 py-3 rounded-full font-bold text-gray-500 hover:bg-gray-100">Batal</button>
                <button type="submit" class="bg-red-600 text-white px-8 py-3 rounded-full font-bold shadow-lg hover:bg-red-400">Simpan</button>
            </div>
        </form>
    </div>

    <!-- Product List -->
    <div class="bg-white rounded-[2rem] shadow-xl overflow-hidden border border-gray-100">
        <table class="w-full text-left">
            <thead class="bg-gray-50 border-b border-gray-100">
                <tr>
                    <th class="p-6 font-bold text-gray-400 text-sm uppercase">ID</th>
                    <th class="p-6 font-bold text-gray-400 text-sm uppercase">Produk</th>
                    <th class="p-6 font-bold text-gray-400 text-sm uppercase">Harga</th>
                    <th class="p-6 font-bold text-gray-400 text-sm uppercase">Stok</th>
                    <th class="p-6 font-bold text-gray-400 text-sm uppercase">Aksi</th>
                </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
                <tr v-for="p in products" :key="p.id" class="hover:bg-gray-50 transition-colors">
                    <td class="p-6 text-gray-500">#{{ p.id }}</td>
                    <td class="p-6 font-bold text-gray-800">{{ p.name }}</td>
                    <td class="p-6 text-neon-green font-bold">RP {{ p.price.toLocaleString() }}</td>
                    <td class="p-6 text-gray-600">{{ p.stock }}</td>
                    <td class="p-6 space-x-4">
                        <button @click="openEdit(p)" class="bg-red-600 text-white font-bold px-4 py-2 rounded-full shadow hover:bg-red-400 transition-colors">Edit</button>
                        <button @click="deleteProduct(p.id)" class="bg-red-600 text-white font-bold px-4 py-2 rounded-full shadow hover:bg-red-400 transition-colors">Hapus</button>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
  </div>
</template>