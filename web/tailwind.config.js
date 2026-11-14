/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        sand: {
          50: '#FAF8F5',
          100: '#F5F1EA',
          200: '#E8E0D2',
          300: '#D9CCB8',
          400: '#B8A48B',
          500: '#948067',
          600: '#75634D',
          700: '#5A4B3A',
          800: '#3D3226',
          900: '#241D16',
          950: '#14100B',
        },
      },
      fontFamily: {
        serif: ['Georgia', 'Cambria', 'serif'],
        sans: ['Inter', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
