import type { Paging } from '@/types/paging'
import type { Raybot, RaybotStatus } from '@/types/raybot'
import http from '@/lib/http'

export type ListRaybotSort = 'name' | 'status' | 'last_connected_at' | 'created_at' | 'updated_at'
export interface ListRaybotParams {
  page?: number
  pageSize?: number
  sort?: Array<`-${ListRaybotSort}` | ListRaybotSort>
  status?: RaybotStatus
}

export interface CreateRaybotParams {
  name: string
}

const raybots = {
  list(p: ListRaybotParams): Promise<Paging<Raybot>> {
    return http.get('/raybots', { params: {
      page: p.page,
      pageSize: p.pageSize,
      sort: p.sort?.join(','),
      status: p.status,
    },
    })
  },
  get(id: string): Promise<Raybot> {
    return http.get(`/raybots/${id}`)
  },
  create(body: CreateRaybotParams): Promise<Raybot> {
    return http.post('/raybots', body)
  },
  delete(id: string): Promise<void> {
    return http.delete(`/raybots/${id}`)
  },
}

export default raybots
