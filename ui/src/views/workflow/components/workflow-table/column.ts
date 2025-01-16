import type { WorkflowWithoutDefinition } from '@/types/workflow'
import type { ColumnDef } from '@tanstack/vue-table'
import DataTableSortableHeader from '@/components/DataTableSortableHeader.vue'
import WorkflowNameLink from './WorkflowNameLink.vue'

export const columns: ColumnDef<WorkflowWithoutDefinition>[] = [
  {
    accessorKey: 'name',
    header: ({ column }) => h(DataTableSortableHeader<WorkflowWithoutDefinition>, { column, title: 'Name' }),
    cell: ({ row }) => {
      const workflow = row.original
      return h('div', { class: 'flex items-center' }, h(WorkflowNameLink, { id: workflow.id, name: workflow.name }))
    },
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => h(DataTableSortableHeader<WorkflowWithoutDefinition>, { column, title: 'Created At' }),
    cell: ({ row }) => {
      const createdAt = row.original.createdAt
      return h('div', { class: 'flex items-center' }, createdAt.toString())
    },
  },
  {
    accessorKey: 'updatedAt',
    header: ({ column }) => h(DataTableSortableHeader<WorkflowWithoutDefinition>, { column, title: 'Updated At' }),
    cell: ({ row }) => {
      const updatedAt = row.original.updatedAt
      return h('div', { class: 'flex items-center' }, updatedAt.toString())
    },
  },
]
