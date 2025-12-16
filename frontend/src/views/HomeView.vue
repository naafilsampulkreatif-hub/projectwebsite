<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../services/api'
import { useRouter } from 'vue-router'

const products = ref<any[]>([])
const router = useRouter()

onMounted(async () => {
  try {
    const res = await api.get('/products')
    products.value = res.data
  } catch (e) {
    console.error("Failed to load products", e)
    // Fallback data for demo if DB is empty
    if (products.value.length === 0) {
      products.value = [
         { id: 1, name: 'Produk 1', description: 'Kualitas terbaik', price: 5000, image_url: '' },
         { id: 2, name: 'Produk 2', description: 'Sangat bagus', price: 15000, image_url: '' },
         { id: 3, name: 'Produk 3', description: 'Murah meriah', price: 25000, image_url: '' },
         { id: 4, name: 'Produk 4', description: 'Limited Edition', price: 50000, image_url: '' },
      ]
    }
  }
})

const addToCart = async (product: any) => {
  try {
     await api.post('/cart', {
        product_id: product.id,
        quantity: 1
     })
     alert('Produk ditambahkan ke keranjang!')
  } catch (e) {
     console.error(e)
     alert('Gagal menambahkan ke keranjang')
  }
}

const getImageUrl = (url: string) => url || 'https://via.placeholder.com/300x200/00FF00/FFFFFF?text=Product'
</script>

<template>
  <div class="bg-[#F3F3F3] text-black pb-20">
    <!-- Hero Section -->
    <div class="bg-white pb-2">
       <div class="bg-neon-green h-[500px] w-full relative flex items-center justify-center">
          <h1 class="text-6xl md:text-8xl font-black text-white tracking-tighter opacity-50 select-none">
            NIGHT STALKERS
          </h1>
       </div>
    </div>

    <!-- Features Section (Overlap) -->
    <section class="container mx-auto px-4 -mt-24 relative z-10 grid grid-cols-1 md:grid-cols-3 gap-8 mb-24">
       <div v-for="i in 3" :key="i" class="bg-white p-10 rounded-[2rem] shadow-xl flex flex-col items-center text-center transition-transform hover:-translate-y-2">
          <div class="w-20 h-20 bg-green-100 rounded-full mb-6 flex items-center justify-center text-neon-green">
             <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
               <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
             </svg>
          </div>
          <h3 class="font-bold text-xl mb-3 text-gray-800">LOREM IPSUM</h3>
          <p class="text-gray-500 text-sm leading-relaxed">
             The quick fox jumps over the lazy dog. The quick fox jumps over the lazy dog.
          </p>
       </div>
    </section>

    <!-- Synopsis Section -->
    <section class="container mx-auto px-4 mb-32">
       <div class="flex flex-col md:flex-row gap-12 items-center">
          <div class="flex-1 bg-white p-12 rounded-[2rem] shadow-lg relative overflow-hidden">
             <div class="absolute top-0 left-0 w-2 h-full bg-red-500"></div>
             <h2 class="text-xl font-bold mb-2 uppercase text-red-500 tracking-widest">Sinopsis</h2>
             <h3 class="text-5xl font-extrabold mb-8 text-gray-900">Night Stalkers</h3>
             <p class="text-gray-600 mb-8 leading-loose">
               Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam.
             </p>
             <button class="text-neon-green font-bold uppercase tracking-widest flex items-center hover:underline">
               Selengkapnya
               <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 ml-2" viewBox="0 0 20 20" fill="currentColor">
                 <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
               </svg>
             </button>
          </div>
          <div class="flex-1 h-96 bg-neon-green rounded-[2rem] shadow-inner opacity-80"></div>
       </div>
    </section>

    <!-- About Section (Tentang) -->
    <section class="container mx-auto px-4 mb-32">
       <div class="flex flex-col md:flex-row-reverse gap-12 items-center">
          <div class="flex-1 bg-white p-12 rounded-[2rem] shadow-lg relative overflow-hidden">
             <div class="absolute top-0 right-0 w-2 h-full bg-red-500"></div>
             <h2 class="text-xl font-bold mb-2 uppercase text-red-500 tracking-widest text-right">Tentang</h2>
             <h3 class="text-5xl font-extrabold mb-8 text-gray-900 text-right">Night Stalkers</h3>
             <p class="text-gray-600 mb-8 leading-loose text-right">
               Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
             </p>
             <div class="flex justify-end">
                <button class="text-neon-green font-bold uppercase tracking-widest flex items-center hover:underline">
                   Selengkapnya
                   <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 ml-2" viewBox="0 0 20 20" fill="currentColor">
                     <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                   </svg>
                </button>
             </div>
          </div>
          <div class="flex-1 h-96 bg-neon-green rounded-[2rem] shadow-inner opacity-80"></div>
       </div>
    </section>

    <!-- Shop Section -->
    <section class="container mx-auto px-4 mb-32">
       <div class="text-center mb-16">
          <h2 class="text-4xl font-extrabold mb-4">TOKO <span class="text-neon-green">Koleksi Kami</span></h2>
          <p class="text-gray-500 mt-2 text-lg">Find your favorite items</p>
       </div>

       <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div v-for="product in products" :key="product.id" class="bg-white rounded-[2rem] shadow-lg overflow-hidden hover:shadow-2xl transition-all duration-300 group">
             <div class="h-64 bg-gray-100 relative overflow-hidden">
                <div class="absolute inset-0 bg-neon-green opacity-0 group-hover:opacity-20 transition-opacity"></div>
                <img :src="getImageUrl(product.image_url)" class="w-full h-full object-cover transform group-hover:scale-110 transition-transform duration-500" />
             </div>
             <div class="p-6 text-center">
                <h3 class="font-bold text-xl mb-2 text-gray-800">{{ product.name }}</h3>
                <p class="text-sm text-gray-500 mb-4 line-clamp-2">{{ product.description }}</p>
                <div class="text-neon-green font-extrabold text-2xl mb-6">RP {{ product.price.toLocaleString() }}</div>
                <button @click="addToCart(product)" class="bg-neon-green text-white px-8 py-3 rounded-full font-bold uppercase tracking-wider hover:bg-[#00cc00] transition-colors shadow-lg hover:shadow-neon-green/50">
                   Beli
                </button>
             </div>
          </div>
       </div>
    </section>

    <!-- Testimonials -->
    <section class="container mx-auto px-4 mb-20">
       <div class="text-center mb-16">
          <h2 class="text-3xl font-bold">Testimoni</h2>
       </div>
       <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div v-for="i in 3" :key="i" class="bg-neon-green p-10 rounded-[2rem] shadow-xl text-white relative mt-12 transition-transform hover:-translate-y-2">
             <div class="w-24 h-24 bg-white rounded-full absolute -top-12 left-1/2 transform -translate-x-1/2 border-4 border-neon-green shadow-sm overflow-hidden flex items-center justify-center">
                <span class="text-neon-green text-3xl font-bold">U{{i}}</span>
             </div>
             <h3 class="mt-12 font-bold text-xl text-center">Orang {{ i }}</h3>
             <p class="text-sm opacity-90 mt-4 text-center italic leading-relaxed">
               "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
             </p>
          </div>
       </div>
    </section>

  </div>
</template>
