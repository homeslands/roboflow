import type { Raybot } from '@/types/raybot'
import type { ColumnDef } from '@tanstack/vue-table'
import DataTableSortableHeader from '@/components/DataTableSortableHeader.vue'
import RaybotNameLink from './RaybotNameLink.vue'

export const columns: ColumnDef<Raybot>[] = [
  {
    accessorKey: 'name',
    header: ({ column }) => h(DataTableSortableHeader<Raybot>, { column, title: 'Name' }),
    cell: ({ row }) => {
      const raybot = row.original
      return h('div', { class: 'flex items-center' }, h(RaybotNameLink, { id: raybot.id, name: raybot.name }))
    },
  },
  {
    accessorKey: 'status',
    header: ({ column }) => h(DataTableSortableHeader<Raybot>, { column, title: 'Status' }),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('status'))
    },
  },
  {
    accessorKey: 'ipAddress',
    enableSorting: false,
    header: ({ column }) => h(DataTableSortableHeader<Raybot>, { column, title: 'IP address' }),
    cell: ({ row }) => {
      const ipAddress = row.original.ipAddress ?? 'N/A'
      return h('span', { class: 'text-right' }, ipAddress)
    },
  },
  {
    accessorKey: 'lastConnectedAt',
    header: ({ column }) => h(DataTableSortableHeader<Raybot>, { column, title: 'Last connected at' }),
    cell: ({ row }) => {
      const lastConnectedAt = row.original.lastConnectedAt?.toLocaleString() ?? 'N/A'
      return h('span', { class: 'text-right' }, lastConnectedAt)
    },
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => h(DataTableSortableHeader<Raybot>, { column, title: 'Created at' }),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('createdAt'))
    },
  },
]
