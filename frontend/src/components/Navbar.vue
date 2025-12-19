<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
</script>

<template>
  <nav class="bg-white text-black p-4 shadow-md sticky top-0 z-50">
    <div class="container mx-auto flex justify-between items-center">
      <div class="flex items-center space-x-3">
        <RouterLink to="/" class="flex items-center space-x-3 group">
          <div class="overflow-hidden rounded-md">
            <img src="/NS-logo.png" alt="NS" class="h-10 w-10 md:h-12 md:w-12 lg:h-14 lg:w-14 object-cover block" />
          </div>
          <span class="text-2xl md:text-3xl font-bold tracking-tighter uppercase group-hover:text-neon-green transition-colors">
            <span class="text-red-600">NIGHT</span>
            <span class="text-black"> STALKERS</span>
          </span>
        </RouterLink>
      </div>
      <div class="hidden md:flex space-x-8 font-semibold text-sm uppercase tracking-wide">
        <RouterLink to="/" class="hover:text-neon-green transition-colors">Home</RouterLink>
        <RouterLink to="/cart" class="hover:text-neon-green transition-colors">Keranjang</RouterLink>
        <RouterLink to="/shop" class="hover:text-neon-green transition-colors">Toko</RouterLink>
      </div>
    </div>
  </nav>
</template>
