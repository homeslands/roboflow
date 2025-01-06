import type { Raybot } from '@/types/raybot'
import type { ColumnDef } from '@tanstack/vue-table'
import RaybotNameBtn from './RaybotNameBtn.vue'

export const columns: ColumnDef<Raybot>[] = [
  {
    accessorKey: 'name',
    header: () => h('span', { class: 'text-right' }, 'Name'),
    cell: ({ row }) => {
      const raybot = row.original
      return h('span', { class: 'text-right' }, h(RaybotNameBtn, { id: raybot.id, name: raybot.name }))
    },
  },
  {
    accessorKey: 'status',
    header: () => h('span', { class: 'text-right' }, 'Status'),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('status'))
    },
  },
  {
    accessorKey: 'ipAddress',
    header: () => h('span', { class: 'text-right' }, 'IP Address'),
    cell: ({ row }) => {
      const ipAddress = row.original.ipAddress ?? 'N/A'
      return h('span', { class: 'text-right' }, ipAddress)
    },
  },
  {
    accessorKey: 'lastConnectedAt',
    header: () => h('span', { class: 'text-right' }, 'Last Connected At'),
    cell: ({ row }) => {
      const lastConnectedAt = row.original.lastConnectedAt?.toLocaleString() ?? 'N/A'
      return h('span', { class: 'text-right' }, lastConnectedAt)
    },
  },
  {
    accessorKey: 'createdAt',
    header: () => h('span', { class: 'text-right' }, 'Created At'),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('createdAt'))
    },
  },
]
