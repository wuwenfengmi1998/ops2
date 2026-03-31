import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueDevTools()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  build: {
    outDir: "../../backend/dist", // 默认是 'dist'，可以修改为你想要的目录名
    assetsDir: "assets", // 静态资源目录（相对于 outDir）
  },
  server: {
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8080",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, "/api"), // 如果需要重写路径
        // 如果后端接口没有 /api 前缀，可以这样写：
        // rewrite: (path) => path.replace(/^\/api/, '')

        // 设置代理超时配置
        // proxyTimeout: 30000, // 代理服务器等待目标服务器响应的超时时间（毫秒）
        // timeout: 30000, // 整个请求的超时时间（毫秒）
        // connectTimeout: 30000, // 连接超时（毫秒，某些版本支持）
      },
    },
  },
});
