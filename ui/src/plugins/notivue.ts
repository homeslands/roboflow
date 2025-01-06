import { createNotivue } from 'notivue'
import 'notivue/animations.css'
import 'notivue/notification.css'

export const notivue = createNotivue({
  notifications: {
    global: {
      duration: 5000,
    },
  },
})
