import type { RaybotCommand } from '@/types/raybot-command'
import type { ColumnDef } from '@tanstack/vue-table'

export const columns: ColumnDef<RaybotCommand>[] = [
  {
    accessorKey: 'type',
    header: () => h('span', { class: 'text-right' }, 'Type'),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('type'))
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
    accessorKey: 'input',
    header: () => h('span', { class: 'text-right' }, 'Input'),
    cell: ({ row }) => {
      const input = row.getValue('input')
      return h('span', { class: 'text-right' }, JSON.stringify(input))
    },
  },
  {
    accessorKey: 'createdAt',
    header: () => h('span', { class: 'text-right' }, 'Created At'),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('createdAt'))
    },
  },
  {
    accessorKey: 'completedAt',
    header: () => h('span', { class: 'text-right' }, 'Completed At'),
    cell: ({ row }) => {
      const completedAt = row.original.completedAt?.toLocaleString() ?? 'N/A'
      return h('span', { class: 'text-right' }, completedAt)
    },
  },
]
