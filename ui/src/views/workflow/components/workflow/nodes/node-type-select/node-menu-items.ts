import type { RaybotCommandType } from '@/types/raybot-command'
import type { ActionNodeType } from '@/types/workflow'
import type { Component } from 'vue'
import { BoxIcon, MoveIcon, ShapesIcon } from 'lucide-vue-next'

export interface SubMenuItem<T extends ActionNodeType = ActionNodeType> {
  label: string
  nodeType: T
  subType: NodeSubTypeMap[T]
}

interface NodeSubTypeMap {
  RAYBOT_CONTROL: RaybotCommandType
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
          { label: 'Move Forward', nodeType: 'RAYBOT_CONTROL', subType: 'MOVE_FORWARD' },
          { label: 'Move Backward', nodeType: 'RAYBOT_CONTROL', subType: 'MOVE_BACKWARD' },
          { label: 'Move To Location', nodeType: 'RAYBOT_CONTROL', subType: 'MOVE_TO_LOCATION' },
        ],
      },
      {
        label: 'Container Box',
        icon: BoxIcon,
        subItems: [
          { label: 'Open Box', nodeType: 'RAYBOT_CONTROL', subType: 'OPEN_BOX' },
          { label: 'Close Box', nodeType: 'RAYBOT_CONTROL', subType: 'CLOSE_BOX' },
          { label: 'Lift Box', nodeType: 'RAYBOT_CONTROL', subType: 'LIFT_BOX' },
          { label: 'Drop Box', nodeType: 'RAYBOT_CONTROL', subType: 'DROP_BOX' },
          { label: 'Wait Get Item', nodeType: 'RAYBOT_CONTROL', subType: 'WAIT_GET_ITEM' },
        ],
      },
      {
        label: 'Another',
        icon: ShapesIcon,
        subItems: [
          { label: 'Check QR Code', nodeType: 'RAYBOT_CONTROL', subType: 'CHECK_QR' },
          { label: 'Scan QR Location', nodeType: 'RAYBOT_CONTROL', subType: 'SCAN_QR_LOCATION' },
          { label: 'Stop', nodeType: 'RAYBOT_CONTROL', subType: 'STOP' },
        ],
      },
    ],
  },
]
