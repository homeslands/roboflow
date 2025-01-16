import type { Node, NodeType } from '@/types/workflow'
import type { ControlRaybotType } from '@/types/workflow/node/definition/control-raybot-node-definition'
import type { TriggerType } from '@/types/workflow/node/definition/trigger-node-definition'
import { useVueFlow, type Edge as VueflowEdge, type Node as VueflowNode } from '@vue-flow/core'
import { v4 } from 'uuid'
import { CONTROL_RAYBOT_DEFINITION_REGISTRY, TRIGGER_DEFINITION_REGISTRY } from '../nodes/node-definition-registry'

export function useAddNode(parentNodeId: string) {
  const { setNodes, setEdges, findNode } = useVueFlow()

  function addTriggerNode(type: TriggerType) {
    const definition = TRIGGER_DEFINITION_REGISTRY[type]
    createNodeAndEdge({ nodeType: 'TRIGGER', definition })
  }

  function addControlRaybotNode(type: ControlRaybotType) {
    const definition = CONTROL_RAYBOT_DEFINITION_REGISTRY[type]
    createNodeAndEdge({ nodeType: 'CONTROL_RAYBOT', definition })
  }

  function createNodeAndEdge({ nodeType, definition }: {
    nodeType: NodeType
    definition: Node['definition']
  }) {
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
    addTriggerNode,
    addControlRaybotNode,
  }
}
