export const RaybotStatusValues = [
  'OFFLINE',
  'IDLE',
  'BUSY',
] as const
export type RaybotStatus = typeof RaybotStatusValues[number]

export interface Raybot {
  id: string
  name: string
  token: string
  status: RaybotStatus
  ipAddress?: string
  lastConnectedAt?: Date
  createdAt: Date
  updatedAt: Date
}
