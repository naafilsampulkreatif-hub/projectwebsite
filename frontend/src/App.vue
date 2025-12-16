<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'
import Navbar from './components/Navbar.vue'
import Footer from './components/Footer.vue'

onMounted(() => {
  if (!localStorage.getItem('session_id')) {
    const sessionId = 'sess_' + Math.random().toString(36).substr(2, 9) + Date.now().toString(36);
    localStorage.setItem('session_id', sessionId);
  }
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-dark-bg text-white font-sans">
    <Navbar />
    <main class="flex-grow">
      <RouterView />
    </main>
    <Footer />
  </div>
</template>

<style>
/* Global overrides if needed */
body {
  background-color: #1a1a1a;
  color: #ffffff;
}
</style>
