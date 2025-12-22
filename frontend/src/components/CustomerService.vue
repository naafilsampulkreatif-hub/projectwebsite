<script setup lang="ts">
import { ref, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import api from '../services/api';

const authStore = useAuthStore();

const messages = ref<any[]>([]);
const newMessage = ref('');
const loading = ref(false);
const isOpen = ref(false);
const unreadCount = ref(0);

const canSendMessage = computed(() => authStore.isAuthenticated);

const toggleChat = async () => {
  if (!canSendMessage.value) return;
  isOpen.value = !isOpen.value;
  if (isOpen.value) {
    unreadCount.value = 0;
  }
};

const sendMessage = async () => {
  if (!newMessage.value.trim() || loading.value) return;

  loading.value = true;
  try {
    const response = await api.post('/support/message', {
      message: newMessage.value.trim(),
      type: 'user'
    });
    
    if (response.data) {
      messages.value.push({
        id: Date.now(),
        message: newMessage.value.trim(),
        type: 'user',
        created_at: new Date().toISOString()
      });
      newMessage.value = '';

      // Simulate support response (in real app, this would be WebSocket or polling)
      setTimeout(() => {
        messages.value.push({
          id: Date.now() + 1,
          message: 'Terima kasih telah menghubungi kami. Tim support akan membalas Anda segera.',
          type: 'support',
          created_at: new Date().toISOString()
        });
      }, 1000);
    }
  } catch (error: any) {
    console.error('Failed to send message:', error);
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <!-- Chat Button -->
  <div v-if="canSendMessage" class="fixed bottom-6 right-6 z-40">
    <button 
      v-if="!isOpen" 
      @click="toggleChat"
      class="text-white rounded-full p-4 shadow-lg hover:shadow-xl transition-all duration-200 w-16 h-16 flex items-center justify-center hover:opacity-80"
      style="background-color: #8BAE66;"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
      </svg>
      <span v-if="unreadCount > 0" class="absolute -top-1 -right-1 bg-sky-400 text-white text-xs font-bold w-5 h-5 rounded-full flex items-center justify-center">
        {{ unreadCount }}
      </span>
    </button>

    <!-- Chat Window -->
    <div v-else class="fixed bottom-24 right-6 w-96 max-w-[calc(100vw-24px)] bg-white rounded-2xl shadow-2xl border border-slate-200 flex flex-col h-96 z-40">
      <!-- Header -->
      <div class="bg-gradient-to-r from-emerald-500 to-emerald-600 text-white p-4 rounded-t-2xl flex justify-between items-center">
        <div>
          <h3 class="font-bold text-lg">Customer Service</h3>
          <p class="text-emerald-100 text-xs">Tim support kami siap membantu</p>
        </div>
        <button @click="toggleChat" class="text-white hover:bg-sky-700 p-1 rounded transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>

      <!-- Messages -->
      <div class="flex-1 overflow-y-auto p-4 space-y-4">
        <div v-if="messages.length === 0" class="text-center text-slate-500 py-8">
          <p class="text-sm">Halo! Bagaimana kami bisa membantu Anda?</p>
        </div>

        <div v-for="msg in messages" :key="msg.id" :class="['flex', msg.type === 'user' ? 'justify-end' : 'justify-start']">
          <div :class="['max-w-xs rounded-lg p-3 text-sm', msg.type === 'user' ? 'bg-sky-500 text-white rounded-br-none' : 'bg-slate-100 text-slate-800 rounded-bl-none']">
            {{ msg.message }}
          </div>
        </div>
      </div>

      <!-- Input -->
      <div class="border-t border-slate-200 p-4">
        <form @submit.prevent="sendMessage" class="flex gap-2">
          <input 
            v-model="newMessage"
            type="text"
            placeholder="Ketik pesan..."
            :disabled="loading"
            class="flex-1 bg-slate-50 border border-slate-200 rounded-lg p-2 text-sm text-black focus:outline-none focus:border-emerald-400 transition-colors disabled:opacity-50"
          />
          <button 
            type="submit"
            :disabled="!newMessage.trim() || loading"
            class="text-white rounded-lg px-3 py-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed hover:opacity-80"
            style="background-color: #547792;"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path>
            </svg>
          </button>
        </form>
      </div>
    </div>
  </div>

</template>
