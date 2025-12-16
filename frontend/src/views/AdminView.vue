<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '../services/api';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

const products = ref<any[]>([]);
const isEditing = ref(false);
const showForm = ref(false);

const form = ref({
    id: 0,
    name: '',
    slug: '',
    description: '',
    price: 0,
    stock: 0,
    image_url: '',
    category_id: 1
});

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
        slug: '',
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
             stock: Number(form.value.stock),
             category_id: Number(form.value.category_id)
        };

        if (isEditing.value) {
            await api.put(`/admin/products/${form.value.id}`, productData);
        } else {
            await api.post('/admin/products', productData);
        }
        await fetchProducts();
        resetForm();
    } catch (e) {
        alert("Error saving product");
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
</script>

<template>
  <div class="min-h-screen bg-[#F3F3F3] flex text-black">
     <!-- Sidebar -->
     <aside class="w-64 bg-[#1a1a1a] text-white min-h-screen p-8 hidden md:flex flex-col fixed h-full">
        <h2 class="text-3xl font-extrabold mb-12 text-neon-green tracking-tighter">NS ADMIN</h2>
        <ul class="space-y-6 font-bold tracking-wide">
           <li><a href="#" class="block text-neon-green flex items-center"><span class="w-2 h-2 bg-neon-green rounded-full mr-3"></span> Products</a></li>
           <li><a href="#" class="block text-gray-400 hover:text-white transition-colors">Orders</a></li>
           <li><a href="#" class="block text-gray-400 hover:text-white transition-colors">Users</a></li>
           <li><a href="#" class="block text-gray-400 hover:text-white transition-colors">Settings</a></li>
        </ul>
        <div class="mt-auto">
           <button @click="authStore.logout(); router.push('/login')" class="text-red-500 hover:text-red-400 font-bold uppercase text-sm tracking-widest">Logout</button>
        </div>
     </aside>

     <!-- Content -->
     <main class="flex-1 p-8 md:ml-64">
        <div class="flex justify-between items-center mb-10">
            <h1 class="text-3xl font-bold">Dashboard Produk</h1>
            <button @click="openAdd" class="bg-neon-green text-white px-6 py-3 rounded-full font-bold uppercase shadow-lg hover:shadow-neon-green/50 transition-all hover:-translate-y-1">
                + Tambah Produk
            </button>
        </div>

        <!-- Product Form -->
        <div v-if="showForm" class="bg-white p-8 rounded-[2rem] shadow-xl mb-10 border border-gray-100">
            <h2 class="text-2xl font-bold mb-6 border-b pb-2">{{ isEditing ? 'Edit' : 'Tambah' }} Produk</h2>
            <form @submit.prevent="saveProduct" class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <input v-model="form.name" placeholder="Nama Produk" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green" required>
                <input v-model="form.slug" placeholder="Slug (URL)" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green" required>
                <input v-model="form.price" type="number" step="100" placeholder="Harga (IDR)" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green" required>
                <input v-model="form.stock" type="number" placeholder="Stok" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green" required>
                <input v-model="form.image_url" placeholder="URL Gambar" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green">
                <input v-model="form.category_id" type="number" placeholder="ID Kategori" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green">
                <textarea v-model="form.description" placeholder="Deskripsi" class="bg-gray-50 border border-gray-200 p-4 rounded-xl focus:outline-none focus:border-neon-green md:col-span-2 h-32"></textarea>

                <div class="md:col-span-2 flex space-x-4 justify-end">
                    <button type="button" @click="resetForm" class="px-6 py-3 rounded-full font-bold text-gray-500 hover:bg-gray-100">Batal</button>
                    <button type="submit" class="bg-neon-green text-black px-8 py-3 rounded-full font-bold shadow-lg hover:bg-[#00cc00]">Simpan</button>
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
                            <button @click="openEdit(p)" class="text-blue-500 font-bold hover:underline">Edit</button>
                            <button @click="deleteProduct(p.id)" class="text-red-500 font-bold hover:underline">Hapus</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
     </main>
  </div>
</template>
