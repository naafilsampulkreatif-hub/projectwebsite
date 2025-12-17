<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../services/api'
import { useToastStore } from '../stores/toast'

const products = ref<any[]>([])
const productsContainer = ref<HTMLElement | null>(null)

const scrollProducts = (direction: 'left' | 'right') => {
   const el = productsContainer.value
   if (!el) return
   const offset = el.clientWidth * 0.7
   const next = direction === 'left' ? el.scrollLeft - offset : el.scrollLeft + offset
   el.scrollTo({ left: next, behavior: 'smooth' })
}
const toast = useToastStore()

// Comments (stored in localStorage)
const commenterName = ref('')
const commenterEmail = ref('')
const commentText = ref('')
const comments = ref<any[]>([])

// Images for placeholders (change these URLs to real uploads or assets)
// Use banner.png placed in `frontend/public/banner.png` (served at `/banner.png`)
const heroImage = ref('/banner.png')
const synopsisImage = ref('/sipnosis.png')
const aboutImage = ref('/tentang.png')
const testimonialImages = ref<string[]>([
   'https://images.unsplash.com/photo-1545996124-1b9a0f0d3bde?auto=format&fit=crop&w=256&q=60',
   'https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=crop&w=256&q=60',
   'https://images.unsplash.com/photo-1547425260-76bcadfb4f2c?auto=format&fit=crop&w=256&q=60'
])

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

   // load saved comments
   try {
      const saved = localStorage.getItem('home_comments')
      if (saved) comments.value = JSON.parse(saved)
   } catch (e) {
      console.error('Failed to load comments', e)
   }
})

const addToCart = async (product: any) => {
  try {
     await api.post('/cart', {
        product_id: product.id,
        quantity: 1
     })
     toast.show('Produk ditambahkan ke keranjang!', 'success')
  } catch (e) {
     console.error(e)
     toast.show('Gagal menambahkan ke keranjang', 'error')
  }
}

const getImageUrl = (url: string) => url || 'https://via.placeholder.com/300x200/00FF00/FFFFFF?text=Product'

const submitComment = () => {
   if (!commenterName.value.trim() || !commenterEmail.value.trim() || !commentText.value.trim()) {
      toast.show('Mohon isi nama, email, dan komentar', 'error')
      return
   }
   if (!commenterEmail.value.includes('@')) {
      toast.show('Email tidak valid', 'error')
      return
   }

   const c = {
      name: commenterName.value.trim(),
      email: commenterEmail.value.trim(),
      text: commentText.value.trim(),
      created_at: new Date().toISOString(),
   }

   comments.value.unshift(c)
   try {
      localStorage.setItem('home_comments', JSON.stringify(comments.value))
   } catch (e) {
      console.error('Failed to save comment', e)
   }

   commenterName.value = ''
   commenterEmail.value = ''
   commentText.value = ''
   toast.show('Komentar berhasil ditambahkan', 'success')
}
</script>

<template>
  <div class="bg-[#F3F3F3] text-black pb-20">
    <!-- Hero Section -->
    <div class="bg-white pb-2">
          <div class="relative h-[500px] w-full overflow-hidden">
             <img :src="heroImage" alt="Hero" class="absolute inset-0 w-full h-full object-cover" />
             <div class="absolute inset-0 bg-neon-green opacity-40"></div>
             <div class="relative z-10 flex items-center justify-center h-full">
                <h1 class="text-6xl md:text-8xl font-black text-white tracking-tighter opacity-90 select-none">HITAM</h1>
             </div>
          </div>
    </div>

   <!-- Features Section (Overlap) -->
   <section class="max-w-5xl mx-auto px-6 lg:px-8 -mt-16 relative z-10 grid grid-cols-1 md:grid-cols-3 gap-6 mb-16 bg-[#fbfbfb] rounded-[1rem] py-6">
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
   <section class="max-w-5xl mx-auto px-6 lg:px-8 mb-12 bg-[#fbfbfb] rounded-[1rem] py-6">
       <div class="flex flex-col md:flex-row gap-12 items-center">
          <div class="md:w-8/12 w-full bg-white p-12 rounded-[2rem] shadow-lg relative overflow-hidden">
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
          <div class="md:w-4/12 w-full rounded-[1.25rem] shadow-inner overflow-hidden flex items-center justify-center bg-gray-100">
             <img :src="synopsisImage" alt="Synopsis" class="w-full h-56 md:h-72 lg:h-80 object-cover block" />
          </div>
       </div>
    </section>

   <!-- About Section (Tentang) -->
   <section class="max-w-5xl mx-auto px-6 lg:px-8 mb-12 bg-[#fbfbfb] rounded-[1rem] py-6">
       <div class="flex flex-col md:flex-row-reverse gap-12 items-center">
          <div class="md:w-8/12 w-full bg-white p-12 rounded-[2rem] shadow-lg relative overflow-hidden">
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
          <div class="md:w-4/12 w-full rounded-[1.25rem] shadow-inner overflow-hidden flex items-center justify-center bg-gray-100">
             <img :src="aboutImage" alt="About" class="w-full h-56 md:h-72 lg:h-80 object-cover block" />
          </div>
       </div>
    </section>

   <!-- Shop Section -->
   <section class="max-w-5xl mx-auto px-6 lg:px-8 mb-12 bg-[#fbfbfb] rounded-[1rem] py-6">
       <div class="text-center mb-16">
          <h2 class="text-4xl font-extrabold mb-4">Koleksi <span class="text-neon-green">Kami</span></h2>
          <p class="text-gray-500 mt-2 text-lg">Temukan yang menarik bagimu</p>
       </div>

          <div class="relative">
             <div class="flex justify-end mb-4 gap-2">
                <button @click="scrollProducts('left')" class="bg-white border rounded-full p-2 shadow hover:bg-gray-50"><</button>
                <button @click="scrollProducts('right')" class="bg-white border rounded-full p-2 shadow hover:bg-gray-50">></button>
             </div>

             <div ref="productsContainer" class="flex gap-6 overflow-x-auto snap-x snap-mandatory scroll-smooth pb-4">
                <div v-for="product in products" :key="product.id" class="min-w-[240px] md:min-w-[260px] lg:min-w-[300px] snap-start bg-white rounded-[1.25rem] shadow-lg overflow-hidden hover:shadow-2xl transition-all duration-300 group">
                   <div class="h-56 bg-gray-100 relative overflow-hidden">
                      <div class="absolute inset-0 bg-neon-green opacity-0 group-hover:opacity-20 transition-opacity"></div>
                      <img :src="getImageUrl(product.image_url)" class="w-full h-full object-cover transform group-hover:scale-110 transition-transform duration-500" />
                   </div>
                   <div class="p-4 text-center">
                      <h3 class="font-bold text-lg mb-1 text-gray-800">{{ product.name }}</h3>
                      <p class="text-sm text-gray-500 mb-2 line-clamp-2">{{ product.description }}</p>
                      <div class="text-neon-green font-extrabold text-xl mb-3">RP {{ product.price.toLocaleString() }}</div>
                      <button @click="addToCart(product)" class="bg-neon-green text-black px-4 py-2 rounded-full font-bold uppercase tracking-wider hover:bg-[#00cc00] transition-colors shadow-sm">Beli</button>
                   </div>
                </div>
             </div>
          </div>
    </section>

   <!-- Testimonials -->
   <section class="max-w-5xl mx-auto px-6 lg:px-8 mb-12 bg-[#fbfbfb] rounded-[1rem] py-6">
       <div class="text-center mb-16">
          <h2 class="text-3xl font-bold">Testimoni</h2>
       </div>
       <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div v-for="(img, idx) in testimonialImages" :key="idx" class="bg-neon-green p-10 rounded-[2rem] shadow-xl text-white relative mt-12 transition-transform hover:-translate-y-2">
             <div class="w-24 h-24 bg-white rounded-full absolute -top-12 left-1/2 transform -translate-x-1/2 border-4 border-neon-green shadow-sm overflow-hidden flex items-center justify-center">
                <img :src="img" alt="avatar" class="w-full h-full object-cover" />
             </div>
             <h3 class="mt-12 font-bold text-xl text-center">Orang {{ idx + 1 }}</h3>
             <p class="text-sm opacity-90 mt-4 text-center italic leading-relaxed text-black">
               "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
             </p>
          </div>
       </div>
    </section>

      <!-- Comment Form -->
      <section class="max-w-5xl mx-auto px-6 lg:px-8 mb-20">
         <div class="max-w-5xl mx-auto bg-white p-8 rounded-[1rem] shadow-lg">
            <h3 class="text-2xl font-bold mb-4">Komentar</h3>

            <form @submit.prevent="submitComment" class="space-y-4">
               <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <input v-model="commenterName" type="text" placeholder="Nama" class="w-full bg-gray-50 border border-gray-200 p-3 rounded-lg" />
                  <input v-model="commenterEmail" type="email" placeholder="Email" class="w-full bg-gray-50 border border-gray-200 p-3 rounded-lg" />
               </div>
               <div>
                  <textarea v-model="commentText" rows="4" placeholder="Tulis komentar Anda..." class="w-full bg-gray-50 border border-gray-200 p-3 rounded-lg"></textarea>
               </div>
               <div class="flex justify-end">
                  <button type="submit" class="bg-neon-green text-black px-6 py-2 rounded-full font-bold hover:bg-[#00cc00]">Kirim Komentar</button>
               </div>
            </form>

            <div class="mt-6">
               <h4 class="font-bold mb-3">Komentar Terbaru</h4>
               <div v-if="comments.length === 0" class="text-gray-500">Belum ada komentar.</div>
               <div v-else class="space-y-4">
                  <div v-for="(c, idx) in comments" :key="idx" class="border border-gray-100 p-4 rounded-lg">
                     <div class="flex justify-between items-center mb-2">
                        <div class="font-bold">{{ c.name }}</div>
                        <div class="text-xs text-gray-400">{{ new Date(c.created_at).toLocaleString() }}</div>
                     </div>
                     <div class="text-sm text-gray-700">{{ c.text }}</div>
                  </div>
               </div>
            </div>
         </div>
      </section>

  </div>
</template>
