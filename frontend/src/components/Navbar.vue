<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import { computed } from 'vue'

const authStore = useAuthStore()
const cartStore = useCartStore()
const route = useRoute()

const totalItems = computed(() => cartStore.totalItems)
</script>

<template>
  <nav class="bg-white text-slate-800 p-4 shadow-sm border-b border-slate-200 sticky top-0 z-50">
    <div class="container mx-auto flex justify-between items-center">
      <div class="flex items-center space-x-3">
        <RouterLink to="/" class="flex items-center space-x-3 group">
          <div class="overflow-hidden rounded-md">
            <img src="/NS-logo.png" alt="NS" class="h-10 w-10 md:h-12 md:w-12 lg:h-14 lg:w-14 object-cover block" />
          </div>
          <span class="text-2xl md:text-3xl font-bold tracking-tighter uppercase group-hover:text-sky-500 transition-colors" style="font-family: 'TT Regard Demo', sans-serif;">
            <span class="text-red-600">NIGHT</span>
            <span class="text-slate-700"> STALKERS</span>
          </span>
        </RouterLink>
      </div>
      <div class="hidden md:flex space-x-8 font-semibold text-sm uppercase tracking-wide items-center">
        <RouterLink to="/" class="hover:text-sky-500 transition-colors">Home</RouterLink>
        <RouterLink to="/shop" class="hover:text-sky-500 transition-colors">Toko</RouterLink>
        <!-- Cart Icon with Notification -->
        <RouterLink to="/cart" class="relative group hover:text-sky-500 transition-colors flex items-center">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"></path>
          </svg>
          <span v-if="totalItems > 0" class="absolute -top-2 -right-2 bg-sky-400 text-white text-xs font-bold w-5 h-5 rounded-full flex items-center justify-center">
            {{ totalItems > 99 ? '99+' : totalItems }}
          </span>
        </RouterLink>

        <!-- User Menu -->
        <div v-if="authStore.isAuthenticated" class="flex items-center space-x-4">
          <!-- Profile Icon - Clickable -->
          <RouterLink to="/profile" class="w-8 h-8 rounded-full bg-gradient-to-br from-sky-300 to-sky-500 flex items-center justify-center text-white text-sm font-bold hover:shadow-lg transition-shadow">
            {{ authStore.user?.name?.charAt(0)?.toUpperCase() }}
          </RouterLink>
        </div>

        <!-- Login/Register -->
        <div v-else class="flex items-center space-x-3">
          <RouterLink to="/login" class="hover:text-sky-500 transition-colors">Login</RouterLink>
          <RouterLink to="/register" class="bg-sky-500 hover:bg-sky-600 text-white px-4 py-2 rounded-lg transition-colors">Daftar</RouterLink>
        </div>
      </div>
    </div>
  </nav>
</template>
