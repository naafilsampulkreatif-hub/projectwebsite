<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
</script>

<template>
  <nav class="bg-white text-black p-4 shadow-md sticky top-0 z-50">
    <div class="container mx-auto flex justify-between items-center">
      <div class="flex items-center space-x-2">
         <span class="text-2xl font-bold tracking-tighter uppercase text-black">
           <span class="text-red-600">NS</span> Night Stalkers
         </span>
      </div>
      <div class="hidden md:flex space-x-8 font-semibold text-sm uppercase tracking-wide">
        <RouterLink to="/" class="hover:text-neon-green transition-colors">Home</RouterLink>
        <RouterLink to="/cart" class="hover:text-neon-green transition-colors">Keranjang</RouterLink>
        <RouterLink to="/" class="hover:text-neon-green transition-colors">Toko</RouterLink>
        <RouterLink to="/orders" class="hover:text-neon-green transition-colors">Pesanan</RouterLink>

        <RouterLink v-if="authStore.isAdmin" to="/admin" class="hover:text-neon-green transition-colors text-red-600">Dashboard</RouterLink>
      </div>
    </div>
  </nav>
</template>
