import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://192.168.1.11:3000',
        changeOrigin: true,
      },
      '/storage': {
        target: 'http://192.168.1.11:3000',
        changeOrigin: true,
      }
    }
  }
})
