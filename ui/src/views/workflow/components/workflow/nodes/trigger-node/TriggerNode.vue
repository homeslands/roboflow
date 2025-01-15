<script setup lang="ts">
import type { TriggerNodeDefinition } from '@/types/workflow/node/trigger-node-definition'
import type { NodeProps } from '@vue-flow/core'
import { Separator } from '@/components/ui/separator'
import BaseNode from '../BaseNode.vue'
import TriggerConfigForm from './TriggerConfigForm.vue'

const props = defineProps<NodeProps>()

const category = computed<string>(() => {
  const type = props.data.definition.type as TriggerNodeDefinition['type']
  if (type === 'ON_DEMAND') {
    return 'Trigger: On Demand'
  }
  return 'Trigger: Schedule'
})
</script>

<template>
  <BaseNode
    :node-props="props"
    :node-category="category"
  >
    <template #popover-content>
      <div class="flex flex-col gap-y-1.5 text-center sm:text-left">
        <h3 class="text-lg font-semibold leading-none tracking-tight">
          On demand trigger
        </h3>
        <p class="text-sm text-muted-foreground">
          Manually run this workflow as a user,
          within another workflow or via an API call.
        </p>
        <Separator />
        <TriggerConfigForm />
      </div>
    </template>
  </BaseNode>
</template>
