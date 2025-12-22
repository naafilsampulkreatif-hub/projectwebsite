<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import api from '../services/api';

const route = useRoute();

const loading = ref(true);
const error = ref('');
const invoice = ref<any>(null);
const items = ref<any[]>([]);

const orderId = route.params.id as string;

onMounted(async () => {
    if (!orderId) {
        error.value = 'Order ID tidak ditemukan';
        loading.value = false;
        return;
    }

    try {
        const response = await api.get(`/invoice/${orderId}`);
        if (response.data) {
            invoice.value = response.data.order;
            items.value = response.data.items || [];
        }
    } catch (e: any) {
        console.error('Failed to fetch invoice:', e);
        error.value = e.response?.data?.message || 'Gagal memuat invoice';
    } finally {
        loading.value = false;
    }
});

const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString('id-ID', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};

const getStatusColor = (status: string) => {
    switch (status) {
        case 'pending':
            return 'bg-yellow-100 text-yellow-800';
        case 'paid':
            return 'bg-blue-100 text-blue-800';
        case 'shipped':
            return 'bg-green-100 text-green-800';
        case 'cancelled':
            return 'bg-red-100 text-red-800';
        default:
            return 'bg-gray-100 text-gray-800';
    }
};

const getStatusLabel = (status: string) => {
    switch (status) {
        case 'pending':
            return 'Menunggu Pembayaran';
        case 'paid':
            return 'Sudah Dibayar';
        case 'shipped':
            return 'Dikirim';
        case 'cancelled':
            return 'Dibatalkan';
        default:
            return status;
    }
};
</script>

<template>
    <div class="min-h-screen bg-[#F3F3F3] text-black py-12">
        <div class="container mx-auto px-4">
            <h1 class="text-4xl font-extrabold mb-10 border-b-4 border-neon-green inline-block pb-2">Invoice Pesanan</h1>

            <div v-if="loading" class="text-center py-12">
                <p class="text-xl font-bold animate-pulse">Memuat invoice...</p>
            </div>

            <div v-else-if="error" class="bg-white p-8 rounded-[2rem] shadow-lg max-w-2xl mx-auto">
                <h2 class="text-2xl font-bold text-red-600 mb-4">Error</h2>
                <p class="text-gray-700 mb-6">{{ error }}</p>
                <router-link to="/cart" class="bg-red-600 text-white px-6 py-3 rounded-full font-bold hover:bg-red-400 transition-colors">
                    Kembali ke Keranjang
                </router-link>
            </div>

            <div v-else-if="invoice" class="grid grid-cols-1 lg:grid-cols-2 gap-12">
                <!-- Invoice Detail -->
                <div class="bg-white p-8 rounded-[2rem] shadow-lg">
                    <div class="mb-6 pb-6 border-b border-gray-200">
                        <h2 class="text-2xl font-bold mb-4">Nomor Pesanan: #{{ invoice.id }}</h2>
                        <p class="text-gray-600 mb-4">{{ formatDate(invoice.created_at) }}</p>
                        <div :class="['px-4 py-2 rounded-full font-bold inline-block', getStatusColor(invoice.status)]">
                            {{ getStatusLabel(invoice.status) }}
                        </div>
                    </div>

                    <!-- Guest Info -->
                    <div v-if="invoice.guest_info" class="mb-8">
                        <h3 class="text-xl font-bold mb-4">Informasi Pengiriman</h3>
                        <div class="bg-gray-50 p-4 rounded-xl space-y-2">
                            <p><span class="font-bold">Nama:</span> {{ invoice.guest_info.full_name }}</p>
                            <p><span class="font-bold">Email:</span> {{ invoice.guest_info.email }}</p>
                            <p><span class="font-bold">No. HP:</span> {{ invoice.guest_info.phone }}</p>
                            <p><span class="font-bold">Alamat:</span> {{ invoice.guest_info.address }}</p>
                            <p><span class="font-bold">Kota:</span> {{ invoice.guest_info.city }}, {{ invoice.guest_info.province }}</p>
                            <p><span class="font-bold">Kode Pos:</span> {{ invoice.guest_info.postal_code }}</p>
                        </div>
                    </div>

                    <!-- Payment Status -->
                    <div class="mb-8 pb-8 border-b border-gray-200">
                        <h3 class="text-xl font-bold mb-4">Status Pembayaran</h3>
                        <div v-if="invoice.status === 'pending'" class="bg-yellow-50 p-4 rounded-xl border border-yellow-200">
                            <p class="font-bold text-yellow-800 mb-2">⏳ Menunggu Pembayaran</p>
                            <p class="text-sm text-gray-700">Total: <span class="font-bold">RP {{ invoice.total_amount.toLocaleString() }}</span></p>
                            <p class="text-sm text-gray-600 mt-2">Hubungi layanan pelanggan untuk informasi pembayaran lebih lanjut.</p>
                        </div>
                        <div v-else-if="invoice.status === 'paid'" class="bg-green-50 p-4 rounded-xl border border-green-200">
                            <p class="font-bold text-green-800">✓ Pembayaran Diterima</p>
                            <p class="text-sm text-gray-700 mt-2">Terima kasih telah melakukan pembayaran. Pesanan Anda sedang diproses.</p>
                        </div>
                        <div v-else-if="invoice.status === 'shipped'" class="bg-blue-50 p-4 rounded-xl border border-blue-200">
                            <p class="font-bold text-blue-800">📦 Pesanan Dikirim</p>
                            <p class="text-sm text-gray-700 mt-2">Pesanan Anda sedang dalam perjalanan. Cek status pengiriman reguler.</p>
                        </div>
                    </div>
                </div>

                <!-- Items Summary -->
                <div class="bg-white p-8 rounded-[2rem] shadow-lg h-fit">
                    <h2 class="text-2xl font-bold mb-6">Ringkasan Pesanan</h2>
                    <div class="space-y-4 mb-8">
                        <div v-for="(item, idx) in items" :key="idx" class="flex justify-between items-center py-3 border-b border-gray-50">
                            <div>
                                <span class="font-bold block">{{ item.name }}</span>
                                <span class="text-xs text-gray-400">x{{ item.quantity }}</span>
                            </div>
                            <div class="text-right">
                                <span class="block text-sm text-gray-600">RP {{ item.price.toLocaleString() }}</span>
                                <span class="font-bold">RP {{ (item.total).toLocaleString() }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="flex justify-between border-t border-dashed border-gray-300 pt-6 text-xl mb-8">
                        <span class="font-bold">Total</span>
                        <span class="font-extrabold text-neon-green">RP {{ invoice.total_amount.toLocaleString() }}</span>
                    </div>

                    <div class="space-y-3">
                        <router-link to="/" class="block w-full bg-gray-200 text-black px-6 py-3 rounded-full font-bold hover:bg-gray-300 transition-colors text-center">
                            Kembali ke Beranda
                        </router-link>
                        <router-link to="/shop" class="block w-full bg-neon-green text-black px-6 py-3 rounded-full font-bold hover:opacity-90 transition-opacity text-center">
                            Lanjut Belanja
                        </router-link>
                    </div>
                </div>
            </div>

            <div v-else class="bg-white p-8 rounded-[2rem] shadow-lg max-w-2xl mx-auto text-center">
                <p class="text-gray-600 mb-6">Invoice tidak ditemukan</p>
                <router-link to="/cart" class="bg-red-600 text-white px-6 py-3 rounded-full font-bold hover:bg-red-400 transition-colors">
                    Kembali ke Keranjang
                </router-link>
            </div>
        </div>
    </div>
</template>
