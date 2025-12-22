<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import { onMounted, computed } from 'vue'
import Navbar from './components/Navbar.vue'
import Footer from './components/Footer.vue'
import Toast from './components/Toast.vue'
import CustomerService from './components/CustomerService.vue'

const route = useRoute();

const showLayout = computed(() => {
  // Hide Navbar and Footer on admin routes (path starts with /admin)
  return !route.path.startsWith('/admin');
});

onMounted(() => {
  if (!localStorage.getItem('session_id')) {
    const sessionId = 'sess_' + Math.random().toString(36).substr(2, 9) + Date.now().toString(36);
    localStorage.setItem('session_id', sessionId);
  }
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-dark-bg text-white font-sans">
    <Navbar v-if="showLayout" />
    <main class="flex-grow">
      <RouterView />
    </main>
    <Footer v-if="showLayout" />
    <Toast />
    <CustomerService />
  </div>
</template>

<style>
/* Global overrides if needed */
body {
  background-color: #1a1a1a;
  color: #ffffff;
}
</style>
