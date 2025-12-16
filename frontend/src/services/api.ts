import axios from 'axios';

// Create an axios instance
const api = axios.create({
  baseURL: 'http://localhost:8080/api', // Backend URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to attach the token and session ID
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token'); // Get token from local storage
    if (token) {
      config.headers.Authorization = `Bearer ${token}`; // Attach to header
    }

    const sessionId = localStorage.getItem('session_id');
    if (sessionId) {
      config.headers['X-Session-ID'] = sessionId;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default api;
