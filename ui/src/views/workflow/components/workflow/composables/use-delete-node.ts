import { useVueFlow } from '@vue-flow/core'

export function useDeleteNode(id: string) {
  /**
   * There are some cases we have to deal with:
   *
   * 1. Deleting a node that has no children (completely delete the node)
   * 2. Deleting a node that has children (delete the node and re-assign the children to the parent of the deleted node)
   * 3. Deleting the root node (TRIGGER node) (turn into empty node)
   * 4. Deleting an EMPTY node (not allowed)
   */
  const { setNodes, setEdges, findNode, edges } = useVueFlow()

  function deleteNode() {
    const node = findNode(id)

    if (!node) {
      console.warn('Node not found')
      return
    }

    // Case 4: Deleting an EMPTY node (not allowed)
    if (node.type === 'EMPTY') {
      console.warn('Deleting an EMPTY node is not allowed')
      return
    }

    // Case 3: Deleting the root node (TRIGGER node)
    if (node.type === 'TRIGGER') {
      setNodes(nodes =>
        nodes.map(n =>
          n.id === id ? { ...n, type: 'EMPTY', data: {}, label: '' } : n,
        ),
      )
      return
    }

    const childEdges = edges.value.filter(edge => edge.source === id)
    const parentEdges = edges.value.filter(edge => edge.target === id)

    // Case 1: Deleting a node that has no children
    if (childEdges.length === 0) {
      setNodes(nodes => nodes.filter(n => n.id !== id))
      setEdges(edges => edges.filter(edge => edge.source !== id && edge.target !== id))
      return
    }

    // Case 2: Deleting a node that has children
    setEdges((edges) => {
      const newEdges = edges.filter(edge => edge.source !== id && edge.target !== id)
      if (parentEdges.length > 0) {
        const parent = parentEdges[0].source // Assume one parent
        childEdges.forEach((childEdge) => {
          newEdges.push({ ...childEdge, source: parent })
        })
      }
      return newEdges
    })

    setNodes(nodes => nodes.filter(n => n.id !== id))
  }

  return { deleteNode }
}
