<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../services/api'

const products = ref<any[]>([])

onMounted(async () => {
  try {
    const res = await api.get('/products')
    products.value = res.data || []
  } catch (e) {
    console.error('Failed to load products for shop', e)
  }
})

const getImageUrl = (url: string) => url || 'https://via.placeholder.com/300x200/00FF00/FFFFFF?text=Product'
</script>

<template>
  <div class="max-w-6xl mx-auto p-6">
    <h1 class="text-4xl font-extrabold mb-6">Toko Kami</h1>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <div v-for="p in products" :key="p.id" class="bg-white rounded-lg shadow p-4">
        <img :src="getImageUrl(p.image_url)" class="w-full h-44 object-cover rounded mb-4" />
        <h3 class="font-bold text-lg">{{ p.name }}</h3>
        <p class="text-gray-600">{{ p.description }}</p>
        <div class="mt-4 flex justify-between items-center">
          <div class="text-neon-green font-bold">Rp {{ p.price.toLocaleString() }}</div>
          <button @click.prevent class="bg-neon-green text-black px-4 py-2 rounded-full font-bold">Beli</button>
        </div>
      </div>
    </div>
  </div>
</template>
