import type { App } from 'vue'
import router from '@/router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { notivue } from './notivue'
import { pinia } from './pinia'
import { queryPluginOpts } from './vue-query'

export function registerPlugins(app: App) {
  app
    .use(router)
    .use(pinia)
    .use(VueQueryPlugin, queryPluginOpts)
    .use(notivue)
}
