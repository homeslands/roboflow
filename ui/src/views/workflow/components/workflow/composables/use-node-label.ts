import { type NodeProps, useVueFlow } from '@vue-flow/core'

export function useNodeLabel(id: NodeProps['id']) {
  const { nodes, updateNodeData } = useVueFlow()
  const node = nodes.value.find(node => node.id === id)

  function setLabel(newLabel: string) {
    if (!node)
      return

    updateNodeData(id, { label: newLabel })
  }

  return {
    label: node?.data.label as string ?? '',
    setLabel,
  }
}
