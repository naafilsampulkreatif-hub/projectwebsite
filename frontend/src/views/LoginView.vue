<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import Header from '../components/Header.vue';

const email = ref('');
const password = ref('');
const error = ref('');
const authStore = useAuthStore();
const router = useRouter();

const handleLogin = async () => {
  error.value = '';
  const success = await authStore.login({ email: email.value, password: password.value });
  if (success) {
    router.push('/');
  } else {
    error.value = 'Invalid credentials';
  }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <Header />
    <div class="container mx-auto px-4 py-16 flex justify-center">
        <div class="bg-white p-8 rounded shadow-md w-full max-w-md">
            <h1 class="text-2xl font-bold mb-6 text-center">Login</h1>

            <form @submit.prevent="handleLogin" class="space-y-4">
                <div>
                    <label class="block text-gray-700 mb-2">Email address</label>
                    <input v-model="email" type="email" class="w-full border p-2 rounded focus:border-primary outline-none" required>
                </div>
                <div>
                    <label class="block text-gray-700 mb-2">Password</label>
                    <input v-model="password" type="password" class="w-full border p-2 rounded focus:border-primary outline-none" required>
                </div>

                <div v-if="error" class="text-red-500 text-sm">{{ error }}</div>

                <button type="submit" class="w-full bg-primary text-white py-2 rounded hover:bg-orange-600 transition-colors">
                    Log In
                </button>
            </form>

            <p class="mt-4 text-center text-sm text-gray-600">
                Don't have an account? <router-link to="/register" class="text-primary hover:underline">Register</router-link>
            </p>
        </div>
    </div>
  </div>
</template>
