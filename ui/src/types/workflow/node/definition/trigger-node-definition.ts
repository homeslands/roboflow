export const TriggerTypeValues = [
  'ON_DEMAND',
  'SCHEDULE',
] as const
export type TriggerType = typeof TriggerTypeValues[number]

export interface TriggerNodeDefinition<
  T extends TriggerType = TriggerType,
> {
  type: T
  label: string
  configs: ConfigMap[T]
}

export interface InputConfig {
  key: string
  inputType: 'text' | 'number'
  defaultValue: string
  required: boolean
}

interface ConfigMap {
  ON_DEMAND: InputConfig[]
  SCHEDULE: undefined
}
