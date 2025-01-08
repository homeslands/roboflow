import type { SortPrefix } from '@/lib/sort'
import type { Paging } from '@/types/paging'
import type { InputMap, RaybotCommand, RaybotCommandType } from '@/types/raybot-command'
import http from '@/lib/http'

export type ListRaybotCommandSort = 'status' | 'created_at' | 'completed_at'
export interface ListRaybotCommandParams {
  page?: number
  pageSize?: number
  sort?: SortPrefix<ListRaybotCommandSort>[]
}

export interface CreateRaybotCommandParams<T extends RaybotCommandType> {
  type: T
  input?: InputMap[T]
}

const raybotCommands = {
  list(raybotId: string, p: ListRaybotCommandParams): Promise<Paging<RaybotCommand>> {
    return http.get(`/raybots/${raybotId}/commands`, { params: {
      page: p.page,
      pageSize: p.pageSize,
      sort: p.sort?.join(','),
    },
    })
  },
  get(id: string): Promise<RaybotCommand> {
    return http.get(`/raybot-commands/${id}`)
  },
  create<T extends RaybotCommandType>(
    raybotId: string,
    body: CreateRaybotCommandParams<T>,
  ): Promise<RaybotCommand<T>> {
    return http.post(`/raybots/${raybotId}/commands`, body)
  },
}

export default raybotCommands
