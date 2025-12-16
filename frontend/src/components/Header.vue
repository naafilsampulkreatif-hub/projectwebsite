<script setup lang="ts">
import { useAuthStore } from '../stores/auth';
import { useCartStore } from '../stores/cart';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const cartStore = useCartStore();
const router = useRouter();

const logout = () => {
  authStore.logout();
  router.push('/login');
};
</script>

<template>
  <header class="bg-white shadow-sm sticky top-0 z-50">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
      <!-- Logo -->
      <router-link to="/" class="text-2xl font-bold text-gray-800">
        Molla
      </router-link>

      <!-- Navigation -->
      <nav class="hidden md:flex space-x-6">
        <router-link to="/" class="text-gray-600 hover:text-primary">Home</router-link>
        <router-link to="/" class="text-gray-600 hover:text-primary">Shop</router-link>
        <router-link to="/" class="text-gray-600 hover:text-primary">Product</router-link>
        <router-link to="/" class="text-gray-600 hover:text-primary">Pages</router-link>
        <router-link to="/" class="text-gray-600 hover:text-primary">Blog</router-link>
        <router-link to="/" class="text-gray-600 hover:text-primary">Elements</router-link>
      </nav>

      <!-- Actions -->
      <div class="flex items-center space-x-4">
        <!-- Search -->
        <button class="text-gray-600 hover:text-primary">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>

        <!-- Cart -->
        <router-link to="/cart" class="relative text-gray-600 hover:text-primary">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          <span v-if="cartStore.totalItems > 0" class="absolute -top-2 -right-2 bg-primary text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
            {{ cartStore.totalItems }}
          </span>
        </router-link>

        <!-- Auth -->
        <div v-if="authStore.isAuthenticated" class="flex items-center space-x-2">
           <span class="text-sm">Hi, {{ authStore.user?.name }}</span>
           <button @click="logout" class="text-sm text-red-500">Logout</button>
           <router-link v-if="authStore.isAdmin" to="/admin" class="text-sm text-blue-500">Admin</router-link>
        </div>
        <div v-else class="space-x-2">
            <router-link to="/login" class="text-gray-600 hover:text-primary">Login</router-link>
            <router-link to="/register" class="text-gray-600 hover:text-primary">Register</router-link>
        </div>
      </div>
    </div>
  </header>
</template>
