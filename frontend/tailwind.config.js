/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'primary': '#cc9966', // Molla orange/gold
        'secondary': '#333333',
        'neon-green': '#84cc16', // Lime green from Tailwind
      }
    },
  },
  plugins: [],
}

