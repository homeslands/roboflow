import type { NodeTypesObject } from '@vue-flow/core'
import { ControlRaybotNode } from './nodes/control-raybot-node'
import { EmptyNode } from './nodes/empty-node'
import { TriggerNode } from './nodes/trigger-node'

export const nodeTypes: NodeTypesObject = {
  // @ts-expect-error - this is a valid node type
  EMPTY: markRaw(EmptyNode),
  TRIGGER: markRaw(TriggerNode),
  CONTROL_RAYBOT: markRaw(ControlRaybotNode),
}
