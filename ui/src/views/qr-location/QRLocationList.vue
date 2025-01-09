<script setup lang="ts">
import type { ListQRLocationParams, ListQRLocationSort } from '@/api/qr-location'
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
import { useListQRLocationQuery } from '@/composables/use-qr-location'
import { PlusIcon, RefreshCcwIcon } from 'lucide-vue-next'
import { CreateQRLocationForm } from './components/create-qr-location-form'
import { columns, QRLocationTable } from './components/qr-location-table'

type ListQRLocationRequiredPagingParams = Omit<ListQRLocationParams, 'page' | 'pageSize'> & {
  page: number
  pageSize: number
}

const params = ref<ListQRLocationRequiredPagingParams>({
  page: 1,
  pageSize: 10,
})
const { data, isPending, refetch } = useListQRLocationQuery(params)

function handleSortingChange(sorts: SortPrefix<ListQRLocationSort>[]) {
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
        <SheetContent class="sm:max-w-xl">
          <SheetHeader>
            <SheetTitle>Create a QR location</SheetTitle>
            <SheetDescription>
              Create a new QR location to build your map.
              You can add metadata fields to store additional
              information about your QR Location.
            </SheetDescription>
          </SheetHeader>
          <CreateQRLocationForm class="mt-4" />
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

    <QRLocationTable
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
