export const RaybotCommandTypeValues = [
  'STOP',
  'MOVE_FORWARD',
  'MOVE_BACKWARD',
  'MOVE_TO_LOCATION',
  'OPEN_BOX',
  'CLOSE_BOX',
  'LIFT_BOX',
  'DROP_BOX',
  'CHECK_QR',
  'WAIT_GET_ITEM',
  'SCAN_QR_LOCATION',
] as const
export type RaybotCommandType = typeof RaybotCommandTypeValues[number]

export const RaybotCommandStatusValues = [
  'PENDING',
  'IN_PROGRESS',
  'SUCCESS',
  'FAILED',
] as const
export type RaybotCommandStatus = typeof RaybotCommandStatusValues[number]

export interface InputMap {
  STOP: undefined
  MOVE_FORWARD: undefined
  MOVE_BACKWARD: undefined
  OPEN_BOX: undefined
  CLOSE_BOX: undefined
  WAIT_GET_ITEM: undefined
  DROP_BOX?: { distance?: number }
  LIFT_BOX?: { distance?: number }
  MOVE_TO_LOCATION: {
    location: string
    direction: 'FORWARD' | 'BACKWARD'
  }
  CHECK_QR: { qr_code: string }
  SCAN_QR_LOCATION: undefined
}

export type SuccessOutputMap = Record<RaybotCommandType, undefined> & {
  SCAN_QR_LOCATION: { locations: string[] }
}

export interface OutputMap {
  FAILED: { reason: string }
  SUCCESS: SuccessOutputMap[RaybotCommandType]
  PENDING: undefined
  IN_PROGRESS: undefined
}

export interface RaybotCommand<
  T extends RaybotCommandType = RaybotCommandType,
  S extends RaybotCommandStatus = RaybotCommandStatus,
> {
  id: string
  raybotId: string
  type: T
  status: RaybotCommandStatus
  input: InputMap[T]
  output: OutputMap[S]
  createdAt: Date
  completedAt?: Date
}
