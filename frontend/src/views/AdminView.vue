<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore();
const router = useRouter();
</script>

<template>
	<div class="min-h-screen bg-[#F3F3F3] flex text-black">
		 <!-- Sidebar -->
		 <aside class="w-64 bg-[#1a1a1a] text-white min-h-screen p-8 hidden md:flex flex-col fixed h-full">
				<h2 class="text-3xl font-extrabold mb-6 text-neon-green tracking-tighter">NS ADMIN</h2>
				<nav aria-label="Admin menu" class="w-full">
					<ul class="flex flex-col gap-2 font-bold tracking-wide mt-6">
						<li>
							<router-link to="/admin/produk" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-red-400/20 hover:text-red-700 transition-all transform hover:translate-x-2" active-class="bg-red-600/20 text-red-700">
								<span class="w-2 h-2 bg-neon-green rounded-full"></span>
								<span class="ml-1">Produk</span>
							</router-link>
						</li>
						<li>
							<router-link to="/admin/pesanan" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-red-400/20 hover:text-red-700 transition-all transform hover:translate-x-2" active-class="bg-red-600/20 text-red-700">
								<span class="w-2 h-2 bg-transparent rounded-full"></span>
								<span class="ml-1">Pesanan</span>
							</router-link>
						</li>
						<li>
							<router-link to="/admin/pengguna" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-red-400/20 hover:text-red-700 transition-all transform hover:translate-x-2" active-class="bg-red-600/20 text-red-700">
								<span class="w-2 h-2 bg-transparent rounded-full"></span>
								<span class="ml-1">Pengguna</span>
							</router-link>
						</li>
						<li>
							<router-link to="/admin/pengaturan" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-red-400/20 hover:text-red-700 transition-all transform hover:translate-x-2" active-class="bg-red-600/20 text-red-700">
								<span class="w-2 h-2 bg-transparent rounded-full"></span>
								<span class="ml-1">Pengaturan</span>
							</router-link>
						</li>
					</ul>
				</nav>
				<div class="mt-auto">
					<button @click="authStore.logout(); router.push('/')" class="w-full text-center bg-red-600 text-white px-4 py-2 rounded-lg font-bold uppercase text-sm tracking-widest hover:bg-red-400 transition-all">Logout</button>
				</div>
		 </aside>

		 <!-- Content -->
		 <main class="flex-1 p-8 md:ml-64">
				<RouterView />
		 </main>
	</div>
</template>

<style scoped>
/* small fix to ensure router-link inherits block display */
a.router-link-active {
	background-color: rgba(16,185,129,0.08);
}
</style>
*** End Patch
