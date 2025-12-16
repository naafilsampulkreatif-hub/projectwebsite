<script setup lang="ts">
import { ref } from 'vue';
import { useCartStore } from '../stores/cart';
import { useRouter } from 'vue-router';

const cartStore = useCartStore();
const router = useRouter();

const form = ref({
    firstName: '',
    lastName: '',
    email: '',
    address: '',
    city: '',
    phone: ''
});

const handleCheckout = async () => {
    // Basic validation
    if (!form.value.firstName || !form.value.email || !form.value.address) {
        alert("Mohon lengkapi data diri anda.");
        return;
    }

    const guestInfo = {
        name: `${form.value.firstName} ${form.value.lastName}`,
        email: form.value.email,
        address: `${form.value.address}, ${form.value.city}`,
        phone: form.value.phone
    };

    const success = await cartStore.checkout(guestInfo);
    if (success) {
        alert("Pesanan berhasil dibuat!");
        router.push('/');
    } else {
        alert("Checkout gagal. Silakan coba lagi.");
    }
};
</script>

<template>
  <div class="min-h-screen bg-[#F3F3F3] text-black">
    <div class="container mx-auto px-4 py-12">
        <h1 class="text-4xl font-extrabold mb-10 border-b-4 border-neon-green inline-block pb-2">Checkout</h1>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12">
             <!-- Billing Details Form -->
             <div class="bg-white p-8 rounded-[2rem] shadow-lg">
                 <h2 class="text-2xl font-bold mb-6">Informasi Pengiriman</h2>
                 <form class="space-y-6" @submit.prevent="handleCheckout">
                     <div class="grid grid-cols-2 gap-6">
                         <div>
                             <label class="block text-sm font-bold mb-2 text-gray-500">Nama Depan</label>
                             <input v-model="form.firstName" type="text" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 focus:outline-none focus:border-neon-green transition-colors" required>
                         </div>
                         <div>
                             <label class="block text-sm font-bold mb-2 text-gray-500">Nama Belakang</label>
                             <input v-model="form.lastName" type="text" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 focus:outline-none focus:border-neon-green transition-colors">
                         </div>
                     </div>
                     <div>
                         <label class="block text-sm font-bold mb-2 text-gray-500">Email</label>
                         <input v-model="form.email" type="email" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 focus:outline-none focus:border-neon-green transition-colors" required>
                     </div>
                     <div>
                         <label class="block text-sm font-bold mb-2 text-gray-500">Alamat</label>
                         <input v-model="form.address" type="text" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 focus:outline-none focus:border-neon-green transition-colors" required>
                     </div>
                     <div class="grid grid-cols-2 gap-6">
                         <div>
                             <label class="block text-sm font-bold mb-2 text-gray-500">Kota</label>
                             <input v-model="form.city" type="text" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 focus:outline-none focus:border-neon-green transition-colors" required>
                         </div>
                         <div>
                             <label class="block text-sm font-bold mb-2 text-gray-500">No. HP</label>
                             <input v-model="form.phone" type="tel" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 focus:outline-none focus:border-neon-green transition-colors" required>
                         </div>
                     </div>
                 </form>
             </div>

             <!-- Order Review -->
             <div class="bg-white p-8 rounded-[2rem] shadow-lg h-fit">
                 <h2 class="text-2xl font-bold mb-6">Ringkasan Pesanan</h2>
                 <div class="space-y-4 mb-8">
                     <div v-for="item in cartStore.items" :key="item.id" class="flex justify-between items-center py-2 border-b border-gray-50">
                         <span class="text-gray-700 font-medium">{{ item.product.name }} <span class="text-xs text-gray-400">x{{ item.quantity }}</span></span>
                         <span class="font-bold">RP {{ (item.product.price * item.quantity).toLocaleString() }}</span>
                     </div>
                 </div>
                 <div class="flex justify-between border-t border-dashed border-gray-300 pt-6 text-xl">
                     <span class="font-bold">Total Pembayaran</span>
                     <span class="font-extrabold text-neon-green">RP {{ cartStore.totalPrice.toLocaleString() }}</span>
                 </div>

                 <div class="mt-8">
                     <button @click="handleCheckout" class="w-full bg-neon-green text-black py-4 rounded-full font-bold uppercase tracking-wider hover:bg-[#00cc00] transition-colors shadow-lg hover:shadow-neon-green/50">
                         Buat Pesanan
                     </button>
                 </div>
             </div>
        </div>
    </div>
  </div>
</template>
