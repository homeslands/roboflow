<script setup lang="ts">
import type { TriggerNodeDefinition, TriggerType } from '@/types/workflow/node/definition'
import { CircleXIcon, LoaderCircleIcon } from 'lucide-vue-next'

interface Props {
  nodeId: string
  type: TriggerType
  definition: TriggerNodeDefinition
}

const props = defineProps<Props>()

const FormAsyncComponent = defineAsyncComponent({
  loader: (): Promise<Component> => {
    if (props.type === 'ON_DEMAND') {
      return import('./OnDemandTrigerConfigForm.vue')
    }
    return import('./ScheduleTriggerForm.vue')
  },
  loadingComponent: LoaderCircleIcon,
  delay: 200,
  errorComponent: CircleXIcon,
  timeout: 5000,
})
</script>

<template>
  <FormAsyncComponent
    :node-id="props.nodeId"
    :definition="props.definition"
  />
</template>
