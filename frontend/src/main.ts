import { createApp } from "vue"
import App from "./App.vue"
import { QueryClient, VueQueryPlugin } from "@tanstack/vue-query"
import router from "./router"
import ui from "@nuxt/ui/vue-plugin"
import "./style.css"

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1 * 60 * 1000,
      refetchOnWindowFocus: false,
      retry: 2
    }
  }
})

const app = createApp(App)
app.use(VueQueryPlugin, { queryClient })
app.use(router)
app.use(ui)

app.mount("#app")
