import { defineStore } from 'pinia';
import api from '../services/api';

// Define the Auth Store
export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user') || 'null'), // Load user from storage
    token: localStorage.getItem('token') || null, // Load token from storage
  }),
  getters: {
    isAuthenticated: (state) => !!state.token, // Check if logged in
    isAdmin: (state) => state.user?.role === 'admin', // Check if admin
  },
  actions: {
    // Login action
    async login(credentials: any) {
      try {
        const response = await api.post('/login', credentials);
        this.token = response.data.token;
        this.user = response.data.user;

        // Save to local storage
        localStorage.setItem('token', this.token as string);
        localStorage.setItem('user', JSON.stringify(this.user));

        return true;
      } catch (error) {
        console.error("Login failed", error);
        return false;
      }
    },
    // Register action
    async register(userData: any) {
      try {
        await api.post('/register', userData);
        return true;
      } catch (error) {
        console.error("Registration failed", error);
        return false;
      }
    },
    // Logout action
    logout() {
      this.token = null;
      this.user = null;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
  }
});
