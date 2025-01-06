<script setup lang="ts">
import type { ListRaybotParams } from '@/api/raybot'
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
import { useRaybotQuery } from '@/composables/use-raybot'
import { RaybotStatusValues } from '@/types/raybot'
import { RefreshCcwIcon, SearchIcon, XIcon } from 'lucide-vue-next'
import { columns, RaybotTable } from './components/raybot-table'

const selectedStatus = ref<RaybotStatus>()
const params = ref<ListRaybotParams>({
  status: selectedStatus.value,
})

const { data, isPending, refetch } = useRaybotQuery(params)

function onStatusChange(status: RaybotStatus) {
  selectedStatus.value = status
}
</script>

<template>
  <div class="space-y-4">
    <!-- Top -->
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
        @click="refetch"
      >
        <SearchIcon />
        Search
      </Button>
      <Button variant="outline">
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
      :columns="columns"
      :data="data?.items"
      :is-loading="isPending"
    />
  </div>
</template>
