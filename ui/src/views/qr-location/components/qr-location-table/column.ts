import type { QRLocation } from '@/types/qr-location'
import type { ColumnDef } from '@tanstack/vue-table'

export const columns: ColumnDef<QRLocation>[] = [
  {
    accessorKey: 'name',
    header: () => h('span', { class: 'text-right' }, 'Name'),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('name'))
    },
  },
  {
    accessorKey: 'qrCode',
    header: () => h('span', { class: 'text-right' }, 'QR Code'),
    cell: ({ row }) => {
      return h('span', { class: 'text-right' }, row.getValue('qrCode'))
    },
  },
  {
    accessorKey: 'metadata',
    header: () => h('span', { class: 'text-right' }, 'Metadata'),
    cell: ({ row }) => {
      const metadata = row.getValue('metadata')
      return h('span', { class: 'text-right' }, JSON.stringify(metadata))
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
