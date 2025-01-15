<script setup lang="ts">
import type { ListWorkflowParams, ListWorkflowSort } from '@/api/workflow'
import type { SortPrefix } from '@/lib/sort'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { useListWorkflowQuery } from '@/composables/use-workflow'
import { CreateWorkflowTempForm } from './components/create-workflow-temp-form'
import { columns, WorkflowTable } from './components/workflow-table'

type ListWorkflowRequiredPagingParams = Omit<ListWorkflowParams, 'page' | 'pageSize'> & {
  page: number
  pageSize: number
}

const params = ref<ListWorkflowRequiredPagingParams>({
  page: 1,
  pageSize: 10,
})

const { data, isPending, refetch } = useListWorkflowQuery(params)

function handleSortingChange(sorts: SortPrefix<ListWorkflowSort>[]) {
  params.value = {
    ...params.value,
    sort: sorts,
  }
  refetch()
}
</script>

<template>
  <div class="flex flex-col items-end space-y-2">
    <div class="space-x-2">
      <Dialog>
        <DialogTrigger as-child>
          <Button variant="outline">
            New workflow
          </Button>
        </DialogTrigger>
        <DialogContent class="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Create new workflow</DialogTitle>
            <DialogDescription>
              Create new workflow with a name.
            </DialogDescription>
          </DialogHeader>
          <CreateWorkflowTempForm />
        </DialogContent>
      </Dialog>
    </div>

    <WorkflowTable
      v-model:page="params.page"
      v-model:page-size="params.pageSize"
      class="w-full"
      :columns="columns"
      :is-loading="isPending"
      :data="data?.items ?? []"
      :total-items="data?.totalItems ?? 0"
      :sorts="params.sort"
      @sorts="handleSortingChange"
    />
  </div>
</template>
