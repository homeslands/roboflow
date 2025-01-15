import type { NodeTypesObject } from '@vue-flow/core'
import { EmptyNode } from './nodes/empty-node'
import { RaybotControlNode } from './nodes/raybot-control-node'
import { TriggerNode } from './nodes/trigger-node'

export const nodeTypes: NodeTypesObject = {
  // @ts-expect-error - this is a valid node type
  EMPTY: markRaw(EmptyNode),
  TRIGGER: markRaw(TriggerNode),
  RAYBOT_CONTROL: markRaw(RaybotControlNode),
}
