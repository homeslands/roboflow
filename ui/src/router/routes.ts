import type { Component } from 'vue'
import {
  BotIcon,
  QrCodeIcon,
  WorkflowIcon,
} from 'lucide-vue-next'

interface Route {
  name: string
  path: string
  icon: Component
}

const routes: Route[] = [
  {
    name: 'Raybots',
    path: '/raybots',
    icon: BotIcon,
  },
  {
    name: 'QR Locations',
    path: '/qr-locations',
    icon: QrCodeIcon,
  },
  {
    name: 'Workflows',
    path: '/workflows',
    icon: WorkflowIcon,
  },
]

export { routes }
