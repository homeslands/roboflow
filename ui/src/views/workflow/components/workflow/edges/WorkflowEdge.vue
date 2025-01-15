<script setup lang="ts">
import type { EdgeProps } from '@vue-flow/core'
import { Button } from '@/components/ui/button'
import { BaseEdge, EdgeLabelRenderer, getBezierPath } from '@vue-flow/core'

const props = defineProps<EdgeProps>()

const path = computed(() => getBezierPath(props))
</script>

<template>
  <BaseEdge :path="path[0]" :marker-end="props.markerEnd" :style="props.style" />
  <EdgeLabelRenderer>
    <div
      :style="{
        position: 'absolute',
        pointerEvents: 'all',
        transform: `translate(-50%, -50%) translate(${path[1]}px, ${path[2]}px)` }"
    >
      <Button
        variant="outline"
        size="icon"
        class="w-5 h-5 text-xs leading-none border rounded-full cursor-pointer hover:shadow-lg"
      >
        +
      </Button>
    </div>
  </EdgeLabelRenderer>
</template>
