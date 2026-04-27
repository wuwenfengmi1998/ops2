import { fileURLToPath, URL } from "node:url";
import { readFileSync, writeFileSync } from "node:fs";
import { resolve, dirname } from "node:path";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";
import tailwindcss from "@tailwindcss/vite";

const __dirname = dirname(fileURLToPath(import.meta.url));
const packageJsonPath = resolve(__dirname, "package.json");

// 读取 package.json
let packageJson = JSON.parse(readFileSync(packageJsonPath, "utf-8"));

// 每次 build 时自动递增 patch 版本
const isBuild = process.argv.includes("build");
if (isBuild) {
  const parts = packageJson.version.split(".");
  const patch = Math.max(0, parseInt(parts[2] || "0", 10));
  parts[2] = String(patch + 1);
  packageJson.version = parts.join(".");
  writeFileSync(packageJsonPath, JSON.stringify(packageJson, null, 2) + "\n", "utf-8");
  console.log(`[bump] version → ${packageJson.version}`);
}

// https://vite.dev/config/
export default defineConfig({
  define: {
    // 全局可用 __APP_VERSION__，取自 package.json 的 version 字段
    __APP_VERSION__: JSON.stringify(packageJson.version),
  },
  plugins: [vue(), vueDevTools(), tailwindcss()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  build: {
    outDir: "../../backend/my_work/dist",
    assetsDir: "assets",
  },
  server: {
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8080",
        changeOrigin: true,
      },
    },
  },
});
