import type { EdgeTypesObject } from '@vue-flow/core'
import { WorkflowEdge } from './edges'

export const edgeTypes: EdgeTypesObject = {
  WORKFLOW: markRaw(WorkflowEdge),
}
