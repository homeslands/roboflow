import type { Paging } from '@/types/paging'
import type { QRLocation } from '@/types/qr-location'
import http from '@/lib/http'

export type ListQRLocationSort = 'name' | 'qr_code' | 'created_at' | 'updated_at'
export interface ListQRLocationParams {
  page?: number
  pageSize?: number
  sort?: Array<`-${ListQRLocationSort}` | ListQRLocationSort>
}

export interface CreateQRLocationParams {
  name: string
  qrCode: string
  metadata: Record<string, string | number | boolean>
}

export type UpdateQRLocationParams = CreateQRLocationParams

const qrLocations = {
  list(p: ListQRLocationParams): Promise<Paging<QRLocation>> {
    return http.get('/qr-locations', { params: {
      page: p.page,
      pageSize: p.pageSize,
      sort: p.sort?.join(','),
    } })
  },
  get(id: string): Promise<QRLocation> {
    return http.get(`/qr-locations/${id}`)
  },
  create(body: CreateQRLocationParams): Promise<QRLocation> {
    return http.post('/qr-locations', body)
  },
  update(id: string, body: UpdateQRLocationParams): Promise<QRLocation> {
    return http.put(`/qr-locations/${id}`, body)
  },
  delete(id: string): Promise<void> {
    return http.delete(`/qr-locations/${id}`)
  },
}

export default qrLocations
