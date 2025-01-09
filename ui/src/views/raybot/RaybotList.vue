<script setup lang="ts">
import type { ListRaybotParams, ListRaybotSort } from '@/api/raybot'
import type { SortPrefix } from '@/lib/sort'
import type { RaybotStatus } from '@/types/raybot'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { useListRaybotQuery } from '@/composables/use-raybot'
import { RaybotStatusValues } from '@/types/raybot'
import { RefreshCcwIcon, SearchIcon, XIcon } from 'lucide-vue-next'
import { columns, RaybotTable } from './components/raybot-table'

type ListRaybotRequiredPagingParams = Omit<ListRaybotParams, 'page' | 'pageSize'> & {
  page: number
  pageSize: number
}

const selectedStatus = ref<RaybotStatus>()
const params = ref<ListRaybotRequiredPagingParams>({
  page: 1,
  pageSize: 10,
})

const { data, isPending, refetch } = useListRaybotQuery(params)

function onStatusChange(status: RaybotStatus) {
  selectedStatus.value = status
}

function handleSearch() {
  params.value.status = selectedStatus.value
  refetch()
}

function handleReset() {
  selectedStatus.value = undefined
  params.value = {
    page: 1,
    pageSize: 10,
    sort: [],
    status: undefined,
  }
  refetch()
}

function handleSortingChange(sorts: SortPrefix<ListRaybotSort>[]) {
  params.value = {
    ...params.value,
    sort: sorts,
  }
  refetch()
}
</script>

<template>
  <div class="space-y-4">
    <!-- Top search config -->
    <div class="flex items-center space-x-4">
      <Select
        v-model="selectedStatus"
        @update:model-value="(status) => onStatusChange(status as RaybotStatus)"
      >
        <SelectTrigger>
          <SelectValue placeholder="Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem
            v-for="status in RaybotStatusValues"
            :key="status"
            :value="status"
          >
            {{ status }}
          </SelectItem>
        </SelectContent>
      </Select>
      <Button
        variant="outline"
        @click="handleSearch"
      >
        <SearchIcon />
        Search
      </Button>
      <Button
        variant="outline"
        @click="handleReset"
      >
        <XIcon />
        Reset
      </Button>
    </div>
    <Separator />
    <div class="flex justify-end">
      <Button
        variant="outline"
        @click="refetch"
      >
        <RefreshCcwIcon />
        Refresh
      </Button>
    </div>

    <!-- Data table -->
    <RaybotTable
      v-model:page="params.page"
      v-model:page-size="params.pageSize"
      :columns="columns"
      :is-loading="isPending"
      :data="data?.items ?? []"
      :total-items="data?.totalItems ?? 0"
      :sorts="params.sort"
      @sorts="handleSortingChange"
    />
  </div>
</template>
