/**
 * uni-app + Vue3 Vite 配置
 *
 * 解决：API 模块目录命名为 @api/，避免与 /api 代理前缀冲突。
 * 真正的后端 API 请求走 /api（代理到 8080）。
 */
import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import path from 'path'

export default defineConfig({
  plugins: [uni()],

  resolve: {
    alias: [
      // @api/* → 本地 @api/ 目录（存放 API 模块文件）
      { find: /^@api\/(.*)/, replacement: path.resolve(__dirname, '@api/$1') },
    ],
  },

  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
        secure: false,
        ws: false,
        bypass(req) {
          // .js 等模块文件不走代理（@api/ 模块不带 /api 前缀，不会走到这里）
          if (req.url && /\.\w+$/.test(req.url)) {
            return false
          }
        },
      },
    },
  },
})
