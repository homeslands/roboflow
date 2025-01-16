<script setup lang="ts">
import type { TriggerNodeDefinition, TriggerType } from '@/types/workflow/node/definition/trigger-node-definition'
import type { NodeProps } from '@vue-flow/core'
import { Separator } from '@/components/ui/separator'
import { ClockIcon, WebhookIcon } from 'lucide-vue-next'
import BaseNode from '../BaseNode.vue'
import FormAsync from './FormAsync.vue'

const props = defineProps<NodeProps>()

const triggerTypeMap: Record<
  TriggerType,
  { icon: Component, category: string, description: string }
> = {
  ON_DEMAND: {
    icon: WebhookIcon,
    category: 'Trigger: On Demand',
    description: 'Manually run this workflow as a user, within another workflow or via an API call.',
  },
  SCHEDULE: {
    icon: ClockIcon,
    category: 'Trigger: Schedule',
    description: 'Run this workflow on a schedule.',
  },
}

const definition = computed<TriggerNodeDefinition>(() => props.data.definition)
const info = computed(() => triggerTypeMap[definition.value.type])
</script>

<template>
  <BaseNode
    :node-props="props"
    :node-category="info.category"
  >
    <template #icon>
      <component :is="info.icon" />
    </template>
    <template #popover-content>
      <div class="flex flex-col gap-y-1.5 text-center sm:text-left">
        <h3 class="text-lg font-semibold leading-none tracking-tight">
          {{ info.category }}
        </h3>
        <p class="text-sm text-muted-foreground">
          {{ info.description }}
        </p>
        <Separator />
        <FormAsync
          :node-id="props.id"
          :type="definition.type"
          :definition="definition"
        />
      </div>
    </template>
  </BaseNode>
</template>
