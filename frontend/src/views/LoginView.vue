<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const email = ref('');
const password = ref('');
const error = ref('');
const authStore = useAuthStore();
const router = useRouter();

const handleLogin = async () => {
  error.value = '';
  const success = await authStore.login({ email: email.value, password: password.value });
  if (success) {
    if (authStore.isAdmin) {
      router.push('/admin');
    } else {
      router.push('/');
    }
  } else {
    error.value = 'Email atau Password salah.';
  }
};
</script>

<template>
  <div class="min-h-screen bg-[#F3F3F3] text-black flex items-center justify-center">
    <div class="w-full max-w-md px-4">
        <div class="bg-white p-10 rounded-[2rem] shadow-xl">
            <h1 class="text-3xl font-extrabold mb-8 text-center text-neon-green">LOGIN</h1>

            <form @submit.prevent="handleLogin" class="space-y-6">
                <div>
                    <label class="block text-sm font-bold mb-2 text-gray-500">Email Address</label>
                    <input v-model="email" type="email" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 focus:outline-none focus:border-neon-green transition-colors" placeholder="user@example.com" required>
                </div>
                <div>
                    <label class="block text-sm font-bold mb-2 text-gray-500">Password</label>
                    <input v-model="password" type="password" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 focus:outline-none focus:border-neon-green transition-colors" placeholder="••••••••" required>
                </div>

                <div v-if="error" class="text-red-500 text-sm font-bold text-center">{{ error }}</div>

                <button type="submit" class="w-full bg-neon-green text-black py-3 rounded-full font-bold uppercase tracking-wider hover:bg-[#00cc00] transition-colors shadow-lg hover:shadow-neon-green/50">
                    Masuk
                </button>
            </form>

             <p class="mt-2 text-center text-sm text-gray-400">
                <router-link to="/" class="hover:underline">Kembali ke Home</router-link>
            </p>
        </div>
    </div>
  </div>
</template>
