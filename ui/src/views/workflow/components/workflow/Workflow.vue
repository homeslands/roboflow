<script setup lang="ts">
import { Background } from '@vue-flow/background'
import { type Edge, type Node, useVueFlow, VueFlow } from '@vue-flow/core'
import { MiniMap } from '@vue-flow/minimap'
import { WorkflowControls } from './controls'
import { edgeTypes } from './edge-type'
import { nodeTypes } from './node-type'

import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'

interface Props {
  nodes: Node[]
  edges: Edge[]
}
const props = defineProps<Props>()
const { onInit } = useVueFlow()

onInit((instance) => {
  instance.fitView()
})
</script>

<template>
  <VueFlow
    class="relative"
    :nodes="props.nodes"
    :edges="props.edges"
    :default-viewport="{ zoom: 1.5 }"
    :min-zoom="0.2"
    :node-types="nodeTypes"
    :edge-types="edgeTypes"
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
