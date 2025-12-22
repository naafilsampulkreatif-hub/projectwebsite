<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../services/api'
import ProductCard from '../components/ProductCard.vue'
import { useToastStore } from '../stores/toast'

const products = ref<any[]>([])
const toast = useToastStore()

onMounted(async () => {
  try {
    const res = await api.get('/products')
    products.value = res.data || []
  } catch (e) {
    console.error('Failed to load products for shop', e)
  }
})
</script>

<template>
  <div class="min-h-screen bg-white text-black">
    <div class="max-w-6xl mx-auto p-4 md:p-8 lg:p-12">
      <h1 class="text-3xl md:text-4xl font-extrabold mb-4 md:mb-6 text-gray-900">Toko Kami</h1>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3 md:gap-6">
        <ProductCard v-for="p in products" :key="p.id" :product="p" @added="(ok) => ok && toast.show('Produk ditambahkan ke keranjang', 'success')" />
      </div>
    </div>
  </div>
</template>
