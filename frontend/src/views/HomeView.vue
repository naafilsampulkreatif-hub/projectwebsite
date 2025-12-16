<script setup lang="ts">
import { onMounted, ref } from 'vue';
import Header from '../components/Header.vue';
import ProductCard from '../components/ProductCard.vue';
import api from '../services/api';

const products = ref<any[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const response = await api.get('/products');
    products.value = response.data;
  } catch (error) {
    console.error("Failed to fetch products", error);
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="min-h-screen bg-white">
    <Header />

    <!-- Hero Section -->
    <section class="relative bg-gray-100 h-[500px] flex items-center">
        <div class="container mx-auto px-4 grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
            <div>
                <p class="text-primary uppercase tracking-widest mb-2">New Arrivals</p>
                <h1 class="text-5xl font-bold text-gray-800 mb-4">Beats Studio3 Wireless</h1>
                <h2 class="text-4xl font-light text-gray-600 mb-6">Experience your music like never before.</h2>
                <div class="flex items-center space-x-4">
                    <span class="text-3xl font-bold text-primary">$299.95</span>
                    <button class="bg-gray-800 text-white px-8 py-3 rounded hover:bg-primary transition-colors">
                        Shop Now
                    </button>
                </div>
            </div>
            <div class="hidden md:block">
               <!-- Placeholder for Hero Image -->
               <img src="https://via.placeholder.com/500" alt="Headphones" class="w-full max-w-md mx-auto">
            </div>
        </div>
    </section>

    <!-- Trending Products -->
    <section class="py-16">
        <div class="container mx-auto px-4">
            <h2 class="text-3xl font-bold text-center mb-12">Trending Products</h2>

            <div v-if="loading" class="text-center py-12">Loading...</div>

            <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
                <ProductCard v-for="product in products" :key="product.id" :product="product" />
            </div>
        </div>
    </section>

    <footer class="bg-gray-800 text-white py-12">
        <div class="container mx-auto px-4 text-center">
            <p>&copy; 2025 Molla Store. All rights reserved.</p>
        </div>
    </footer>
  </div>
</template>
