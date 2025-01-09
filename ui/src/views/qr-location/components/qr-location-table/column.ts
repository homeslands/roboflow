import type { QRLocation } from '@/types/qr-location'
import type { ColumnDef } from '@tanstack/vue-table'
import DataTableSortableHeader from '@/components/DataTableSortableHeader.vue'

export const columns: ColumnDef<QRLocation>[] = [
  {
    accessorKey: 'name',
    header: ({ column }) => h(DataTableSortableHeader<QRLocation>, { column, title: 'Name' }),
    cell: ({ row }) => {
      const name = row.original.name
      return h('span', { class: 'max-w-[500px]' }, name)
    },
  },
  {
    accessorKey: 'qrCode',
    header: ({ column }) => h(DataTableSortableHeader<QRLocation>, { column, title: 'QR code' }),
    cell: ({ row }) => {
      const qrCode = row.original.qrCode
      return h('div', { class: 'flex items-center' }, qrCode)
    },
  },
  {
    accessorKey: 'metadata',
    header: ({ column }) => h(DataTableSortableHeader<QRLocation>, { column, title: 'Metadata' }),
    cell: ({ row }) => {
      const metadata = row.getValue('metadata')
      return h('div', { class: 'flex items-center' }, h('pre', { class: 'whitespace-pre-wrap' }, JSON.stringify(metadata, null, 2)))
    },
    enableSorting: false,
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => h(DataTableSortableHeader<QRLocation>, { column, title: 'Created at' }),
    cell: ({ row }) => {
      const createdAt = row.original.createdAt
      return h('div', { class: 'flex items-center' }, createdAt.toString())
    },
  },
]
