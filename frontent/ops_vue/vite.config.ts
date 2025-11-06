import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueJsx(), vueDevTools()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    outDir: '../../backend/dist', // 默认是 'dist'，可以修改为你想要的目录名
    assetsDir: 'assets', // 静态资源目录（相对于 outDir）
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api'), // 如果需要重写路径
        // 如果后端接口没有 /api 前缀，可以这样写：
        // rewrite: (path) => path.replace(/^\/api/, '')
      },
    },
  },
})
