<script setup lang="ts">
import { useToastStore } from '../stores/toast';

const toast = useToastStore();
</script>

<template>
  <div v-if="toast.visible" class="fixed bottom-4 right-4 px-6 py-4 rounded-lg shadow-lg text-white font-bold transition-opacity z-50"
       :class="toast.type === 'success' ? 'bg-neon-green text-black' : 'bg-red-500'">
      {{ toast.message }}
  </div>
</template>
