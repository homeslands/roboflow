import type { WorkflowDefinition } from './workflow'

export const WorkflowExecutionStatusValues = [
  'PENDING',
  'RUNNING',
  'COMPLETED',
  'FAILED',
  'CANCELED',
] as const
export type WorkflowExecutionStatus = typeof WorkflowExecutionStatusValues[number]

export interface WorkflowExecution {
  id: string
  workflowId: string
  status: WorkflowExecutionStatus
  env: {
    [key: string]: string
  }
  definition: WorkflowDefinition
  createdAt: Date
  startedAt?: Date
  completedAt?: Date
}
