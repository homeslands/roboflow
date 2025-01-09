<script setup lang="ts">
import { useColorMode } from '@vueuse/core'
import { Notification, Notivue } from 'notivue'

const mode = useColorMode()
mode.value = 'dark'
</script>

<template>
  <Notivue v-slot="item">
    <Notification :item="item" />
  </Notivue>

  <RouterView />
</template>
