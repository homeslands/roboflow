<script setup lang="ts" generic="TData, TValue">
import type { ColumnDef } from '@tanstack/vue-table'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  FlexRender,
  getCoreRowModel,
  useVueTable,
} from '@tanstack/vue-table'
import { LoaderCircleIcon } from 'lucide-vue-next'
import DataTablePagination from './DataTablePagination.vue'

interface Props {
  isLoading: boolean
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  totalItems: number
  pageSizeOptions?: number[]
}
const props = withDefaults(defineProps<Props>(), {
  pageSizeOptions: () => [10, 20, 50, 100],
})
const page = defineModel<number>('page', { required: true })
const pageSize = defineModel<number>('pageSize', { required: true })

const pageSizeStr = computed({
  get: () => pageSize.value.toString(),
  set: value => pageSize.value = Number(value),
})

const table = useVueTable({
  get data() { return props.data },
  get columns() { return props.columns },
  getCoreRowModel: getCoreRowModel(),

  // Server-side pagination
  manualPagination: true,
  rowCount: props.totalItems,
  state: {
    pagination: {
      pageIndex: page.value - 1,
      pageSize: pageSize.value,
    },
  },
  onPaginationChange: (updater) => {
    const newState = typeof updater === 'function'
      ? updater(table.getState().pagination)
      : updater

    page.value = newState.pageIndex
    pageSize.value = newState.pageSize
  },
})
</script>

<template>
  <div class="flex flex-col space-y-4">
    <div class="border rounded-md">
      <Table>
        <TableHeader class="bg-muted-foreground/10">
          <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
            <TableHead v-for="header in headerGroup.headers" :key="header.id">
              <FlexRender
                v-if="!header.isPlaceholder"
                :render="header.column.columnDef.header"
                :props="header.getContext()"
              />
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="props.isLoading">
            <TableRow>
              <TableCell
                :colspan="columns.length"
                class="relative h-24"
              >
                <div class="absolute inset-0 flex items-center justify-center">
                  <LoaderCircleIcon
                    class="w-8 h-8 animate-spin text-primary"
                    aria-label="Loading..."
                  />
                </div>
              </TableCell>
            </TableRow>
          </template>
          <template v-else-if="table.getRowModel().rows?.length">
            <TableRow
              v-for="row in table.getRowModel().rows" :key="row.id"
              :data-state="row.getIsSelected() ? 'selected' : undefined"
            >
              <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
              </TableCell>
            </TableRow>
          </template>
          <template v-else>
            <TableRow>
              <TableCell :colspan="columns.length" class="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          </template>
        </TableBody>
      </Table>
    </div>

    <div class="flex ml-auto space-x-6">
      <!-- Page size selector -->
      <div class="flex items-center gap-1">
        <span class="text-sm sr-only text-muted-foreground sm:not-sr-only">
          Items per page:
        </span>
        <Select
          v-model="pageSizeStr"
          :disabled="props.isLoading"
        >
          <SelectTrigger class="w-20">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem
                v-for="option in props.pageSizeOptions"
                :key="option"
                :value="option.toString()"
              >
                {{ option }}
              </SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>

      <!-- Pagination -->
      <DataTablePagination
        v-if="!props.isLoading"
        :current-page="page"
        :page-size="pageSize"
        :total-items="props.totalItems ?? 0"
        @page-change="page = $event"
      />
    </div>
  </div>
</template>
