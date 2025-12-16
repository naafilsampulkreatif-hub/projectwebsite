<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import Header from '../components/Header.vue';

const name = ref('');
const email = ref('');
const password = ref('');
const error = ref('');
const authStore = useAuthStore();
const router = useRouter();

const handleRegister = async () => {
  error.value = '';
  const success = await authStore.register({ name: name.value, email: email.value, password: password.value });
  if (success) {
    router.push('/login');
  } else {
    error.value = 'Registration failed';
  }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <Header />
    <div class="container mx-auto px-4 py-16 flex justify-center">
        <div class="bg-white p-8 rounded shadow-md w-full max-w-md">
            <h1 class="text-2xl font-bold mb-6 text-center">Register</h1>

            <form @submit.prevent="handleRegister" class="space-y-4">
                <div>
                    <label class="block text-gray-700 mb-2">Full Name</label>
                    <input v-model="name" type="text" class="w-full border p-2 rounded focus:border-primary outline-none" required>
                </div>
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
                    Sign Up
                </button>
            </form>

            <p class="mt-4 text-center text-sm text-gray-600">
                Already have an account? <router-link to="/login" class="text-primary hover:underline">Login</router-link>
            </p>
        </div>
    </div>
  </div>
</template>
