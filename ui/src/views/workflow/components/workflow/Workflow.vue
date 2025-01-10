<script setup lang="ts">
import { Background } from '@vue-flow/background'
import { type Node, useVueFlow, VueFlow } from '@vue-flow/core'
import { MiniMap } from '@vue-flow/minimap'
import { v4 } from 'uuid'
import { WorkflowControls } from './controls'

import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'

const defaultNodes: Node[] = [
  {
    id: v4(),
    data: { label: 'Start' },
    position: { x: 0, y: 0 },
    type: 'abc',
  },
]

const { onInit } = useVueFlow()

const nodes = ref(defaultNodes)
const edges = ref([])

onInit((instance) => {
  instance.fitView()
})
</script>

<template>
  <VueFlow
    class="relative"
    :nodes="nodes"
    :edges="edges"
    :default-viewport="{ zoom: 1.5 }"
    :min-zoom="0.2"
  >
    <Background :gap="8" class="opacity-25" />
    <MiniMap
      :zoomable="true"
      :pannable="true"
      class="hidden sm:block !bg-background rounded-lg [&>svg]:rounded-lg"
    />
    <WorkflowControls />
  </VueFlow>
</template>
