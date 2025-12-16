<script setup lang="ts">
import { useCartStore } from '../stores/cart';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const props = defineProps<{
  product: any;
}>();

const cartStore = useCartStore();
const authStore = useAuthStore();
const router = useRouter();

const addToCart = () => {
  if (!authStore.isAuthenticated) {
    router.push('/login');
    return;
  }
  cartStore.addToCart(props.product.id);
};
</script>

<template>
  <div class="bg-white group overflow-hidden">
    <div class="relative h-64 bg-gray-100 mb-4 overflow-hidden">
        <img :src="product.image_url || 'https://via.placeholder.com/300'" :alt="product.name" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300">

        <div class="absolute bottom-0 left-0 w-full p-4 translate-y-full group-hover:translate-y-0 transition-transform duration-300">
            <button @click="addToCart" class="w-full bg-primary text-white py-2 hover:bg-orange-600 transition-colors">
                Add to Cart
            </button>
        </div>
    </div>

    <div class="text-center">
        <p class="text-gray-500 text-sm mb-1">Category</p>
        <h3 class="text-lg text-gray-800 font-medium mb-1">{{ product.name }}</h3>
        <p class="text-primary font-bold">${{ product.price }}</p>
    </div>
  </div>
</template>
