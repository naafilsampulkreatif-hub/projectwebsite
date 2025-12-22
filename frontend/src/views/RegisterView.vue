<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

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
    alert("Registrasi berhasil! Silakan login.");
    router.push('/login');
  } else {
    error.value = 'Registrasi gagal. Email mungkin sudah terdaftar.';
  }
};
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-sky-50 to-slate-50 flex items-center justify-center py-12">
    <div class="w-full max-w-md px-4">
      <div class="space-y-8">
        <!-- Header -->
        <div class="text-center">
          <h1 class="text-4xl font-extrabold bg-gradient-to-r from-sky-600 to-sky-500 bg-clip-text text-transparent mb-2">
            Daftar
          </h1>
          <p class="text-slate-800 font-semibold">Bergabunglah dengan Night Stalkers hari ini</p>
        </div>

        <!-- Register Card -->
        <div class="bg-white rounded-3xl shadow-2xl border border-slate-200 p-8">
          <form @submit.prevent="handleRegister" class="space-y-6">
            <!-- Name Input -->
            <div class="space-y-2">
              <label class="block text-sm font-bold text-slate-900">Nama Lengkap</label>
              <input 
                v-model="name" 
                type="text" 
                placeholder="John Doe"
                class="w-full bg-gradient-to-b from-slate-50 to-slate-100 border-2 border-slate-200 rounded-xl p-4 focus:outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-200 transition-all duration-200 text-slate-900 placeholder:text-slate-700 font-medium"
                required
              />
            </div>

            <!-- Email Input -->
            <div class="space-y-2">
              <label class="block text-sm font-bold text-slate-900">Email Address</label>
              <input 
                v-model="email" 
                type="email" 
                placeholder="user@example.com"
                class="w-full bg-gradient-to-b from-slate-50 to-slate-100 border-2 border-slate-200 rounded-xl p-4 focus:outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-200 transition-all duration-200 text-slate-900 placeholder:text-slate-700 font-medium"
                required
              />
            </div>

            <!-- Password Input -->
            <div class="space-y-2">
              <label class="block text-sm font-bold text-slate-900">Password</label>
              <input 
                v-model="password" 
                type="password" 
                placeholder="••••••••"
                class="w-full bg-gradient-to-b from-slate-50 to-slate-100 border-2 border-slate-200 rounded-xl p-4 focus:outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-200 transition-all duration-200 text-slate-900 placeholder:text-slate-700 font-medium"
                required
              />
            </div>

            <!-- Error Message -->
            <div v-if="error" class="bg-red-50 border-2 border-red-200 text-red-700 text-sm font-bold p-4 rounded-xl">
              {{ error }}
            </div>

            <!-- Register Button -->
            <button 
              type="submit" 
              class="w-full bg-gradient-to-r from-sky-500 to-sky-600 hover:from-sky-600 hover:to-sky-700 text-white py-4 rounded-xl font-bold uppercase tracking-wider transition-all duration-200 shadow-lg hover:shadow-sky-300/50 transform hover:scale-105"
            >
              Daftar
            </button>

            <!-- Divider -->
            <div class="relative py-4">
              <div class="absolute inset-0 flex items-center">
                <div class="w-full border-t border-slate-200"></div>
              </div>
              <div class="relative flex justify-center text-sm">
                <span class="px-2 bg-white text-slate-500">atau</span>
              </div>
            </div>

            <!-- Login Link Button -->
            <router-link 
              to="/login"
              class="block w-full text-center bg-slate-100 hover:bg-slate-200 text-slate-700 py-4 rounded-xl font-bold uppercase tracking-wider transition-all duration-200"
            >
              Sudah Punya Akun?
            </router-link>
          </form>

          <!-- Footer Link -->
          <div class="mt-8 pt-6 border-t border-slate-200 text-center text-sm text-slate-800">
            <p>Dengan mendaftar, Anda menyetujui <span class="font-bold text-slate-900">Syarat & Ketentuan</span> kami</p>
          </div>
        </div>

        <!-- Features -->
        <div class="grid grid-cols-3 gap-4">
          <div class="text-center">
            <div class="bg-sky-100 rounded-full w-12 h-12 flex items-center justify-center mx-auto mb-2">
              <svg class="w-6 h-6 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 10h4.764a2 2 0 011.789 2.894l-3.647 7.292a2 2 0 01-1.788 1.106H7m0 0a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m0 0a2 2 0 110-4m0 4a2 2 0 100-4"></path>
              </svg>
            </div>
            <p class="text-xs font-bold text-slate-900">Gratis Ongkir</p>
          </div>
          <div class="text-center">
            <div class="bg-sky-100 rounded-full w-12 h-12 flex items-center justify-center mx-auto mb-2">
              <svg class="w-6 h-6 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192l-3.536 3.536M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z"></path>
              </svg>
            </div>
            <p class="text-xs font-bold text-slate-900">Rewards</p>
          </div>
          <div class="text-center">
            <div class="bg-sky-100 rounded-full w-12 h-12 flex items-center justify-center mx-auto mb-2">
              <svg class="w-6 h-6 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
              </svg>
            </div>
            <p class="text-xs font-bold text-slate-900">Support</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
