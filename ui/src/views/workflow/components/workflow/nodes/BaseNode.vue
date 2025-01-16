<script setup lang="ts">
import type { NodeProps } from '@vue-flow/core'
import { Button } from '@/components/ui/button'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from '@/components/ui/context-menu'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import { Handle, Position, useVueFlow } from '@vue-flow/core'
import { CheckIcon } from 'lucide-vue-next'
import { useDeleteNode } from '../composables/use-delete-node'
import { NodeTypeSelect } from './node-type-select'
import NodePopover from './NodePopover.vue'

interface Props {
  nodeProps: NodeProps
  nodeCategory: string
}

const props = defineProps<Props>()
const { edges, updateNodeData } = useVueFlow()

const isPopoverOpen = defineModel<boolean>('isPopoverOpen', { default: true })
const hasNextNode = computed(() => edges.value.some(edge => edge.source === props.nodeProps.id))
const isRootNode = computed(() => props.nodeProps.type === 'TRIGGER')

const isEditLabel = ref<boolean>(false)
const nodeId = computed(() => props.nodeProps.id)
const definition = computed(() => props.nodeProps.data.definition)
const nodeLabel = computed(() => definition.value.label)
const label = ref(nodeLabel.value)
watch(
  () => props.nodeProps.data.definition.label,
  newLabel => label.value = newLabel,
)

function handleSaveLabel() {
  isEditLabel.value = false
  label.value = label.value.trim()
  if (label.value === nodeLabel.value) {
    return
  }
  else if (label.value === '') {
    label.value = nodeLabel.value
    return
  }

  updateNodeData(nodeId.value, {
    definition: {
      ...definition.value,
      label: label.value,
    },
  })
}

const { deleteNode: handleDelete } = useDeleteNode(props.nodeProps.id)
</script>

<template>
  <div class="flex flex-col items-center space-y-2">
    <div class="relative flex items-center">
      <ContextMenu>
        <ContextMenuTrigger>
          <Button
            class="w-32 h-32 border-purple-400 [&_svg]:size-16"
            size="icon"
            variant="outline"
            @click="isPopoverOpen = !isPopoverOpen"
          >
            <slot name="icon">
              No icon provided
            </slot>
          </Button>
        </ContextMenuTrigger>
        <ContextMenuContent>
          <ContextMenuItem @select.stop="isPopoverOpen = true">
            Settings
          </ContextMenuItem>
          <ContextMenuItem @select.stop="isEditLabel = true">
            Rename
          </ContextMenuItem>
          <ContextMenuItem @select.stop="handleDelete">
            Delete
          </ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
      <NodePopover v-model:is-open="isPopoverOpen">
        <template #content>
          <slot name="popover-content">
            No content provided
          </slot>
        </template>
      </NodePopover>
      <NodeTypeSelect
        v-if="!hasNextNode"
        v-bind="props.nodeProps"
      />

      <Handle
        v-if="!isRootNode"
        :class="cn('workflow-handle', { invisible: false })"
        type="target"
        :position="Position.Left"
        :connectable="false"
      />
      <Handle
        :class="cn('workflow-handle', { invisible: !hasNextNode })"
        type="source"
        :position="Position.Right"
        :connectable="false"
      />
    </div>

    <div class="flex flex-col items-center space-y-1">
      <span v-if="!isEditLabel" class="text-xs">
        {{ label }}
      </span>
      <div v-else class="flex items-center space-x-1">
        <Input
          id="label-input"
          v-model="label"
          class="w-32"
          @blur="handleSaveLabel"
          @keydown.enter="handleSaveLabel"
        />
        <Button
          variant="ghost"
          size="icon"
          class="absolute -right-9 shrink-0"
          @click="handleSaveLabel"
        >
          <CheckIcon />
        </Button>
      </div>
      <span class="text-xs text-gray-500">
        {{ props.nodeCategory }}
      </span>
    </div>
  </div>
</template>

<style lang="css" scoped>
.workflow-handle {
  height: 2px;
  width: 2px;
  background: gray;
  cursor: pointer !important;
  border: gray;
}
</style>
