import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    // Sign-in returns to PUBLIC_URL, so a silent move to another port would break it.
    port: 5173,
    strictPort: true,
    proxy: { '/api': 'http://localhost:8080' },
  },
})
