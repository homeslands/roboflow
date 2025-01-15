import type { RaybotCommandType } from '@/types/raybot-command'
import type { ActionNode, ActionNodeType } from '@/types/workflow'
import {
  type NodeProps,
  useVueFlow,
  type Edge as VueflowEdge,
  type Node as VueflowNode,
} from '@vue-flow/core'
import { v4 } from 'uuid'

interface AddActionNodeParams {
  nodeType: ActionNodeType
  definition: ActionNode['definition']
}

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

function getLabel({ nodeType, definition }: AddActionNodeParams) {
  if (nodeType === 'RAYBOT_CONTROL') {
    return raybotCommandLabelMap[definition.type]
  }

  return nodeType
}

export function useAddNode(parentNodeId: NodeProps['id']) {
  const { setNodes, setEdges, findNode } = useVueFlow()

  function addActionNode({ nodeType, definition }: AddActionNodeParams) {
    const parentNode = findNode(parentNodeId)
    if (!parentNode)
      return

    // Create a new node
    const nodeId = v4()
    const newNode: VueflowNode = {
      id: nodeId,
      type: nodeType,
      position: {
        x: parentNode.position.x + 300,
        y: parentNode.position.y,
      },
      data: {
        label: getLabel({ nodeType, definition }),
        definition,
      },
    }

    // Create a new edge between the parent node and the new node
    const newEdge: VueflowEdge = {
      id: `${parentNode.id}->${nodeId}`,
      source: parentNode.id,
      target: nodeId,
      type: 'WORKFLOW',
      animated: true,
    }

    setNodes(nodes => [...nodes, newNode])
    setEdges(edges => [...edges, newEdge])
  }

  return {
    addActionNode,
  }
}
