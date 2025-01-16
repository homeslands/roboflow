<script setup lang="ts">
import type { ControlRaybotNodeDefinition } from '@/types/workflow/node/definition'
import type { NodeProps } from '@vue-flow/core'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { PlusCircleIcon } from 'lucide-vue-next'
import { useAddNode } from '../../composables/use-add-node'
import { nodeMenuItems, type SubMenuItem } from './node-menu-items'

const props = defineProps<NodeProps>()

const { addControlRaybotNode } = useAddNode(props.id)

function handleNewNode(item: SubMenuItem) {
  if (item.nodeType === 'CONTROL_RAYBOT') {
    const definition = item.definition as ControlRaybotNodeDefinition
    addControlRaybotNode(definition.type)
  }
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger>
      <Tooltip>
        <TooltipTrigger as-child>
          <div class="absolute flex items-center pb-1 -translate-y-1/2 top-1/2 left-full group">
            <div class="border-t-[2px] w-4 border-muted-foreground" />
            <PlusCircleIcon class="rounded-full size-4 text-muted-foreground group-hover:text-primary group-hover:shadow-md" />
          </div>
        </TooltipTrigger>
        <TooltipContent>
          Add another node
        </TooltipContent>
      </Tooltip>
    </DropdownMenuTrigger>
    <DropdownMenuContent class="absolute w-56 translate-x-20 -translate-y-1/2 left-full top-1/2">
      <DropdownMenuLabel class="px-1">
        <Input
          id="search"
          placeholder="Search..."
        />
      </DropdownMenuLabel>
      <ScrollArea>
        <div v-for="(nodeType, index) in nodeMenuItems" :key="index">
          <DropdownMenuLabel class="text-xs font-semibold">
            {{ nodeType.label }}
          </DropdownMenuLabel>
          <DropdownMenuSub v-for="(subMenu, subMenuIndex) in nodeType.subMenu" :key="subMenuIndex">
            <DropdownMenuSubTrigger>
              <component :is="subMenu.icon" class="w-4 h-4 p-0.5 rounded mr-2 text-black bg-white" />
              {{ subMenu.label }}
            </DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent>
                <DropdownMenuItem
                  v-for="(item, itemIndex) in subMenu.subItems"
                  :key="itemIndex"
                  @click="() => handleNewNode(item)"
                >
                  {{ item.label }}
                </DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        </div>
      </ScrollArea>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
