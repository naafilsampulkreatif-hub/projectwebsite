<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Header from '../components/Header.vue';
import api from '../services/api';

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
    const res = await api.get('/products');
    products.value = res.data;
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
        if (isEditing.value) {
            await api.put(`/admin/products/${form.value.id}`, form.value);
        } else {
            await api.post('/admin/products', form.value);
        }
        await fetchProducts();
        resetForm();
    } catch (e) {
        alert("Error saving product");
    }
};

const deleteProduct = async (id: number) => {
    if (!confirm("Are you sure?")) return;
    try {
        await api.delete(`/admin/products/${id}`);
        await fetchProducts();
    } catch (e) {
        alert("Error deleting product");
    }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <Header />
    <div class="container mx-auto px-4 py-12">
        <div class="flex justify-between items-center mb-8">
            <h1 class="text-3xl font-bold">Admin Dashboard</h1>
            <button @click="openAdd" class="bg-primary text-white px-4 py-2 rounded hover:bg-orange-600">Add Product</button>
        </div>

        <!-- Product Form -->
        <div v-if="showForm" class="bg-white p-6 rounded shadow mb-8">
            <h2 class="text-xl font-bold mb-4">{{ isEditing ? 'Edit' : 'Add' }} Product</h2>
            <form @submit.prevent="saveProduct" class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <input v-model="form.name" placeholder="Name" class="border p-2 rounded" required>
                <input v-model="form.slug" placeholder="Slug" class="border p-2 rounded" required>
                <input v-model="form.price" type="number" step="0.01" placeholder="Price" class="border p-2 rounded" required>
                <input v-model="form.stock" type="number" placeholder="Stock" class="border p-2 rounded" required>
                <input v-model="form.image_url" placeholder="Image URL" class="border p-2 rounded">
                <input v-model="form.category_id" type="number" placeholder="Category ID" class="border p-2 rounded">
                <textarea v-model="form.description" placeholder="Description" class="border p-2 rounded md:col-span-2"></textarea>

                <div class="md:col-span-2 flex space-x-4">
                    <button type="submit" class="bg-green-500 text-white px-4 py-2 rounded">Save</button>
                    <button type="button" @click="resetForm" class="bg-gray-300 text-black px-4 py-2 rounded">Cancel</button>
                </div>
            </form>
        </div>

        <!-- Product List -->
        <div class="bg-white rounded shadow overflow-x-auto">
            <table class="w-full text-left">
                <thead class="bg-gray-100 border-b">
                    <tr>
                        <th class="p-4">ID</th>
                        <th class="p-4">Name</th>
                        <th class="p-4">Price</th>
                        <th class="p-4">Stock</th>
                        <th class="p-4">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="p in products" :key="p.id" class="border-b">
                        <td class="p-4">{{ p.id }}</td>
                        <td class="p-4">{{ p.name }}</td>
                        <td class="p-4">${{ p.price }}</td>
                        <td class="p-4">{{ p.stock }}</td>
                        <td class="p-4 space-x-2">
                            <button @click="openEdit(p)" class="text-blue-500 hover:underline">Edit</button>
                            <button @click="deleteProduct(p.id)" class="text-red-500 hover:underline">Delete</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
  </div>
</template>
