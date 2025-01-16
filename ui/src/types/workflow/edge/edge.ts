import type { Edge as VueflowEdge } from '@vue-flow/core'

export const EdgeTypeValues = [
  'WORKFLOW',
] as const
export type EdgeType = typeof EdgeTypeValues[number]

export type Edge = Omit<VueflowEdge, 'type'> & {
  type: EdgeType
}
