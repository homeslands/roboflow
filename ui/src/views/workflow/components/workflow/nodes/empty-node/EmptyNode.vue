<script setup lang="ts">
import type { TriggerType } from '@/types/workflow/node/trigger-node-definition'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useVueFlow } from '@vue-flow/core'
import { ClockIcon, PlusIcon, WebhookIcon } from 'lucide-vue-next'
import { v4 } from 'uuid'

const { setNodes } = useVueFlow()

// Replace current node with a new node
function replaceNode(type: TriggerType) {
  setNodes((_) => {
    const newNode = {
      id: v4(),
      type: 'TRIGGER',
      position: { x: 0, y: 0 },
      data: {
        label: type === 'ON_DEMAND' ? 'On Demand' : 'Schedule',
        definition: {
          type,
        },
      },
    }
    return [newNode]
  })
}
</script>

<template>
  <div class="flex flex-col items-center justify-center space-y-2">
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button variant="outline" size="icon" class="p-6 border-dashed">
          <PlusIcon class="w-4 h-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent class="min-w-[325px]">
        <DropdownMenuLabel>
          Select a trigger type
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem @click="replaceNode('ON_DEMAND')">
          <WebhookIcon class="w-4 h-4 p-0.5 rounded text-black bg-white" />
          On demand
        </DropdownMenuItem>
        <DropdownMenuItem disabled>
          <ClockIcon class="w-4 h-4 p-0.5 rounded text-black bg-white" />
          Schedule (coming soon)
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
    <p class="text-xs text-gray-500">
      Select a trigger type
    </p>
  </div>
</template>
