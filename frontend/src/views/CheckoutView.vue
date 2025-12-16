<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useCartStore } from '../stores/cart';
import { useRouter } from 'vue-router';
import { useToastStore } from '../stores/toast';

const cartStore = useCartStore();
const router = useRouter();
const toast = useToastStore();

// Steps: 1 = Check Stock, 2 = Form, 3 = Confirmation
const step = ref(1);
const stockLoading = ref(true);
const stockValid = ref(false);
const stockIssues = ref<any[]>([]);

const form = ref({
    firstName: '',
    lastName: '',
    email: '',
    address: '',
    city: '',
    phone: ''
});

const paymentMethod = ref('COD');

// Validate Stock on Mount
const checkStock = async () => {
    stockLoading.value = true;
    stockIssues.value = [];
    try {
        const sessionId = localStorage.getItem('session_id');
        const headers: any = {};
        if (sessionId) headers['X-Session-ID'] = sessionId;

        const res = await fetch(`${import.meta.env.VITE_API_URL}/cart/validate-stock`, { headers });
        const data = await res.json();

        if (data.valid) {
            stockValid.value = true;
            step.value = 2; // Move to Form
        } else {
            stockValid.value = false;
            stockIssues.value = data.issues;
            // Stay on Step 1 (Error View)
        }
    } catch (e) {
        toast.show("Gagal mengecek stok", "error");
    } finally {
        stockLoading.value = false;
    }
};

onMounted(() => {
    if (cartStore.items.length === 0) {
        router.push('/cart');
        return;
    }
    checkStock();
});

const goToConfirmation = () => {
    if (!form.value.firstName || !form.value.email || !form.value.address) {
        toast.show("Mohon lengkapi data diri.", "error");
        return;
    }
    step.value = 3;
};

const handleCheckout = async () => {
    const guestInfo = {
        name: `${form.value.firstName} ${form.value.lastName}`,
        email: form.value.email,
        address: `${form.value.address}, ${form.value.city}`,
        phone: form.value.phone
    };

    // Final attempt
    try {
        const data = await cartStore.checkout(guestInfo);
        if (data && data.order_id) {
            toast.show("Pesanan berhasil dibuat!", "success");
            router.push(`/invoice/${data.order_id}`);
        }
    } catch (e: any) {
        // If 409 Conflict (Stock), go back to step 1
        if (e.response && e.response.status === 409) {
            toast.show(e.response.data || "Stok tidak mencukupi.", "error");
            step.value = 1;
            checkStock(); // Re-validate
        } else {
            toast.show("Checkout gagal.", "error");
        }
    }
};
</script>

<template>
  <div class="min-h-screen bg-[#F3F3F3] text-black">
    <div class="container mx-auto px-4 py-12">
        <h1 class="text-4xl font-extrabold mb-10 border-b-4 border-neon-green inline-block pb-2">Checkout</h1>

        <!-- Step 1: Stock Validation Loading/Error -->
        <div v-if="step === 1" class="bg-white p-8 rounded-[2rem] shadow-lg max-w-2xl mx-auto text-center">
            <div v-if="stockLoading">
                <p class="text-xl font-bold animate-pulse">Mengecek ketersediaan stok...</p>
            </div>
            <div v-else-if="!stockValid">
                <h2 class="text-2xl font-bold text-red-600 mb-4">Stok Tidak Mencukupi</h2>
                <div class="bg-red-50 p-4 rounded-xl mb-6 text-left">
                    <p class="font-medium mb-2">Item berikut tidak tersedia dalam jumlah yang diminta:</p>
                    <ul class="list-disc list-inside space-y-1">
                        <li v-for="issue in stockIssues" :key="issue.product_name">
                            <span class="font-bold">{{ issue.product_name }}</span>
                            (Diminta: {{ issue.requested }}, Tersedia: {{ issue.available }})
                        </li>
                    </ul>
                </div>
                <div class="flex gap-4 justify-center">
                    <router-link to="/cart" class="bg-gray-200 px-6 py-3 rounded-full font-bold">Kembali ke Keranjang</router-link>
                    <button @click="checkStock" class="bg-neon-green px-6 py-3 rounded-full font-bold">Coba Lagi</button>
                </div>
            </div>
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-12">

             <!-- Step 2: Form Data -->
             <div v-if="step === 2" class="bg-white p-8 rounded-[2rem] shadow-lg">
                 <h2 class="text-2xl font-bold mb-6 flex items-center gap-2">
                    <span class="bg-black text-white w-8 h-8 rounded-full flex items-center justify-center text-sm">1</span>
                    Informasi Pengiriman
                 </h2>
                 <form class="space-y-6" @submit.prevent="goToConfirmation">
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

                     <!-- Payment Method Section inside form for Step 2 flow or separate? -->
                     <!-- Flowchart says: Form -> Validate -> Payment -> Confirm -->

                     <div class="pt-6 border-t border-gray-100">
                        <h3 class="text-xl font-bold mb-4">Metode Pembayaran</h3>
                        <label class="flex items-center gap-3 p-4 border border-neon-green bg-green-50 rounded-xl cursor-pointer">
                            <input type="radio" v-model="paymentMethod" value="COD" checked class="accent-neon-green w-5 h-5">
                            <div>
                                <span class="font-bold block">COD (Cash on Delivery)</span>
                                <span class="text-sm text-gray-500">Bayar ditempat saat kurir tiba.</span>
                            </div>
                        </label>
                     </div>

                     <button type="submit" class="w-full bg-black text-white py-4 rounded-full font-bold hover:bg-gray-800 transition-colors mt-6">
                         Lanjut ke Konfirmasi
                     </button>
                 </form>
             </div>

             <!-- Step 3: Confirmation -->
             <div v-if="step === 3" class="bg-white p-8 rounded-[2rem] shadow-lg">
                <h2 class="text-2xl font-bold mb-6 flex items-center gap-2">
                    <span class="bg-black text-white w-8 h-8 rounded-full flex items-center justify-center text-sm">2</span>
                    Konfirmasi Pesanan
                 </h2>

                 <div class="space-y-6 mb-8">
                     <div class="bg-gray-50 p-4 rounded-xl">
                         <h3 class="font-bold text-gray-500 text-sm uppercase mb-2">Dikirim ke:</h3>
                         <p class="font-bold">{{ form.firstName }} {{ form.lastName }}</p>
                         <p>{{ form.address }}, {{ form.city }}</p>
                         <p>{{ form.phone }}</p>
                         <button @click="step = 2" class="text-sm text-neon-green font-bold mt-2 underline">Ubah Data</button>
                     </div>

                     <div class="bg-gray-50 p-4 rounded-xl">
                        <h3 class="font-bold text-gray-500 text-sm uppercase mb-2">Pembayaran:</h3>
                        <p class="font-bold">{{ paymentMethod }}</p>
                    </div>
                 </div>

                 <button @click="handleCheckout" class="w-full bg-neon-green text-black py-4 rounded-full font-bold uppercase tracking-wider hover:bg-[#00cc00] transition-colors shadow-lg hover:shadow-neon-green/50">
                     Konfirmasi & Proses Pesanan
                 </button>
             </div>

             <!-- Order Summary (Always Visible) -->
             <div class="bg-white p-8 rounded-[2rem] shadow-lg h-fit">
                 <h2 class="text-2xl font-bold mb-6">Ringkasan</h2>
                 <div class="space-y-4 mb-8">
                     <div v-for="item in cartStore.items" :key="item.id" class="flex justify-between items-center py-2 border-b border-gray-50">
                         <span class="text-gray-700 font-medium">{{ item.product.name }} <span class="text-xs text-gray-400">x{{ item.quantity }}</span></span>
                         <span class="font-bold">RP {{ (item.product.price * item.quantity).toLocaleString() }}</span>
                     </div>
                 </div>
                 <div class="flex justify-between border-t border-dashed border-gray-300 pt-6 text-xl">
                     <span class="font-bold">Total</span>
                     <span class="font-extrabold text-neon-green">RP {{ cartStore.totalPrice.toLocaleString() }}</span>
                 </div>
             </div>
        </div>
    </div>
  </div>
</template>
