import type { SortPrefix } from '@/lib/sort'
import type { Paging } from '@/types/paging'
import type { WorkflowExecution } from '@/types/workflow-execution'
import http from '@/lib/http'

export type ListWorkflowExecutionSort = 'status' | 'created_at' | 'started_at' | 'completed_at'
export interface ListWorkflowExecutionParams {
  page?: number
  pageSize?: number
  sort?: SortPrefix<ListWorkflowExecutionSort>[]
}

const workflowExecutions = {
  list(workflowId: string, p: ListWorkflowExecutionParams): Promise<Paging<WorkflowExecution>> {
    return http.get(`/workflows/${workflowId}/executions`, { params: {
      page: p.page,
      pageSize: p.pageSize,
      sort: p.sort?.join(','),
    },
    })
  },
  get(id: string): Promise<WorkflowExecution> {
    return http.get(`/workflow-executions/${id}`)
  },
}

export default workflowExecutions
