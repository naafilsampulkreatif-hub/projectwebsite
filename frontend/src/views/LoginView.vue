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
  try {
    const success = await authStore.login({ email: email.value, password: password.value });
    if (success) {
      // Reset form on success
      email.value = '';
      password.value = '';
      if (authStore.isAdmin) {
        router.push('/admin');
      } else {
        router.push('/');
      }
      return;
    }
    error.value = 'Email atau Password salah.';
  } catch (err) {
    console.error('Login error', err);
    error.value = 'Terjadi kesalahan saat login.';
  } finally {
    // Clear password in all cases so it's not left in the field
    password.value = '';
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center py-6 sm:py-12" style="background-color: #F9DFDF;">
    <div class="w-full max-w-md px-3 sm:px-4">
      <div class="space-y-4 sm:space-y-6 md:space-y-8">
        <!-- Header -->
        <div class="text-center">
          <h1 class="text-3xl sm:text-4xl font-extrabold bg-gradient-to-r from-sky-600 to-sky-500 bg-clip-text text-transparent mb-1 sm:mb-2">
            Masuk
          </h1>
          <p class="text-slate-800 font-semibold text-sm sm:text-base">Selamat datang kembali ke Night Stalkers</p>
        </div>

        <!-- Login Card -->
        <div class="rounded-2xl sm:rounded-3xl shadow-xl border border-slate-200 p-5 sm:p-6 md:p-8" style="background-color: #FCF8F8;">
          <form @submit.prevent="handleLogin" class="space-y-3 sm:space-y-4 md:space-y-6">
            <!-- Email Input -->
            <div class="space-y-1 sm:space-y-2">
              <label class="block text-xs sm:text-sm font-bold text-slate-900">Email Address</label>
              <input 
                v-model="email" 
                type="email" 
                placeholder="user@example.com"
                class="w-full border-2 border-slate-200 rounded-lg sm:rounded-xl p-2 sm:p-3 md:p-4 focus:outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-200 transition-all duration-200 text-slate-900 placeholder:text-slate-700 font-medium text-xs sm:text-sm"
                style="background-color: #FBEFEF;"
                required
              />
            </div>

            <!-- Password Input -->
            <div class="space-y-1 sm:space-y-2">
              <label class="block text-xs sm:text-sm font-bold text-slate-900">Password</label>
              <input 
                v-model="password" 
                type="password" 
                placeholder="••••••••"
                class="w-full border-2 border-slate-200 rounded-lg sm:rounded-xl p-2 sm:p-3 md:p-4 focus:outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-200 transition-all duration-200 text-slate-900 placeholder:text-slate-700 font-medium text-xs sm:text-sm"
                style="background-color: #FBEFEF;"
                required
              />
            </div>

            <!-- Error Message -->
            <div v-if="error" class="bg-red-50 border-2 border-red-200 text-red-700 text-xs sm:text-sm font-bold p-2 sm:p-3 md:p-4 rounded-lg sm:rounded-xl">
              {{ error }}
            </div>

            <!-- Login Button -->
            <button 
              type="submit" 
              class="w-full bg-gradient-to-r from-sky-500 to-sky-600 hover:from-sky-600 hover:to-sky-700 text-white py-2 sm:py-3 md:py-4 rounded-lg sm:rounded-xl font-bold text-sm sm:text-base uppercase tracking-wider transition-all duration-200 shadow-lg hover:shadow-sky-300/50"
            >
              Masuk
            </button>

            <!-- Divider -->
            <div class="relative py-2 sm:py-3">
              <div class="absolute inset-0 flex items-center">
                <div class="w-full border-t border-slate-200"></div>
              </div>
              <div class="relative flex justify-center text-xs sm:text-sm">
                <span class="px-2 text-slate-600" style="background-color: #FCF8F8;">atau</span>
              </div>
            </div>

            <!-- Register Link Button -->
            <router-link 
              to="/register"
              class="block w-full text-center bg-slate-100 hover:bg-slate-200 text-slate-900 py-2 sm:py-3 md:py-4 rounded-lg sm:rounded-xl font-bold text-sm sm:text-base uppercase tracking-wider transition-all duration-200"
            >
              Buat Akun Baru
            </router-link>
          </form>

          <!-- Footer Link -->
          <div class="mt-4 sm:mt-5 md:mt-8 pt-3 sm:pt-4 md:pt-6 border-t border-slate-200 text-center text-xs sm:text-sm text-slate-800">
            <p>Dengan masuk, Anda menyetujui <span class="font-bold text-slate-900">Syarat & Ketentuan</span> kami</p>
          </div>
        </div>

        <!-- Benefits -->
        <div class="grid grid-cols-3 gap-2 sm:gap-3 md:gap-4">
          <div class="text-center p-2 sm:p-3">
            <div class="bg-sky-100 rounded-full w-10 sm:w-12 md:w-14 h-10 sm:h-12 md:h-14 flex items-center justify-center mx-auto mb-1 sm:mb-2">
              <svg class="w-5 sm:w-6 md:w-7 h-5 sm:h-6 md:h-7 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
            </div>
            <p class="text-xs font-bold text-slate-900">Checkout Mudah</p>
          </div>
          <div class="text-center p-2 sm:p-3">
            <div class="bg-sky-100 rounded-full w-10 sm:w-12 md:w-14 h-10 sm:h-12 md:h-14 flex items-center justify-center mx-auto mb-1 sm:mb-2">
              <svg class="w-5 sm:w-6 md:w-7 h-5 sm:h-6 md:h-7 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
              </svg>
            </div>
            <p class="text-xs font-bold text-slate-900">Akses Cepat</p>
          </div>
          <div class="text-center p-2 sm:p-3">
            <div class="bg-sky-100 rounded-full w-10 sm:w-12 md:w-14 h-10 sm:h-12 md:h-14 flex items-center justify-center mx-auto mb-1 sm:mb-2">
              <svg class="w-5 sm:w-6 md:w-7 h-5 sm:h-6 md:h-7 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
              </svg>
            </div>
            <p class="text-xs font-bold text-slate-900">Aman</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
