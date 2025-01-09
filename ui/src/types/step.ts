import type { Node } from './workflow'

export const StepStatusValues = [
  'PENDING',
  'RUNNING',
  'COMPLETED',
  'FAILED',
] as const
export type StepStatus = typeof StepStatusValues[number]

export interface Step {
  id: string
  workflowExecutionId: string
  env: {
    [key: string]: string
  }
  node: Node
  status: StepStatus
  startedAt?: Date
  completedAt?: Date
}
