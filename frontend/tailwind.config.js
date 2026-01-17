/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  // Important: Disable Tailwind's preflight to avoid conflicts with Ant Design
  corePlugins: {
    preflight: false,
  },
  theme: {
    extend: {
      // You can extend Tailwind theme here if needed
      // Ant Design will handle most component styling
    },
  },
  plugins: [],
}
