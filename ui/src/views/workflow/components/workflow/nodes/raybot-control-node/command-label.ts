import type { RaybotCommandType } from '@/types/raybot-command'

type RaybotCommandLabelMap = {
  [key in RaybotCommandType]: string;
}

export const raybotCommandLabelMap: RaybotCommandLabelMap = {
  STOP: 'Stop',
  MOVE_FORWARD: 'Move To Location',
  MOVE_BACKWARD: 'Move Backward',
  MOVE_TO_LOCATION: 'Move To Location',
  OPEN_BOX: 'Open Box',
  CLOSE_BOX: 'Close Box',
  LIFT_BOX: 'Lift Box',
  DROP_BOX: 'Drop Box',
  CHECK_QR: 'Check QR Code',
  WAIT_GET_ITEM: 'Wait Get Item',
  SCAN_QR_LOCATION: 'Scan QR Location',
}
