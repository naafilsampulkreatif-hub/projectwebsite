<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useCartStore } from '../stores/cart';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import { useToastStore } from '../stores/toast';
import api from '../services/api';
import { getProvincesArray, getCitiesForProvince } from '../data/indonesianRegions';

const cartStore = useCartStore();
const authStore = useAuthStore();
const router = useRouter();
const toast = useToastStore();

// Pastikan session_id selalu ada di localStorage
if (!localStorage.getItem('session_id')) {
    localStorage.setItem('session_id', 'session-' + Math.random().toString(36).substring(2, 12));
}

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
    province: '',
    city: '',
    zipCode: '',
    phone: ''
});

const provinces = ref<string[]>([]);
const citiesForProvince = computed(() => {
    return form.value.province ? getCitiesForProvince(form.value.province) : [];
});

const selectedShippingCost = computed(() => {
    const method = shippingMethods.value.find(m => m.id === selectedShippingMethod.value);
    return method?.cost || 0;
});

const totalWithShipping = computed(() => {
    return cartStore.totalPrice + selectedShippingCost.value;
});

const paymentMethod = ref('COD');
const shippingMethods = ref<any[]>([]);
const selectedShippingMethod = ref<number>(1);
const loadingShipping = ref(true);

// Validate Stock on Mount
const checkStock = async () => {
    stockLoading.value = true;
    stockIssues.value = [];
    try {
        const res = await api.get('/cart/validate-stock');
        const data = res.data;

        if (data.valid) {
            stockValid.value = true;
            step.value = 2; // Move to Form
        } else {
            stockValid.value = false;
            stockIssues.value = data.issues;
            // Stay on Step 1 (Error View)
        }
    } catch (e: any) {
        console.error('Validate stock error', e);
        toast.show("Gagal mengecek stok", "error");
    } finally {
        stockLoading.value = false;
    }
};

onMounted(async () => {
    // Ensure cart is loaded before validating
    await cartStore.fetchCart();
    if (cartStore.items.length === 0) {
        router.push('/cart');
        return;
    }
    provinces.value = getProvincesArray();
    
    // Auto-fill data dari user profile jika sudah login
    if (authStore.isAuthenticated && authStore.user) {
        form.value.firstName = authStore.user.name?.split(' ')[0] || '';
        form.value.lastName = authStore.user.name?.split(' ').slice(1).join(' ') || '';
        form.value.email = authStore.user.email || '';
        form.value.phone = authStore.user.phone || '';
        form.value.address = authStore.user.address || '';
        form.value.province = authStore.user.province || '';
        form.value.city = authStore.user.city || '';
        form.value.zipCode = authStore.user.postal_code || '';
    }
    
    // Fetch shipping methods
    try {
        const res = await api.get('/shipping-methods');
        shippingMethods.value = res.data || [];
        if (shippingMethods.value.length > 0) {
            selectedShippingMethod.value = shippingMethods.value[0].id;
        }
    } catch (e) {
        console.error('Failed to load shipping methods:', e);
        toast.show('Gagal memuat metode pengiriman', 'error');
    } finally {
        loadingShipping.value = false;
    }
    
    checkStock();
});

const goToConfirmation = () => {
    if (!form.value.firstName || !form.value.email || !form.value.address || !form.value.province || !form.value.city || !form.value.zipCode || !form.value.phone) {
        toast.show("Mohon lengkapi semua data diri.", "error");
        return;
    }
    step.value = 3;
};

const handleCheckout = async () => {
    const guestInfo = {
        full_name: `${form.value.firstName} ${form.value.lastName}`,
        email: form.value.email,
        address: form.value.address,
        province: form.value.province,
        city: form.value.city,
        postal_code: form.value.zipCode,
        phone: form.value.phone,
        shipping_method_id: selectedShippingMethod.value
    };

    try {
        console.log('Starting checkout with guest info:', guestInfo);
        const data = await cartStore.checkout(guestInfo);
        console.log('Checkout response:', data);
        
        if (data && data.order_id) {
            toast.show("Pesanan berhasil dibuat!", "success");
            console.log('Order created:', data.order_id);
            // Coba redirect ke invoice, jika gagal redirect ke halaman pembangunan
            try {
                await router.push(`/invoice/${data.order_id}`);
            } catch (err) {
                console.error('Failed to redirect to invoice:', err);
                router.push('/under-construction');
            }
        } else {
            console.warn('No order_id in response:', data);
            toast.show("Error: Pesanan tidak mendapat ID", "error");
            router.push('/under-construction');
        }
    } catch (e: any) {
        console.error('Checkout error:', e.response?.data || e.message);
        const errorMsg = e.response?.data?.message || e.message || 'Gagal memproses pesanan';
        toast.show(`Error: ${errorMsg}`, 'error');
        // Redirect ke halaman pembangunan setelah error
        setTimeout(() => {
            router.push('/under-construction');
        }, 2000);
    }
};
</script>

<template>
  <div class="min-h-screen bg-white text-black">
    <div class="container mx-auto px-4 py-6 md:py-12">
<h1 class="text-3xl md:text-4xl font-extrabold mb-6 md:mb-10 border-b-4 border-blue-500 inline-block pb-2">Checkout</h1>

        <!-- Step 1: Stock Validation Loading/Error -->
        <div v-if="step === 1" class="bg-white p-4 md:p-8 rounded-[2rem] shadow-lg max-w-2xl mx-auto text-center">
            <div v-if="stockLoading">
                <p class="text-lg font-bold animate-pulse">Mengecek ketersediaan stok...</p>
            </div>
            <div v-else-if="!stockValid">
                <h2 class="text-xl md:text-2xl font-bold text-red-600 mb-4">Stok Tidak Mencukupi</h2>
                <div class="bg-red-50 p-3 md:p-4 rounded-xl mb-6 text-left text-sm md:text-base">
                    <p class="font-medium mb-2">Item berikut tidak tersedia dalam jumlah yang diminta:</p>
                    <ul class="list-disc list-inside space-y-1">
                        <li v-for="issue in stockIssues" :key="issue.product_name">
                            <span class="font-bold">{{ issue.product_name }}</span>
                            (Diminta: {{ issue.requested }}, Tersedia: {{ issue.available }})
                        </li>
                    </ul>
                </div>
                <div class="flex gap-4 justify-center">
                    <router-link to="/cart" class="bg-slate-200 hover:bg-slate-300 px-6 py-3 rounded-full font-bold text-slate-800">Kembali ke Keranjang</router-link>
                    <button @click="checkStock" class="text-white px-6 py-3 rounded-full font-bold transition-colors hover:opacity-80" style="background-color: #547792;">Coba Lagi</button>
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
                             <label class="block text-sm font-bold mb-2 text-gray-500">Nama Depan <span class="text-blue-500">*</span></label>
                             <input v-model="form.firstName" type="text" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 text-black focus:outline-none focus:border-blue-400 transition-colors" required>
                         </div>
                         <div>
                             <label class="block text-sm font-bold mb-2 text-gray-500">Nama Belakang</label>
                             <input v-model="form.lastName" type="text" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 text-black focus:outline-none focus:border-blue-400 transition-colors">
                         </div>
                     </div>
                     <div>
                         <label class="block text-sm font-bold mb-2 text-gray-500">Email <span class="text-blue-500">*</span></label>
                         <input v-model="form.email" type="email" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 text-black focus:outline-none focus:border-blue-400 transition-colors" required>
                     </div>
                     <div>
                         <label class="block text-sm font-bold mb-2 text-gray-500">Alamat <span class="text-blue-500">*</span></label>
                         <input v-model="form.address" type="text" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 text-black focus:outline-none focus:border-blue-400 transition-colors" required>
                     </div>
                     <div>
                         <label class="block text-sm font-bold mb-2 text-gray-500">Provinsi <span class="text-blue-500">*</span></label>
                         <select v-model="form.province" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 text-black focus:outline-none focus:border-blue-400 transition-colors" required>
                             <option value="">Pilih Provinsi</option>
                             <option v-for="prov in provinces" :key="prov" :value="prov">{{ prov }}</option>
                         </select>
                     </div>
                     <div v-if="form.province">
                         <label class="block text-sm font-bold mb-2 text-gray-500">Kota/Kabupaten <span class="text-red-500">*</span></label>
                         <select v-model="form.city" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 text-black focus:outline-none focus:border-blue-400 transition-colors" required>
                             <option value="">Pilih Kota</option>
                             <option v-for="city in citiesForProvince" :key="city" :value="city">{{ city }}</option>
                         </select>
                     </div>
                     <div class="grid grid-cols-2 gap-6">
                        <div>
                             <label class="block text-sm font-bold mb-2 text-gray-500">Kode Pos <span class="text-red-500">*</span></label>
                             <input v-model="form.zipCode" type="text" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 text-black focus:outline-none focus:border-blue-400 transition-colors" required>
                         </div>
                         <div>
                             <label class="block text-sm font-bold mb-2 text-gray-500">No. HP <span class="text-red-500">*</span></label>
                             <input v-model="form.phone" type="tel" class="w-full bg-gray-50 border border-gray-200 rounded-xl p-3 text-black focus:outline-none focus:border-blue-400 transition-colors" required>
                         </div>
                     </div>

                     <!-- Payment Method Section inside form for Step 2 flow or separate? -->
                     <!-- Flowchart says: Form -> Validate -> Payment -> Confirm -->

                     <div class="pt-6 border-t border-gray-100">
                        <h3 class="text-xl font-bold mb-4">Metode Pengiriman</h3>
                        <div v-if="loadingShipping" class="text-center py-4">
                            <p class="text-gray-500">Memuat metode pengiriman...</p>
                        </div>
                        <div v-else-if="shippingMethods.length === 0" class="text-center py-4">
                            <p class="text-gray-500">Tidak ada metode pengiriman tersedia</p>
                        </div>
                        <div v-else class="space-y-3">
                            <label v-for="method in shippingMethods" :key="method.id" class="flex items-center gap-3 p-4 border border-gray-300 rounded-xl cursor-pointer hover:border-blue-400 hover:bg-blue-50 transition-all"
                                   :class="selectedShippingMethod === method.id ? 'border-blue-400 bg-blue-50' : ''">
                                <input type="radio" :value="method.id" v-model="selectedShippingMethod" class="accent-blue-500 w-5 h-5">
                                <div class="flex-1">
                                    <span class="font-bold block">{{ method.name }}</span>
                                    <span class="text-sm text-gray-500 block">{{ method.description }}</span>
                                    <span class="font-bold text-blue-600">RP {{ method.cost.toLocaleString() }}</span>
                                </div>
                            </label>
                        </div>
                     </div>

                     <div class="pt-6 border-t border-gray-100">
                        <h3 class="text-xl font-bold mb-4">Metode Pembayaran</h3>
                        <label class="flex items-center gap-3 p-4 border border-blue-400 bg-blue-50 rounded-xl cursor-pointer">
                            <input type="radio" v-model="paymentMethod" value="COD" checked class="accent-blue-500 w-5 h-5">
                            <div>
                                <span class="font-bold block">COD (Bayar di Tempat)</span>
                                <span class="text-sm text-gray-500">Bayar di tempat saat kurir tiba.</span>
                            </div>
                        </label>
                     </div>

                     <button type="submit" class="w-full text-white py-4 rounded-full font-bold transition-colors mt-6 hover:opacity-80" style="background-color: #547792;">
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
                         <p class="text-sm">{{ form.address }}</p>
                         <p class="text-sm">{{ form.city }}, {{ form.province }} {{ form.zipCode }}</p>
                         <p class="text-sm">{{ form.phone }}</p>
                         <button @click="step = 2" class="text-sm text-white font-bold mt-2 px-3 py-1 rounded transition-colors hover:opacity-80" style="background-color: #547792;">Ubah Data</button>
                     </div>

                     <div class="bg-gray-50 p-4 rounded-xl">
                        <h3 class="font-bold text-gray-500 text-sm uppercase mb-2">Metode Pengiriman:</h3>
                        <p class="font-bold">{{ shippingMethods.find(m => m.id === selectedShippingMethod)?.name }}</p>
                        <p class="text-sm">{{ shippingMethods.find(m => m.id === selectedShippingMethod)?.description }}</p>
                        <p class="font-bold text-blue-600">RP {{ selectedShippingCost.toLocaleString() }}</p>
                        <button @click="step = 2" class="text-sm text-white font-bold mt-2 px-3 py-1 rounded transition-colors hover:opacity-80" style="background-color: #547792;">Ubah Metode</button>
                    </div>

                     <div class="bg-gray-50 p-4 rounded-xl">
                        <h3 class="font-bold text-gray-500 text-sm uppercase mb-2">Pembayaran:</h3>
                        <p class="font-bold">{{ paymentMethod }}</p>
                    </div>
                 </div>

                 <button @click="handleCheckout" class="w-full text-white font-bold py-3 rounded-full transition-colors hover:opacity-80" style="background-color: #547792;">
                    Konfirmasi & Proses Pesanan
                 </button>
             </div>

             <!-- Order Summary (Always Visible) -->
             <div class="bg-white p-8 rounded-[2rem] shadow-lg h-fit">
                 <h2 class="text-2xl font-bold mb-6">Ringkasan Pesanan</h2>
                 <div class="space-y-4 mb-8">
                     <div v-for="item in cartStore.items" :key="item.id" class="flex justify-between items-center py-2 border-b border-gray-50">
                         <span class="text-gray-700 font-medium">{{ item.product.name }} <span class="text-xs text-gray-400">x{{ item.quantity }}</span></span>
                         <span class="font-bold text-black">RP {{ (item.product.price * item.quantity).toLocaleString() }}</span>
                     </div>
                 </div>
                 <div class="space-y-2 py-4 border-t border-dashed border-gray-300">
                    <div class="flex justify-between">
                        <span class="text-gray-700">Subtotal</span>
                        <span class="font-bold text-black">RP {{ cartStore.totalPrice.toLocaleString() }}</span>
                    </div>
                    <div class="flex justify-between">
                        <span class="text-gray-700">Pengiriman</span>
                        <span class="font-bold text-blue-600">RP {{ selectedShippingCost.toLocaleString() }}</span>
                    </div>
                 </div>
                 <div class="flex justify-between border-t border-dashed border-gray-300 pt-6 text-xl">
                     <span class="font-bold">Total</span>
                     <span class="font-extrabold text-black">RP {{ totalWithShipping.toLocaleString() }}</span>
                 </div>
             </div>
        </div>
    </div>
  </div>
</template>
