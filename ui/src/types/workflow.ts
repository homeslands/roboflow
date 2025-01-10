import type { FlowExportObject, Edge as VueflowEdge, Node as VueflowNode } from '@vue-flow/core'

export const TaskTypeValues = [
  'VALIDATE_STATE',
  'MOVE_FORWARD',
  'MOVE_BACKWARD',
  'MOVE_TO_LOCATION',
  'OPEN_BOX',
  'CLOSE_BOX',
  'LIFT_BOX',
  'DROP_BOX',
  'CHECK_QR',
  'WAIT_GET_ITEM',
] as const
export type TaskType = typeof TaskTypeValues[number]

export type NodeType = 'RaybotNode'

export interface NodeField {
  useEnv: boolean
  key?: string
  value?: string
}

export interface NodeDefinition {
  type: TaskType
  fields: {
    [key: string]: NodeField
  }
  timeoutSec: number
}

export type Node = Omit<VueflowNode, 'type'> & {
  type: NodeType
  definition: NodeDefinition
}

export type Edge = VueflowEdge

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
