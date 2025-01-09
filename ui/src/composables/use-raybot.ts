import type { ListRaybotParams } from '@/api/raybot'
import raybots from '@/api/raybot'
import { keepPreviousData, useQuery } from '@tanstack/vue-query'

const QUERY_KEY = 'raybots'

export function useRaybotQuery(p: Ref<ListRaybotParams>) {
  return useQuery({
    queryKey: [QUERY_KEY, toValue(p)],
    queryFn: () => raybots.list(toValue(p)),
    placeholderData: keepPreviousData,
  })
}
