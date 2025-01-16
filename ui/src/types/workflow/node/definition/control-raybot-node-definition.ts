import type { RaybotCommandType } from '@/types/raybot-command'

export type ControlRaybotType = RaybotCommandType

export interface ControlRaybotNodeDefinition {
  type: ControlRaybotType
  label: string
}
