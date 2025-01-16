import type { Node, NodeType } from '@/types/workflow'
import type { Component } from 'vue'
import { BoxIcon, MoveIcon, ShapesIcon } from 'lucide-vue-next'
import { CONTROL_RAYBOT_DEFINITION_REGISTRY } from '../node-definition-registry'

type ActionNodeType = Exclude<NodeType, 'EMPTY' | 'TRIGGER'>
export interface SubMenuItem<T extends ActionNodeType = ActionNodeType> {
  label: string
  nodeType: T
  definition: Node['definition']
}

interface MenuItem {
  label: string
  subMenu: {
    label: string
    icon: Component
    subItems: SubMenuItem[]
  }[]
}

export const nodeMenuItems: MenuItem[]
= [
  {
    label: 'Raybot Control',
    subMenu: [
      {
        label: 'Movement',
        icon: MoveIcon,
        subItems: [
          {
            label: 'Move Forward',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.MOVE_FORWARD,
          },
          {
            label: 'Move Backward',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.MOVE_BACKWARD,
          },
          {
            label: 'Move To Location',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.MOVE_TO_LOCATION,
          },
        ],
      },
      {
        label: 'Container Box',
        icon: BoxIcon,
        subItems: [
          {
            label: 'Open Box',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.OPEN_BOX,
          },
          {
            label: 'Close Box',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.CLOSE_BOX,
          },
          {
            label: 'Lift Box',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.LIFT_BOX,
          },
          {
            label: 'Drop Box',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.DROP_BOX,
          },
          {
            label: 'Wait Get Item',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.WAIT_GET_ITEM,
          },
        ],
      },
      {
        label: 'Another',
        icon: ShapesIcon,
        subItems: [
          {
            label: 'Check QR Code',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.CHECK_QR,
          },
          {
            label: 'Scan QR Location',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.SCAN_QR_LOCATION,
          },
          {
            label: 'Stop',
            nodeType: 'CONTROL_RAYBOT',
            definition: CONTROL_RAYBOT_DEFINITION_REGISTRY.STOP,
          },
        ],
      },
    ],
  },
]
