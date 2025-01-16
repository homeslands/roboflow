<script setup lang="ts">
import type { ControlRaybotNodeDefinition } from '@/types/workflow/node/control-raybot-node-definition'
import type { NodeProps } from '@vue-flow/core'
import BaseNode from '../BaseNode.vue'
import { raybotCommandLabelMap } from './command-label'

const props = defineProps<NodeProps>()

const category = computed<string>(() => {
  const type = props.data.definition.type as ControlRaybotNodeDefinition['type']
  return `Raybot Control: ${raybotCommandLabelMap[type]}`
})
</script>

<template>
  <BaseNode
    :node-props="props"
    :node-category="category"
  >
    <template #popover-content>
      <div>
        Hello
      </div>
    </template>
  </BaseNode>
</template>
