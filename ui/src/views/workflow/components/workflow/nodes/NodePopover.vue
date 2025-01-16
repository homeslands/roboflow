<script setup lang="ts">
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerTitle,
} from '@/components/ui/drawer'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { useMediaQuery } from '@vueuse/core'
import { VisuallyHidden } from 'radix-vue'

const isOpen = defineModel<boolean>('isOpen', { required: true })
const isMobile = useMediaQuery('(max-width: 767px)')
const isTablet = useMediaQuery('(min-width: 768px) and (max-width: 1023px)')
</script>

<template>
  <Drawer v-if="isMobile" v-model:open="isOpen">
    <DrawerContent class="w-[100dvw] flex">
      <VisuallyHidden>
        <DrawerTitle />
        <DrawerDescription />
      </VisuallyHidden>
      <slot name="content" />
    </DrawerContent>
  </Drawer>

  <Dialog v-else-if="isTablet" v-model:open="isOpen">
    <DialogContent class="w-[100dvw] max-w-[calc(100dvw-20px)] sm:max-w-[600px] flex">
      <VisuallyHidden>
        <DialogTitle />
        <DialogDescription />
      </VisuallyHidden>
      <slot name="content" />
    </DialogContent>
  </Dialog>

  <Popover v-else v-model:open="isOpen" :modal="true">
    <!-- I don't know why we need to add PopoverTrigger to make it works -->
    <PopoverTrigger />
    <PopoverContent
      class="w-[500px] flex"
      side="right"
      :side-offset="30"
    >
      <slot name="content" />
    </PopoverContent>
  </Popover>
</template>
