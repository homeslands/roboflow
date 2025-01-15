import type { Node as VueflowNode } from '@vue-flow/core'
import type { ControlRaybotNodeDefinition } from './control-raybot-node-definition'
import type { TriggerNodeDefinition } from './trigger-node-definition'

export const NodeTypeValues = [
  'EMPTY',
  'TRIGGER',
  'RAYBOT_CONTROL',
] as const
export type NodeType = typeof NodeTypeValues[number]
export type ActionNodeType = Exclude<NodeType, 'EMPTY' | 'TRIGGER'>

export type Node<T extends NodeType = NodeType> = Omit<VueflowNode, 'type'> & {
  type: T
  label: string
  definition: NodeDefinitionMap[T]
}

export type ActionNode<T extends ActionNodeType = ActionNodeType> = Omit<Node<T>, 'type'> & {
  type: T
  definition: NodeDefinitionMap[T]
}

export interface NodeDefinitionMap {
  EMPTY: undefined
  TRIGGER: TriggerNodeDefinition
  RAYBOT_CONTROL: ControlRaybotNodeDefinition
}
