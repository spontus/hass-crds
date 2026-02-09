/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Home Assistant inspired dark palette
        'ha-blue': '#03a9f4',
        'ha-cyan': '#41bdf5',
        'ha-yellow': '#ffc107',
        'ha-green': '#4caf50',
        'ha-red': '#f44336',
        'slate': {
          850: '#172033',
          950: '#0b1120',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
    },
  },
  plugins: [],
}
