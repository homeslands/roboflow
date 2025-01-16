import type { FlowExportObject } from '@vue-flow/core'
import type { Edge } from './edge'
import type { Node } from './node'

export type WorkflowDefinition = Omit<FlowExportObject, 'nodes' | 'edges'> & {
  nodes: Node[]
  edges: Edge[]
}

export interface WorkflowWithoutDefinition {
  id: string
  name: string
  description?: string
  createdAt: Date
  updatedAt: Date
}

export interface Workflow extends WorkflowWithoutDefinition {
  definition: WorkflowDefinition
}
