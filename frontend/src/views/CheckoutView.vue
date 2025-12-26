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
  <div class="min-h-screen text-black" style="background-color: #F9DFDF;">
    <div class="container mx-auto px-3 sm:px-4 py-6 md:py-12">
<h1 class="text-2xl sm:text-3xl md:text-4xl font-extrabold mb-4 sm:mb-6 md:mb-10 border-b-4 border-blue-500 inline-block pb-1 sm:pb-2">Checkout</h1>

        <!-- Step 1: Stock Validation Loading/Error -->
        <div v-if="step === 1" class="p-4 sm:p-6 md:p-8 rounded-lg sm:rounded-2xl md:rounded-[2rem] shadow-lg max-w-2xl mx-auto text-center" style="background-color: #FBEFEF;">
            <div v-if="stockLoading">
                <p class="text-base sm:text-lg font-bold animate-pulse">Mengecek ketersediaan stok...</p>
            </div>
            <div v-else-if="!stockValid">
                <h2 class="text-lg sm:text-xl md:text-2xl font-bold text-red-600 mb-3 sm:mb-4">Stok Tidak Mencukupi</h2>
                <div class="bg-red-50 p-2 sm:p-3 md:p-4 rounded-lg sm:rounded-xl mb-4 sm:mb-6 text-left text-xs sm:text-sm md:text-base">
                    <p class="font-medium mb-1 sm:mb-2">Item berikut tidak tersedia dalam jumlah yang diminta:</p>
                    <ul class="list-disc list-inside space-y-0.5 sm:space-y-1">
                        <li v-for="issue in stockIssues" :key="issue.product_name">
                            <span class="font-bold">{{ issue.product_name }}</span>
                            (Diminta: {{ issue.requested }}, Tersedia: {{ issue.available }})
                        </li>
                    </ul>
                </div>
                <div class="flex flex-col sm:flex-row gap-2 sm:gap-3 md:gap-4 justify-center">
                    <router-link to="/cart" class="bg-slate-200 hover:bg-slate-300 px-3 sm:px-6 py-2 sm:py-3 rounded-full font-bold text-xs sm:text-base text-slate-800">Kembali ke Keranjang</router-link>
                    <button @click="checkStock" class="text-white px-3 sm:px-6 py-2 sm:py-3 rounded-full font-bold text-xs sm:text-base transition-colors hover:opacity-80" style="background-color: #547792;">Coba Lagi</button>
                </div>
            </div>
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6 md:gap-12">

             <!-- Step 2: Form Data -->
             <div v-if="step === 2" class="p-4 sm:p-6 md:p-8 rounded-lg sm:rounded-2xl md:rounded-[2rem] shadow-lg" style="background-color: #FBEFEF;">
                 <h2 class="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 flex items-center gap-2">
                    <span class="bg-black text-white w-7 sm:w-8 h-7 sm:h-8 rounded-full flex items-center justify-center text-xs sm:text-sm">1</span>
                    Informasi Pengiriman
                 </h2>
                 <form class="space-y-3 sm:space-y-4 md:space-y-6" @submit.prevent="goToConfirmation">
                     <div class="grid grid-cols-2 gap-2 sm:gap-3 md:gap-6">
                         <div>
                             <label class="block text-xs sm:text-sm font-bold mb-1 sm:mb-2 text-gray-500">Nama Depan <span class="text-blue-500">*</span></label>
                             <input v-model="form.firstName" type="text" class="w-full border border-gray-200 rounded-lg sm:rounded-xl p-2 sm:p-3 text-black text-xs sm:text-sm focus:outline-none focus:border-blue-400 transition-colors" style="background-color: #FCF8F8;" required>
                         </div>
                         <div>
                             <label class="block text-xs sm:text-sm font-bold mb-1 sm:mb-2 text-gray-500">Nama Belakang</label>
                             <input v-model="form.lastName" type="text" class="w-full border border-gray-200 rounded-lg sm:rounded-xl p-2 sm:p-3 text-black text-xs sm:text-sm focus:outline-none focus:border-blue-400 transition-colors" style="background-color: #FCF8F8;">
                         </div>
                     </div>
                     <div>
                         <label class="block text-xs sm:text-sm font-bold mb-1 sm:mb-2 text-gray-500">Email <span class="text-blue-500">*</span></label>
                         <input v-model="form.email" type="email" class="w-full border border-gray-200 rounded-lg sm:rounded-xl p-2 sm:p-3 text-black text-xs sm:text-sm focus:outline-none focus:border-blue-400 transition-colors" style="background-color: #FCF8F8;" required>
                     </div>
                     <div>
                         <label class="block text-xs sm:text-sm font-bold mb-1 sm:mb-2 text-gray-500">Alamat <span class="text-blue-500">*</span></label>
                         <input v-model="form.address" type="text" class="w-full border border-gray-200 rounded-lg sm:rounded-xl p-2 sm:p-3 text-black text-xs sm:text-sm focus:outline-none focus:border-blue-400 transition-colors" style="background-color: #FCF8F8;" required>
                     </div>
                     <div>
                         <label class="block text-xs sm:text-sm font-bold mb-1 sm:mb-2 text-gray-500">Provinsi <span class="text-blue-500">*</span></label>
                         <select v-model="form.province" class="w-full border border-gray-200 rounded-lg sm:rounded-xl p-2 sm:p-3 text-black text-xs sm:text-sm focus:outline-none focus:border-blue-400 transition-colors" style="background-color: #FCF8F8;" required>
                             <option value="">Pilih Provinsi</option>
                             <option v-for="prov in provinces" :key="prov" :value="prov">{{ prov }}</option>
                         </select>
                     </div>
                     <div v-if="form.province">
                         <label class="block text-xs sm:text-sm font-bold mb-1 sm:mb-2 text-gray-500">Kota/Kabupaten <span class="text-red-500">*</span></label>
                         <select v-model="form.city" class="w-full border border-gray-200 rounded-lg sm:rounded-xl p-2 sm:p-3 text-black text-xs sm:text-sm focus:outline-none focus:border-blue-400 transition-colors" style="background-color: #FCF8F8;" required>
                             <option value="">Pilih Kota</option>
                             <option v-for="city in citiesForProvince" :key="city" :value="city">{{ city }}</option>
                         </select>
                     </div>
                     <div class="grid grid-cols-2 gap-2 sm:gap-3 md:gap-6">
                        <div>
                             <label class="block text-xs sm:text-sm font-bold mb-1 sm:mb-2 text-gray-500">Kode Pos <span class="text-red-500">*</span></label>
                             <input v-model="form.zipCode" type="text" class="w-full border border-gray-200 rounded-lg sm:rounded-xl p-2 sm:p-3 text-black text-xs sm:text-sm focus:outline-none focus:border-blue-400 transition-colors" style="background-color: #FCF8F8;" required>
                         </div>
                         <div>
                             <label class="block text-xs sm:text-sm font-bold mb-1 sm:mb-2 text-gray-500">No. HP <span class="text-red-500">*</span></label>
                             <input v-model="form.phone" type="tel" class="w-full border border-gray-200 rounded-lg sm:rounded-xl p-2 sm:p-3 text-black text-xs sm:text-sm focus:outline-none focus:border-blue-400 transition-colors" style="background-color: #FCF8F8;" required>
                         </div>
                     </div>

                     <!-- Shipping Method Section -->
                     <div class="pt-4 sm:pt-6 border-t border-gray-100">
                        <h3 class="text-lg sm:text-xl font-bold mb-3 sm:mb-6">Metode Pembayaran</h3>
                        
                        <!-- Shipping Method -->
                        <div>
                            <div v-if="loadingShipping" class="text-center py-2 sm:py-3 md:py-4">
                                <p class="text-xs sm:text-sm text-gray-500">Memuat metode pembayaran...</p>
                            </div>
                            <div v-else-if="shippingMethods.length === 0" class="text-center py-2 sm:py-3 md:py-4">
                                <p class="text-xs sm:text-sm text-gray-500">Tidak ada metode pembayaran tersedia</p>
                            </div>
                            <div v-else class="space-y-2 sm:space-y-3">
                                <label v-for="method in shippingMethods" :key="method.id" class="flex items-center gap-2 sm:gap-3 p-2 sm:p-3 md:p-4 border border-gray-300 rounded-lg sm:rounded-xl cursor-pointer hover:border-blue-400 hover:bg-blue-50 transition-all text-xs sm:text-sm"
                                       :class="selectedShippingMethod === method.id ? 'border-blue-400 bg-blue-50' : ''">
                                    <input type="radio" :value="method.id" v-model="selectedShippingMethod" class="accent-blue-500 w-4 h-4 sm:w-5 sm:h-5">
                                    <div class="flex-1">
                                        <span class="font-bold block">{{ method.name }}</span>
                                        <span class="text-xs sm:text-sm text-gray-500 block">{{ method.description }}</span>
                                        <span class="font-bold text-blue-600 text-xs sm:text-sm">RP {{ method.cost.toLocaleString() }}</span>
                                    </div>
                                </label>
                            </div>
                        </div>
                     </div>

                     <button type="submit" class="w-full text-white py-2 sm:py-3 md:py-4 rounded-full font-bold text-sm sm:text-base transition-colors mt-4 sm:mt-6 hover:opacity-80" style="background-color: #547792;">
                         Lanjut ke Konfirmasi
                     </button>
                 </form>
             </div>

             <!-- Step 3: Confirmation -->
             <div v-if="step === 3" class="p-4 sm:p-6 md:p-8 rounded-lg sm:rounded-2xl md:rounded-[2rem] shadow-lg" style="background-color: #FBEFEF;">
                <h2 class="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 flex items-center gap-2">
                    <span class="bg-black text-white w-7 sm:w-8 h-7 sm:h-8 rounded-full flex items-center justify-center text-xs sm:text-sm">2</span>
                    Konfirmasi Pesanan
                 </h2>

                 <div class="space-y-3 sm:space-y-4 md:space-y-6 mb-4 sm:mb-8">
                     <div class="bg-gray-50 p-2 sm:p-3 md:p-4 rounded-lg sm:rounded-xl">
                         <h3 class="font-bold text-gray-500 text-xs sm:text-sm uppercase mb-1 sm:mb-2">Dikirim ke:</h3>
                         <p class="font-bold text-sm sm:text-base">{{ form.firstName }} {{ form.lastName }}</p>
                         <p class="text-xs sm:text-sm">{{ form.address }}</p>
                         <p class="text-xs sm:text-sm">{{ form.city }}, {{ form.province }} {{ form.zipCode }}</p>
                         <p class="text-xs sm:text-sm">{{ form.phone }}</p>
                         <button @click="step = 2" class="text-xs sm:text-sm text-white font-bold mt-1 sm:mt-2 px-2 sm:px-3 py-1 rounded transition-colors hover:opacity-80" style="background-color: #547792;">Ubah Data</button>
                     </div>

                     <div class="bg-gray-50 p-2 sm:p-3 md:p-4 rounded-lg sm:rounded-xl">
                        <h3 class="font-bold text-gray-500 text-xs sm:text-sm uppercase mb-1 sm:mb-2">Metode Pembayaran:</h3>
                        <p class="font-bold text-sm sm:text-base">{{ shippingMethods.find(m => m.id === selectedShippingMethod)?.name }}</p>
                        <p class="text-xs sm:text-sm">{{ shippingMethods.find(m => m.id === selectedShippingMethod)?.description }}</p>
                        <p class="font-bold text-blue-600 text-sm sm:text-base">RP {{ selectedShippingCost.toLocaleString() }}</p>
                        <button @click="step = 2" class="text-xs sm:text-sm text-white font-bold mt-2 sm:mt-3 px-2 sm:px-3 py-1 rounded transition-colors hover:opacity-80" style="background-color: #547792;">Ubah Metode</button>
                    </div>
                 </div>

                 <button @click="handleCheckout" class="w-full text-white font-bold py-2 sm:py-3 md:py-4 rounded-full text-sm sm:text-base transition-colors hover:opacity-80" style="background-color: #547792;">
                    Konfirmasi & Proses Pesanan
                 </button>
             </div>

             <!-- Order Summary (Always Visible) -->
             <div class="bg-white p-4 sm:p-6 md:p-8 rounded-lg sm:rounded-2xl md:rounded-[2rem] shadow-lg h-fit">
                 <h2 class="text-lg sm:text-xl md:text-2xl font-bold mb-4 sm:mb-6">Ringkasan Pesanan</h2>
                 <div class="space-y-2 sm:space-y-3 md:space-y-4 mb-4 sm:mb-8">
                     <div v-for="item in cartStore.items" :key="item.id" class="flex justify-between items-center py-1 sm:py-2 border-b border-gray-50">
                         <span class="text-gray-700 font-medium text-xs sm:text-sm">{{ item.product.name }} <span class="text-xs text-gray-400">x{{ item.quantity }}</span></span>
                         <span class="font-bold text-black text-xs sm:text-sm">RP {{ (item.product.price * item.quantity).toLocaleString() }}</span>
                     </div>
                 </div>
                 <div class="space-y-1 sm:space-y-2 py-3 sm:py-4 border-t border-dashed border-gray-300">
                    <div class="flex justify-between text-xs sm:text-sm">
                        <span class="text-gray-700">Subtotal</span>
                        <span class="font-bold text-black">RP {{ cartStore.totalPrice.toLocaleString() }}</span>
                    </div>
                    <div class="flex justify-between text-xs sm:text-sm">
                        <span class="text-gray-700">Pengiriman</span>
                        <span class="font-bold text-blue-600">RP {{ selectedShippingCost.toLocaleString() }}</span>
                    </div>
                 </div>
                 <div class="flex justify-between border-t border-dashed border-gray-300 pt-3 sm:pt-4 md:pt-6 text-base sm:text-lg md:text-xl">
                     <span class="font-bold">Total</span>
                     <span class="font-extrabold text-black">RP {{ totalWithShipping.toLocaleString() }}</span>
                 </div>
             </div>
        </div>
    </div>
  </div>
</template>
