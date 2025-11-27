<script setup lang="ts">
import { RouterView } from "vue-router"
import { Toaster } from "vue-sonner"
import "vue-sonner/style.css"
</script>

<template>
  <Toaster rich-colors position="top-center" close-button :duration="2000" />
  <RouterView />
</template>
