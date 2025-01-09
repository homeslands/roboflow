import type { ListWorkflowParams } from '@/api/workflow'
import workflows from '@/api/workflow'
import { keepPreviousData, useQuery } from '@tanstack/vue-query'

const QUERY_KEY = 'workflows'

export function useListWorkflowQuery(p: Ref<ListWorkflowParams>) {
  return useQuery({
    queryKey: [QUERY_KEY, toValue(p)],
    queryFn: () => workflows.list(toValue(p)),
    placeholderData: keepPreviousData,
  })
}
