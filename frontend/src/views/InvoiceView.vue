<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import api from '../services/api';

const route = useRoute();

const loading = ref(true);
const error = ref('');
const invoice = ref<any>(null);
const items = ref<any[]>([]);
const shippingCost = ref(0);
const shippingMethodName = ref('COD (Bayar di Tempat)');

const orderId = route.params.id as string;

const printPDF = () => {
  window.print();
};

const formatDateOnly = (date: string) => {
    return new Date(date).toLocaleDateString('id-ID', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
    });
};

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
            shippingCost.value = response.data.shipping_cost || 0;
            
            // Fetch shipping method details to get the name
            if (invoice.value.shipping_method_id) {
                try {
                    const methodRes = await api.get(`/shipping-methods/${invoice.value.shipping_method_id}/cost`);
                    // We only have the cost, let's fetch all methods to get the name
                    const allMethods = await api.get('/shipping-methods');
                    const selectedMethod = allMethods.data.find((m: any) => m.id === invoice.value.shipping_method_id);
                    if (selectedMethod) {
                        shippingMethodName.value = `${selectedMethod.name}${selectedMethod.description ? ' - ' + selectedMethod.description : ''}`;
                    }
                } catch (e) {
                    console.error('Failed to fetch shipping method details:', e);
                    // Keep default name
                }
            }
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
    <div class="min-h-screen bg-gray-100 text-black py-8">
        <!-- Loading State -->
        <div v-if="loading" class="container mx-auto px-4">
            <div class="text-center py-12">
                <p class="text-xl font-bold animate-pulse">Memuat invoice...</p>
            </div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="container mx-auto px-4">
            <div class="bg-white p-8 rounded-lg shadow max-w-2xl mx-auto">
                <h2 class="text-2xl font-bold text-red-600 mb-4">Error</h2>
                <p class="text-gray-700 mb-6">{{ error }}</p>
                <router-link to="/cart" class="inline-block text-white px-6 py-3 rounded-full font-bold transition-colors hover:opacity-80" style="background-color: #FF5555;">
                    Kembali ke Keranjang
                </router-link>
            </div>
        </div>

        <!-- Invoice Found -->
        <div v-else-if="invoice" class="container mx-auto px-4 print:p-0">
            <!-- Invoice Content -->
            <div class="bg-white p-8 rounded-lg shadow print:shadow-none print:rounded-none print:p-0">
                <div class="max-w-4xl mx-auto">
                    <!-- Header -->
                    <div class="mb-6 pb-6 border-b border-gray-300 print:mb-4 print:pb-4">
                        <div class="flex justify-between items-center">
                            <div>
                                <h1 class="text-2xl font-bold text-red-600">NIGHT STALKERS</h1>
                                <p class="text-xs text-gray-600 mt-1">PT Sampulkreativ</p>
                            </div>
                            <div class="text-right">
                                <h2 class="text-3xl font-bold text-gray-800">INVOICE</h2>
                                <p class="text-xs text-gray-500 mt-1">#{{ invoice.id }}</p>
                            </div>
                        </div>
                    </div>

                    <!-- Invoice Info -->
                    <div class="grid grid-cols-2 gap-6 mb-6 print:mb-4">
                        <div>
                            <h3 class="font-bold text-xs text-gray-600 mb-2 uppercase">Tagihan Kepada:</h3>
                            <div v-if="invoice?.guest_info" class="text-sm space-y-0.5">
                                <p class="font-bold">{{ invoice.guest_info.full_name || '-' }}</p>
                                <p class="text-xs">{{ invoice.guest_info.address || '-' }}</p>
                                <p class="text-xs">{{ invoice.guest_info.city || '-' }}, {{ invoice.guest_info.province || '-' }} {{ invoice.guest_info.postal_code || '-' }}</p>
                                <p class="text-xs text-gray-600">{{ invoice.guest_info.email || '-' }}</p>
                                <p class="text-xs text-gray-600">{{ invoice.guest_info.phone || '-' }}</p>
                            </div>
                            <div v-else class="text-xs text-gray-500">Informasi tidak tersedia</div>
                        </div>
                        <div class="text-right">
                            <div class="mb-2">
                                <span class="text-xs text-gray-600">TANGGAL:</span>
                                <p class="font-bold text-sm">{{ formatDateOnly(invoice.created_at) }}</p>
                            </div>
                        </div>
                    </div>

                    <!-- Items Table -->
                    <div class="mb-4 print:mb-3">
                        <table class="w-full text-xs border-collapse">
                            <thead>
                                <tr class="border-b border-gray-800">
                                    <th class="text-left py-1 px-2 font-bold">PRODUK</th>
                                    <th class="text-center py-1 px-2 font-bold">QTY</th>
                                    <th class="text-right py-1 px-2 font-bold">HARGA</th>
                                    <th class="text-right py-1 px-2 font-bold">JUMLAH</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(item, idx) in items" :key="idx" class="border-b border-gray-200">
                                    <td class="py-2 px-2">{{ item.name }}</td>
                                    <td class="text-center py-2 px-2">{{ item.quantity }}</td>
                                    <td class="text-right py-2 px-2">Rp {{ item.price.toLocaleString() }}</td>
                                    <td class="text-right py-2 px-2 font-bold">Rp {{ (item.total).toLocaleString() }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <!-- Shipping Info -->
                    <div class="mb-4 pb-4 border-b border-gray-300 print:mb-3 print:pb-3">
                        <div class="grid grid-cols-2 gap-4 text-xs">
                            <div>
                                <span class="text-gray-600">METODE PENGIRIMAN:</span>
                                <p class="font-bold text-sm">{{ shippingMethodName }}</p>
                            </div>
                            <div class="text-right">
                                <span class="text-gray-600">BIAYA PENGIRIMAN:</span>
                                <p class="font-bold text-sm">Rp {{ shippingCost.toLocaleString() }}</p>
                            </div>
                        </div>
                    </div>

                    <!-- Summary -->
                    <div class="flex justify-end mb-4 print:mb-2">
                        <div class="w-56 text-xs">
                            <div class="flex justify-between py-1 px-2 bg-gray-50 border-b border-gray-300">
                                <span class="font-semibold">SUBTOTAL</span>
                                <span class="font-semibold">Rp {{ invoice.total_amount.toLocaleString() }}</span>
                            </div>
                            <div class="flex justify-between py-1 px-2 bg-gray-50 border-b border-gray-300">
                                <span class="font-semibold">PENGIRIMAN</span>
                                <span class="font-semibold">Rp {{ shippingCost.toLocaleString() }}</span>
                            </div>
                            <div class="flex justify-between py-2 px-2 bg-orange-500 text-white font-bold">
                                <span>TOTAL</span>
                                <span>Rp {{ (invoice.total_amount + shippingCost).toLocaleString() }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- Notes -->
                    <div class="mb-2 pb-2 border-b border-gray-300 print:mb-1 print:pb-1">
                        <h3 class="font-bold text-xs text-gray-600 mb-1">CATATAN:</h3>
                        <p class="text-xs text-gray-700">
                            Silakan transfer ke rekening kami atau bayar di tempat (COD). Terima kasih!
                        </p>
                    </div>

                    <!-- Signature -->
                    <div class="flex justify-end print:mt-2">
                        <div class="text-center text-xs">
                            <p class="text-gray-600 mb-4 print:mb-2">Hormat,</p>
                            <p class="font-bold text-sm">PT Sampulkreativ</p>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Action Buttons (Hidden when printing) -->
            <div class="print:hidden mt-4 grid grid-cols-3 gap-2">
                <router-link to="/" class="block bg-gray-200 text-black px-3 py-2 rounded-lg font-bold hover:bg-gray-300 transition-colors text-center text-xs">
                    Beranda
                </router-link>
                <router-link to="/shop" class="block text-white px-3 py-2 rounded-lg font-bold hover:opacity-80 transition-opacity text-center text-xs" style="background-color: #8BAE66;">
                    Belanja
                </router-link>
                <button @click="printPDF" class="block text-white px-3 py-2 rounded-lg font-bold hover:opacity-80 transition-opacity text-xs" style="background-color: #547792;">
                    Cetak
                </button>
            </div>
        </div>

        <!-- Not Found State -->
        <div v-else class="container mx-auto px-4">
            <div class="bg-white p-8 rounded-lg shadow max-w-2xl mx-auto text-center">
                <p class="text-gray-600 mb-6">Invoice tidak ditemukan</p>
                <router-link to="/cart" class="inline-block text-white px-6 py-3 rounded-full font-bold transition-colors hover:opacity-80" style="background-color: #FF5555;">
                    Kembali ke Keranjang
                </router-link>
            </div>
        </div>
    </div>
</template>

<!-- Print Styles -->
<style>
    @media print {
        body {
            background: white;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 100% !important;
            padding: 0 !important;
        }
        * {
            page-break-inside: avoid !important;
            break-inside: avoid !important;
        }
        html, body {
            height: 100%;
            margin: 0;
            padding: 0;
        }
    }
    
    @page {
        size: A4;
        margin: 10mm;
    }
</style>
