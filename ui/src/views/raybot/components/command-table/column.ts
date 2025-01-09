import type { RaybotCommand } from '@/types/raybot-command'
import type { ColumnDef } from '@tanstack/vue-table'
import DataTableSortableHeader from '@/components/DataTableSortableHeader.vue'

export const columns: ColumnDef<RaybotCommand>[] = [
  {
    accessorKey: 'type',
    enableSorting: false,
    header: ({ column }) => h(DataTableSortableHeader<RaybotCommand>, { column, title: 'Type' }),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('type'))
    },
  },
  {
    accessorKey: 'status',
    header: ({ column }) => h(DataTableSortableHeader<RaybotCommand>, { column, title: 'Status' }),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('status'))
    },
  },
  {
    accessorKey: 'input',
    enableSorting: false,
    header: ({ column }) => h(DataTableSortableHeader<RaybotCommand>, { column, title: 'Input' }),
    cell: ({ row }) => {
      const input = row.getValue('input')
      return h('span', { class: 'text-right' }, JSON.stringify(input))
    },
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => h(DataTableSortableHeader<RaybotCommand>, { column, title: 'Created at' }),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('createdAt'))
    },
  },
  {
    accessorKey: 'completedAt',
    header: ({ column }) => h(DataTableSortableHeader<RaybotCommand>, { column, title: 'Completed at' }),
    cell: ({ row }) => {
      const completedAt = row.original.completedAt?.toLocaleString() ?? 'N/A'
      return h('span', { class: 'text-right' }, completedAt)
    },
  },
]
