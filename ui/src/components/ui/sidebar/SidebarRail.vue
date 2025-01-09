<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <button
    data-sidebar="rail"
    aria-label="Sidebar Rail"
    title="Sidebar Rail"
    :class="cn(props.class)"
  >
    <slot />
  </button>
</template>
