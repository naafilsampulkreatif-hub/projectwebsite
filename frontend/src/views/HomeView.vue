<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../services/api'
import { useToastStore } from '../stores/toast'
import ProductCard from '../components/ProductCard.vue'

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
   'Grace.jpg',
   'rouf.jpg',
   'zildan.jpg'
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

  // Load comments from backend
  try {
    const res = await api.get('/comments')
    comments.value = res.data || []
  } catch (e) {
    console.error('Failed to load comments', e)
    // Fallback to localStorage if backend fails
    try {
      const saved = localStorage.getItem('home_comments')
      if (saved) comments.value = JSON.parse(saved)
    } catch (ex) {
      console.error('Failed to load comments from storage', ex)
    }
  }
})



const submitComment = async () => {
   if (!commenterName.value.trim() || !commenterEmail.value.trim() || !commentText.value.trim()) {
      toast.show('Mohon isi nama, email, dan komentar', 'error')
      return
   }
   if (!commenterEmail.value.includes('@')) {
      toast.show('Email tidak valid', 'error')
      return
   }

   try {
      // Submit to backend API
      const res = await api.post('/comments', {
         name: commenterName.value.trim(),
         email: commenterEmail.value.trim(),
         text: commentText.value.trim()
      })

         toast.show(res.data.message || 'Komentar dipublikasikan', 'success')

         // If backend returned the created comment, add it immediately
         if (res.data && res.data.comment) {
            comments.value.unshift(res.data.comment)
         } else {
            // Reset form
            commenterName.value = ''
            commenterEmail.value = ''
            commentText.value = ''

            // Reload comments from backend
            const getRes = await api.get('/comments')
            comments.value = getRes.data || []
         }

         // Reset form (ensure cleared even if comment added)
         commenterName.value = ''
         commenterEmail.value = ''
         commentText.value = ''
   } catch (e) {
      console.error('Failed to submit comment', e)
         toast.show('Gagal mengirim komentar ke server, disimpan lokal', 'error')
         // fallback: save locally so it's not lost
         const c = {
           name: commenterName.value.trim(),
           email: commenterEmail.value.trim(),
           text: commentText.value.trim(),
           created_at: new Date().toISOString(),
         }
         comments.value.unshift(c)
         try { localStorage.setItem('home_comments', JSON.stringify(comments.value)) } catch (ex) { console.error(ex) }
   }
}
</script>

<template>
  <div class="bg-white text-black pb-20">
    <!-- Hero Section -->
      <div class="bg-white pb-2">
         <div class="relative h-[320px] sm:h-[400px] md:h-[500px] w-full overflow-hidden">
            <img :src="heroImage" alt="Hero" class="absolute inset-0 w-full h-full object-cover" />
            <div class="absolute inset-0 bg-gradient-to-b from-black/60 via-blue-500/40 to-black/20"></div>
            <div class="relative z-10 flex flex-col items-center justify-center h-full px-2">
               <h1 class="text-4xl sm:text-5xl md:text-7xl font-black text-white drop-shadow-lg tracking-tight text-center mb-4 animate-fade-in" style="font-family: 'TT Regard Demo', sans-serif;">
                  <span class="block">
                    <span class="text-red-600">Night</span><span>Stalkers</span>
                  </span>
                  <span class="block text-lg md:text-2xl font-medium text-blue-500 mt-2 tracking-wide">Misteri | Thriller</span>
               </h1>
                  <a href="https://www.webtoons.com/id/canvas/night-stalkers/list?title_no=1089184" target="_blank" rel="noopener" class="inline-block bg-red-600 text-white font-bold px-6 py-2 rounded-full shadow hover:bg-red-400 transition-colors text-base md:text-lg">Baca di Webtoon</a>
            </div>
         </div>
      </div>

   <!-- Features Section (Overlap) -->
    <section class="max-w-5xl mx-auto px-2 sm:px-4 md:px-6 lg:px-8 -mt-10 md:-mt-16 relative z-10 grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-6 mb-10 md:mb-16 bg-[#fbfbfb] rounded-[1rem] py-4 md:py-6">
         <div class="bg-white p-5 sm:p-7 md:p-10 rounded-2xl md:rounded-[2rem] shadow-xl flex flex-col items-center text-center transition-transform hover:-translate-y-2">
            <div class="w-14 h-14 sm:w-16 sm:h-16 md:w-20 md:h-20 bg-yellow-100 rounded-full mb-4 md:mb-6 flex items-center justify-center">
               <!-- Lampu/Ide Icon -->
               <svg viewBox="0 0 24 24" class="h-12 w-12" fill="none">
                  <circle cx="12" cy="12" r="8" fill="#fde047"/>
                  <path d="M12 7v4l3 3" stroke="#f59e42" stroke-width="2" stroke-linecap="round"/>
                  <rect x="10.5" y="17" width="3" height="2" rx="1" fill="#fbbf24"/>
                  <path d="M8 21h8" stroke="#fbbf24" stroke-width="2" stroke-linecap="round"/>
               </svg>
            </div>
            <h3 class="font-bold text-xl mb-3 text-gray-800">Teknologi</h3>
            <p class="text-gray-500 text-sm leading-relaxed">
               NightStalkers dibuat dengan memaksimalkan penggunaan teknologi AI.
            </p>
         </div>
         <div class="bg-white p-5 sm:p-7 md:p-10 rounded-2xl md:rounded-[2rem] shadow-xl flex flex-col items-center text-center transition-transform hover:-translate-y-2">
            <div class="w-14 h-14 sm:w-16 sm:h-16 md:w-20 md:h-20 bg-pink-100 rounded-full mb-4 md:mb-6 flex items-center justify-center">
               <!-- Heart/Chat Icon (Perfect Heart) -->
               <svg viewBox="0 0 24 24" class="h-12 w-12" fill="none">
                 <ellipse cx="12" cy="13" rx="7" ry="5.5" fill="#f9a8d4" fill-opacity=".4"/>
                 <path d="M12 20.5s-6.5-4.5-8.5-7.5C1.5 9.5 4 6.5 7.5 7.5c2 .6 2.5 2.5 2.5 2.5s.5-1.9 2.5-2.5c3.5-1 6 2 4 5.5-2 3-8.5 7.5-8.5 7.5z" fill="#f472b6" stroke="#ec4899" stroke-width="2"/>
               </svg>
            </div>
            <h3 class="font-bold text-xl mb-3 text-gray-800">Respon</h3>
            <p class="text-gray-500 text-sm leading-relaxed">
               NightStalkers mendapat respon yang baik dari pembaca.
            </p>
         </div>
         <div class="bg-white p-5 sm:p-7 md:p-10 rounded-2xl md:rounded-[2rem] shadow-xl flex flex-col items-center text-center transition-transform hover:-translate-y-2">
            <div class="w-14 h-14 sm:w-16 sm:h-16 md:w-20 md:h-20 bg-blue-100 rounded-full mb-4 md:mb-6 flex items-center justify-center">
               <!-- Buku Icon -->
               <svg viewBox="0 0 24 24" class="h-12 w-12" fill="none">
                  <rect x="4" y="5" width="7" height="14" rx="2" fill="#60a5fa" stroke="#2563eb" stroke-width="2"/>
                  <rect x="13" y="5" width="7" height="14" rx="2" fill="#93c5fd" stroke="#2563eb" stroke-width="2"/>
                  <path d="M7 8h3" stroke="#2563eb" stroke-width="2" stroke-linecap="round"/>
                  <path d="M16 8h3" stroke="#2563eb" stroke-width="2" stroke-linecap="round"/>
               </svg>
            </div>
            <h3 class="font-bold text-xl mb-3 text-gray-800">Buku</h3>
            <p class="text-gray-500 text-sm leading-relaxed">
               NightStalkers sudah menerbitkan buku fisiknya.
            </p>
         </div>
      </section>

   <!-- Synopsis Section -->
   <section class="max-w-5xl mx-auto px-6 lg:px-8 mb-12 bg-[#fbfbfb] rounded-[1rem] py-6">
      <section class="max-w-5xl mx-auto px-2 sm:px-4 md:px-6 lg:px-8 mb-8 md:mb-12 bg-[#fbfbfb] rounded-[1rem] py-4 md:py-6">
       <div class="flex flex-col md:flex-row gap-12 items-center">
          <div class="md:w-8/12 w-full bg-white p-12 rounded-[2rem] shadow-lg relative overflow-hidden">
             <div class="absolute top-0 left-0 w-2 h-full bg-red-500"></div>
             <h2 class="text-xl font-bold mb-2 uppercase text-red-500 tracking-widest">Sinopsis</h2>
                   <h3 class="text-5xl font-extrabold mb-8 text-gray-900" style="font-family: 'TT Regard Demo', sans-serif;">
                      <span class="text-red-600">Night</span> Stalkers
                   </h3>
             <p class="text-gray-600 mb-8 leading-loose">
               Ketika murid SMA di sebuah sekolah terpencil bernama Anita mendadak meninggal misterius, Askara dan teman-temannya mendapati serangkaian kejadian aneh mulai menghantui mereka. Mulai dari penampakan, mimpi-mimpi ganjil, hingga munculnya bunga mawar putih di jendela Askara. Saat sahabat mereka, Ifal, juga jatuh sakit dengan gejala serupa, mereka mencium kehadiran sosok misterius yang menghubungkan semua kejadian itu. Dihantui rasa takut, mereka menyelidiki kebenaran di balik kematian Anita.
             </p>
          </div>
          <div class="md:w-4/12 w-full rounded-[1.25rem] shadow-inner overflow-hidden flex items-center justify-center bg-gray-100">
             <img :src="synopsisImage" alt="Synopsis" class="w-full h-56 md:h-72 lg:h-80 object-cover block" />
          </div>
       </div>
    </section>
   </section>

   <!-- About Section (Tentang) -->
   <section class="max-w-5xl mx-auto px-6 lg:px-8 mb-12 bg-[#fbfbfb] rounded-[1rem] py-6">
      <section class="max-w-5xl mx-auto px-2 sm:px-4 md:px-6 lg:px-8 mb-8 md:mb-12 bg-[#fbfbfb] rounded-[1rem] py-4 md:py-6">
       <div class="flex flex-col md:flex-row-reverse gap-12 items-center">
          <div class="md:w-8/12 w-full bg-white p-12 rounded-[2rem] shadow-lg relative overflow-hidden">
             <div class="absolute top-0 right-0 w-2 h-full bg-red-500"></div>
             <h2 class="text-xl font-bold mb-2 uppercase text-red-500 tracking-widest text-right">Tentang</h2>
                   <h3 class="text-5xl font-extrabold mb-8 text-gray-900 text-right" style="font-family: 'TT Regard Demo', sans-serif;">
                      <span class="text-red-600">Night</span> Stalkers
                   </h3>
             <p class="text-gray-600 mb-8 leading-loose text-right">
               NightStalkers dibuat oleh Risya Farizi yang ditenagai oleh keinginannya untuk berkreasi
             </p>
            <div class="flex justify-end">
             </div>
          </div>
          <div class="md:w-4/12 w-full rounded-[1.25rem] shadow-inner overflow-hidden flex items-center justify-center bg-gray-100">
             <img :src="aboutImage" alt="About" class="w-full h-56 md:h-72 lg:h-80 object-cover block" />
          </div>
       </div>
    </section>
   </section>

   <!-- Shop Section -->
   <section class="max-w-5xl mx-auto px-6 lg:px-8 mb-12 bg-[#fbfbfb] rounded-[1rem] py-6">
      <div class="max-w-5xl mx-auto px-2 sm:px-4 md:px-6 lg:px-8 mb-8 md:mb-12 bg-[#fbfbfb] rounded-[1rem] py-4 md:py-6">
       <div class="text-center mb-16">
          <h2 class="text-4xl font-extrabold mb-4">Koleksi <span class="text-blue-500">Kami</span></h2>
          <p class="text-gray-500 mt-2 text-lg">Temukan yang menarik bagimu</p>
       </div>

          <div class="relative">
             <!-- Tombol slide di tengah area produk, transparan dan tidak mengganggu kartu -->
             <div class="absolute inset-y-0 left-0 right-0 flex items-center justify-between px-4 pointer-events-none z-10">
               <button
                 @click="scrollProducts('left')"
                 class="bg-red-600/70 text-white rounded-full p-2 shadow hover:bg-red-400/80 transition-colors pointer-events-auto backdrop-blur-sm"
                 style="opacity:0.7;"
                 aria-label="Geser kiri"
               >
                 &lt;
               </button>
               <button
                 @click="scrollProducts('right')"
                 class="bg-red-600/70 text-white rounded-full p-2 shadow hover:bg-red-400/80 transition-colors pointer-events-auto backdrop-blur-sm"
                 style="opacity:0.7;"
                 aria-label="Geser kanan"
               >
                 &gt;
               </button>
             </div>

             <div ref="productsContainer" class="flex gap-6 overflow-x-auto snap-x snap-mandatory scroll-smooth pb-4">
                <div v-for="product in products" :key="product.id" class="min-w-[280px] max-w-xs w-full flex-shrink-0">
                  <ProductCard :product="product" @added="() => toast.show('Produk ditambahkan ke keranjang', 'success')" />
                </div>
             </div>
          </div>
    </div>
   </section>

   <!-- Testimonials -->
   <section class="max-w-5xl mx-auto px-6 lg:px-8 mb-12 bg-[#fbfbfb] rounded-[1rem] py-6">
      <div class="max-w-5xl mx-auto px-2 sm:px-4 md:px-6 lg:px-8 mb-8 md:mb-12 bg-[#fbfbfb] rounded-[1rem] py-4 md:py-6">
       <div class="text-center mb-16">
          <h2 class="text-3xl font-bold">Testimoni</h2>
       </div>
       <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div v-for="(img, idx) in testimonialImages" :key="idx" class="bg-white p-10 rounded-[2rem] shadow-xl text-gray-900 relative mt-12 transition-transform hover:-translate-y-2 border border-slate-100">
             <div class="w-24 h-24 bg-white rounded-full absolute -top-12 left-1/2 transform -translate-x-1/2 border-4 border-slate-200 shadow-sm overflow-hidden flex items-center justify-center">
                <img :src="img" alt="avatar" class="w-full h-full object-cover" />
             </div>
             <!-- Nama pemberi testimoni -->
             <h3 class="mt-12 font-bold text-xl text-center text-slate-900 md:text-slate-800">
               {{ idx === 0 ? 'Grace' : idx === 1 ? 'Rouf' : 'Zildan' }}
             </h3>
             <p v-if="idx === 0" class="text-sm opacity-90 mt-4 text-center italic leading-relaxed text-slate-700">
                "Novel NightStalkers benar-benar membawa pengalaman membaca yang berbeda. Ceritanya seru dan penuh kejutan!"
             </p>
             <p v-else-if="idx === 1" class="text-sm opacity-90 mt-4 text-center italic leading-relaxed text-slate-700">
                "Komik NightStalkers ilustrasinya keren banget, alur ceritanya juga mudah diikuti dan bikin penasaran. Highly recommended!"
             </p>
             <p v-else class="text-sm opacity-90 mt-4 text-center italic leading-relaxed text-slate-700">
                "Saya suka banget dengan karakter-karakter di NightStalkers. Baik novel maupun komiknya, dua-duanya bikin nagih!"
             </p>
          </div>
       </div>
    </div>
   </section>

   <!-- End of main container -->

      <!-- Comment Form -->
      <section class="max-w-5xl mx-auto px-2 sm:px-4 md:px-6 lg:px-8 mb-20">
          <div class="max-w-5xl mx-auto bg-white p-3 sm:p-4 md:p-8 rounded-[1rem] shadow-lg">
            <h3 class="text-2xl font-bold mb-4">Komentar</h3>

            <form @submit.prevent="submitComment" class="space-y-3 sm:space-y-4">
               <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 sm:gap-4">
                  <input v-model="commenterName" type="text" placeholder="Nama" class="w-full bg-gray-50 border border-gray-200 p-2 sm:p-3 rounded-lg text-sm" />
                  <input v-model="commenterEmail" type="email" placeholder="Email" class="w-full bg-gray-50 border border-gray-200 p-2 sm:p-3 rounded-lg text-sm" />
               </div>
               <div>
                  <textarea v-model="commentText" rows="4" placeholder="Tulis komentar Anda..." class="w-full bg-gray-50 border border-gray-200 p-2 sm:p-3 rounded-lg text-sm resize-y max-h-32 min-h-[48px] overflow-y-auto"></textarea>
               </div>
               <div class="flex flex-col sm:flex-row justify-end gap-2">
                  <button type="submit" class="bg-red-600 text-white px-4 py-2 rounded-full font-bold hover:bg-red-400 transition-colors w-full sm:w-auto">Kirim Komentar</button>
               </div>
            </form>

            <div class="mt-6">
               <h4 class="font-bold mb-3">Komentar Terbaru</h4>
               <div v-if="comments.length === 0" class="text-gray-500">Belum ada komentar.</div>
               <div v-else class="space-y-2 sm:space-y-3 max-h-52 sm:max-h-72 md:max-h-80 overflow-y-auto pr-1">
                  <div v-for="(c, idx) in comments" :key="idx" class="border border-gray-100 p-2 sm:p-3 rounded-lg bg-gray-50">
                     <div class="flex justify-between items-center mb-1 sm:mb-2">
                        <div class="font-bold text-xs sm:text-sm">{{ c.name }}</div>
                        <div class="text-xs text-gray-400">{{ c.created_at ? new Date(c.created_at).toLocaleString() : '' }}</div>
                     </div>
                     <div class="text-gray-700 text-xs sm:text-sm break-words">{{ c.text }}</div>
                  </div>
               </div>
            </div>
         </div>
      </section>
   </div>
</template>
