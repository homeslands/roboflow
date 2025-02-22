<script setup lang="ts">
import type { ListRaybotCommandParams, ListRaybotCommandSort } from '@/api/raybot-command'
import type { SortPrefix } from '@/lib/sort'
import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import { useListRaybotCommandQuery } from '@/composables/use-raybot-command'
import { PlusIcon, RefreshCcwIcon } from 'lucide-vue-next'
import { columns, CommandTable } from './components/command-table'
import { CreateCommandForm } from './components/create-command-form'

type ListRaybotCommandRequiredPagingParams = Omit<ListRaybotCommandParams, 'page' | 'pageSize'> & {
  page: number
  pageSize: number
}

const route = useRoute()
const raybotId = computed(() => route.params.id as string)
const params = ref<ListRaybotCommandRequiredPagingParams>({
  page: 1,
  pageSize: 10,
})

const { data, isPending, refetch } = useListRaybotCommandQuery(raybotId, params)

function handleSortingChange(sorts: SortPrefix<ListRaybotCommandSort>[]) {
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
      <Sheet>
        <SheetTrigger>
          <Button variant="outline">
            <PlusIcon />
            Create
          </Button>
        </SheetTrigger>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Create new command</SheetTitle>
            <SheetDescription>
              Create a new command to send to raybot.
              Raybot can only process <strong>one</strong> command at a time.
            </SheetDescription>
          </SheetHeader>
          <CreateCommandForm class="mt-4" />
        </SheetContent>
      </Sheet>

      <Button
        variant="outline"
        @click="refetch"
      >
        <RefreshCcwIcon />
        Refresh
      </Button>
    </div>

    <!-- Data table -->
    <CommandTable
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
