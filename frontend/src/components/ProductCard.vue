<script setup lang="ts">
import { useCartStore } from '../stores/cart';
import { useToastStore } from '../stores/toast';
import { ref } from 'vue';

const props = defineProps<{
  product: any;
}>();

const emit = defineEmits<{
  (e: 'added', product: any): void
}>();

const cartStore = useCartStore();
const toast = useToastStore();
const loading = ref(false);

const addToCart = async () => {
  if (loading.value) return;
  loading.value = true;
  try {
    // Allow guests to add to cart; session_id is created automatically by api.ts
    await cartStore.addToCart(props.product.id);
    // Refresh cart to ensure UI reflects new items
    await cartStore.fetchCart();
    // Notify parent components that the product was added
    emit('added', true);
    toast.show('Produk ditambahkan ke keranjang', 'success');
  } catch (e) {
    // Emit failure event for parent to handle if needed
    emit('added', false);
    console.error('Add to cart failed', e);
    const msg = (e as any)?.response?.data?.message || (e as any)?.message || 'Gagal menambah ke keranjang';
    toast.show(String(msg), 'error');
  } finally {
    loading.value = false;
  }
};

// Helper: ensure image URL is absolute or starts with '/'
function formatImageUrl(url: string) {
  if (!url) return url;
  if (url.startsWith('http://') || url.startsWith('https://')) return url;
  if (url.startsWith('/')) return url;
  // If URL contains 'uploads', ensure it starts with /
  if (url.includes('uploads')) {
    return url.startsWith('/') ? url : '/' + url;
  }
  return '/' + url;
}

// Format price safely
function formatPrice(v: any) {
  const n = Number(v) || 0;
  return n.toLocaleString();
}
</script>

<template>
  <div class="w-full bg-white group overflow-hidden rounded-2xl shadow-md flex flex-col hover:shadow-xl transition-shadow duration-300">
    <!-- Image area -->
    <div class="w-full aspect-[3/4] bg-gray-100 overflow-hidden rounded-t-2xl flex items-center justify-center">
      <img
        :src="(product.image_url && formatImageUrl(product.image_url)) || 'https://via.placeholder.com/400x533?text=Product'"
        :alt="product.name"
        class="w-full h-full object-cover object-center transform group-hover:scale-105 transition-transform duration-300"
        @error="console.log('Image failed to load:', product.image_url)"
      />
    </div>

    <!-- Content area -->
    <div class="p-4 flex-1 flex flex-col">
      <div class="mb-3">
        <div class="text-xs text-gray-400 mb-1">{{ product.category || 'Kategori' }}</div>
        <h3 class="text-lg text-gray-900 font-semibold mb-2 line-clamp-2">{{ product.name }}</h3>
        <div class="text-black font-extrabold text-lg">RP {{ formatPrice(product.price) }}</div>
      </div>

      <div class="mt-auto">
        <button :disabled="loading" @click="addToCart" aria-label="Add to cart" 
          class="w-full text-white rounded-full font-bold py-2 px-4 mt-2 transition-colors duration-200 disabled:opacity-60 disabled:cursor-not-allowed hover:opacity-80" 
          style="background-color: #547792;">
          <span v-if="loading">Menambahkan...</span>
          <span v-else>Tambah ke Keranjang</span>
        </button>
      </div>
    </div>
  </div>
</template>
