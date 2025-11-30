import path from "node:path"
import Components from "unplugin-vue-components/vite"
import { defineConfig } from "vite"
import tailwindcss from "@tailwindcss/vite"
import vue from "@vitejs/plugin-vue"

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    Components({
      dirs: ["src/components","src/pages"],
      extensions: ["vue"],
      deep: true,
      dts: "src/components.d.ts"
    })
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src")
    }
  },
  server: {
    port: 3000,
    open: true
  }
})
