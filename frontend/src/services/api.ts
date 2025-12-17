import axios from 'axios';

// Generate a simple session id for guest sessions
function generateSessionId(): string {
  return 'sess_' + Date.now() + '_' + Math.random().toString(36).slice(2, 9);
}

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

    // Ensure a session_id exists for guest users (session-based cart)
    let sessionId = localStorage.getItem('session_id');
    if (!sessionId) {
      sessionId = generateSessionId();
      localStorage.setItem('session_id', sessionId);
    }

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
