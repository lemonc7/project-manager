import { defineStore } from "pinia"
import { useWindowSize } from "@vueuse/core"
import { computed } from "vue"

// 定义常用的断点值
const BREAKPOINTS = {
  SM: 640,
  MD: 768,
  LG: 1024,
  XL: 1280
}

export const useScreenStore = defineStore("screen", () => {
  const { width } = useWindowSize()
  const isMobile = computed(() => width.value < BREAKPOINTS.MD)

  return {
    width,
    isMobile
  }
})
