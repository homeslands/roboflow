<script setup lang="ts">
import type { Edge, Node } from '@vue-flow/core'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { XIcon } from 'lucide-vue-next'
import { v4 } from 'uuid'
import { Workflow, WorkflowContainer } from './components/workflow'

const defaultNode: Node = {
  id: v4(),
  type: 'EMPTY',
  position: { x: 0, y: 0 },
}

const route = useRoute()

const workflowName = ref<string>(route.query.name as string | undefined || 'Untitled Workflow')

const nodes = ref<Node[]>([defaultNode])
const edges = ref<Edge[]>([])
</script>

<template>
  <div>
    <div class="flex items-center justify-between">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink>
              <RouterLink to="/workflows" class="flex items-center text-xs">
                Workflows
              </RouterLink>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <Input
              type="text"
              placeholder="Enter name"
              :model-value="workflowName"
            />
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <div class="flex items-center space-x-2">
        <Button variant="outline">
          <XIcon />
          Reset
        </Button>
        <Button>
          Save
        </Button>
      </div>
    </div>

    <div class="relative flex h-full mt-3">
      <WorkflowContainer>
        <Workflow
          :nodes="nodes"
          :edges="edges"
        />
      </WorkflowContainer>
    </div>
  </div>
</template>
