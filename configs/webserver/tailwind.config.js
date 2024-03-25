/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend: {},
    fontFamily: {
      nunito: ['Nunito', 'sans-serif'],
      mono: ['JetBrains Mono', 'monospace', 'ui-monospace', 'SFMono-Regular'],
    }
  },
  plugins: [],
}

