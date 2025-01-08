<script setup>
import { useSidebar } from '@/components/ui/sidebar'

const { open } = useSidebar()
</script>

<template>
  <img alt="Vue logo" src="@/assets/logo.svg" width="32" height="32">
  <Transition name="fade">
    <span v-if="open" class="ml-2 font-bold">Roboflow</span>
  </Transition>
</template>

<style lang="css">
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
