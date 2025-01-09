import type { CreateQRLocationParams, ListQRLocationParams } from '@/api/qr-location'
import qrLocations from '@/api/qr-location'
import { keepPreviousData, useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

const QUERY_KEY = 'qrLocations'

export function useListQRLocationQuery(p: MaybeRef<ListQRLocationParams>) {
  return useQuery({
    queryKey: [QUERY_KEY, toValue(p)],
    queryFn: () => qrLocations.list(toValue(p)),
    placeholderData: keepPreviousData,
  })
}

export function useCreateQRLocationMutation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (body: CreateQRLocationParams) => qrLocations.create(body),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEY],
      })
    },
  })
}
