import type { SortPrefix } from '@/lib/sort'
import type { Paging } from '@/types/paging'
import type { Step } from '@/types/step'
import http from '@/lib/http'

export type ListStepSort = 'status' | 'started_at' | 'completed_at'
export interface ListStepParams {
  sort?: SortPrefix<ListStepSort>[]
}

const steps = {
  list(workflowExecutionId: string, p: ListStepParams): Promise<Paging<Step>> {
    return http.get(`/workflow-executions/${workflowExecutionId}/steps`, { params: {
      sort: p.sort?.join(','),
    },
    })
  },
  get(id: string): Promise<Step> {
    return http.get(`/steps/${id}`)
  },
}

export default steps
