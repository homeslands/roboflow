<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { ChevronLeftIcon, ChevronRightIcon } from 'lucide-vue-next'
import { useSidebar } from './utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()

const { toggleSidebar, open } = useSidebar()

function getIcon() {
  return open.value ? ChevronLeftIcon : ChevronRightIcon
}
</script>

<template>
  <Button
    data-sidebar="trigger"
    variant="outline"
    size="icon"
    :class="cn('h-7 w-7', props.class)"
    @click="toggleSidebar"
  >
    <component :is="getIcon()" />
    <span class="sr-only">Toggle Sidebar</span>
  </Button>
</template>
