import type { ControlRaybotNodeDefinition, ControlRaybotType } from '@/types/workflow/node/definition/control-raybot-node-definition'
import type { TriggerNodeDefinition, TriggerType } from '@/types/workflow/node/definition/trigger-node-definition'

export const TRIGGER_DEFINITION_REGISTRY: Record<TriggerType, TriggerNodeDefinition> = {
  ON_DEMAND: {
    type: 'ON_DEMAND',
    label: 'On Demand',
    configs: [],
  },
  SCHEDULE: {
    type: 'SCHEDULE',
    label: 'Schedule',
    configs: undefined,
  },
}

export const CONTROL_RAYBOT_DEFINITION_REGISTRY: Record<ControlRaybotType, ControlRaybotNodeDefinition> = {
  STOP: {
    type: 'STOP',
    label: 'Stop',
  },
  MOVE_FORWARD: {
    type: 'MOVE_FORWARD',
    label: 'Move Forward',
  },
  MOVE_BACKWARD: {
    type: 'MOVE_BACKWARD',
    label: 'Move Backward',
  },
  MOVE_TO_LOCATION: {
    type: 'MOVE_TO_LOCATION',
    label: 'Move to Location',
  },
  OPEN_BOX: {
    type: 'OPEN_BOX',
    label: 'Open Box',
  },
  CLOSE_BOX: {
    type: 'CLOSE_BOX',
    label: 'Close Box',
  },
  LIFT_BOX: {
    type: 'LIFT_BOX',
    label: 'Lift Box',
  },
  DROP_BOX: {
    type: 'DROP_BOX',
    label: 'Drop Box',
  },
  CHECK_QR: {
    type: 'CHECK_QR',
    label: 'Check QR',
  },
  WAIT_GET_ITEM: {
    type: 'WAIT_GET_ITEM',
    label: 'Wait Get Item',
  },
  SCAN_QR_LOCATION: {
    type: 'SCAN_QR_LOCATION',
    label: 'Scan QR Location',
  },
}
