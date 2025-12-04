import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import ui from "@nuxt/ui/vite"

export default defineConfig({
  plugins: [
    vue(),
    ui({
      ui: {
        colors: {
          primary: "neutral"
        }
      }
    })
  ],
  server: {
    host: "0.0.0.0",
    port: 3000,
    open: true
  }
})
