import type { Paging } from '@/types/paging'
import type { Workflow, WorkflowDefinition, WorkflowWithoutDefinition } from '@/types/workflow'
import http from '@/lib/http'

export type ListWorkflowSort = 'name' | 'created_at' | 'updated_at'
export interface ListWorkflowParams {
  page?: number
  pageSize?: number
  sort?: Array<`-${ListWorkflowSort}` | ListWorkflowSort>
}

export interface CreateRaybotParams {
  name: string
  description?: string
  definition: WorkflowDefinition
}

export type UpdateRaybotParams = CreateRaybotParams

export interface RunWorkflowParams {
  env: {
    [key: string]: string
  }
}

const workflows = {
  list(p: ListWorkflowParams): Promise<Paging<WorkflowWithoutDefinition>> {
    return http.get('/workflows', { params: {
      page: p.page,
      pageSize: p.pageSize,
      sort: p.sort?.join(','),
    },
    })
  },
  get(id: string): Promise<Workflow> {
    return http.get(`/workflows/${id}`)
  },
  create(body: CreateRaybotParams): Promise<Workflow> {
    return http.post('/workflows', body)
  },
  update(id: string, body: UpdateRaybotParams): Promise<Workflow> {
    return http.put(`/workflows/${id}`, body)
  },
  delete(id: string): Promise<void> {
    return http.delete(`/workflows/${id}`)
  },
  run(id: string, body: RunWorkflowParams): Promise<void> {
    return http.post(`/workflows/${id}/run`, body)
  },
}

export default workflows
