import type { Node as VueflowNode } from '@vue-flow/core'
import type { ControlRaybotNodeDefinition, EmptyNodeDefinition, TriggerNodeDefinition } from './definition'

export const NodeTypeValues = [
  'EMPTY',
  'TRIGGER',
  'CONTROL_RAYBOT',
] as const
export type NodeType = typeof NodeTypeValues[number]

export type Node<T extends NodeType = NodeType> = Omit<VueflowNode, 'type'> & {
  type: T
  definition: NodeDefinitionMap[T]
}

export interface NodeDefinitionMap {
  EMPTY: EmptyNodeDefinition
  TRIGGER: TriggerNodeDefinition
  CONTROL_RAYBOT: ControlRaybotNodeDefinition
}
