<template>
  <Teleport v-if="isMounted" :to="props.to">
    <slot />
  </Teleport>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"

const props = defineProps<{ to: string }>()
const isMounted = ref(false)
onMounted(() => {
  isMounted.value = true
})
</script>
