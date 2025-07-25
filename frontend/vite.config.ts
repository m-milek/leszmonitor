import { defineConfig } from 'vite'
import viteReact from '@vitejs/plugin-react'

import {tanstackRouter} from '@tanstack/router-plugin/vite'
import {resolve} from "pathe";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    tanstackRouter(),
    viteReact(),
  ],
  //@ts-ignore
  test: {
    globals: true,
    environment: 'jsdom',
  },
  resolve: {
    alias: {
      //@ts-ignore
      '@': resolve(__dirname, './src'),
    },
  },
  server: {
    watch: {
      usePolling: true
    }
  }
})
