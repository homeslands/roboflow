import type { CreateRaybotCommandParams, ListRaybotCommandParams } from '@/api/raybot-command'
import type { RaybotCommandType } from '@/types/raybot-command'
import raybotCommands from '@/api/raybot-command'
import { keepPreviousData, useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

const QUERY_KEY = 'raybotCommands'

export function useListRaybotCommandQuery(raybotId: MaybeRef<string>, p: MaybeRef<ListRaybotCommandParams>) {
  return useQuery({
    queryKey: [QUERY_KEY, toValue(p)],
    queryFn: () => raybotCommands.list(toValue(raybotId), toValue(p)),
    placeholderData: keepPreviousData,
  })
}

export function useCreateRaybotCommandMutation(raybotId: MaybeRef<string>) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (body: CreateRaybotCommandParams<RaybotCommandType>) =>
      raybotCommands.create(toValue(raybotId), body),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEY],
      })
    },
  })
}
