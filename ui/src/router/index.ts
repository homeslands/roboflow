import { useNProgress } from '@/lib/nprogress'
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import 'nprogress/nprogress.css'

const MainLayout = () => import('@/layouts/main-layout/MainLayout.vue')

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/raybots',
    meta: {
      name: 'Home',
    },
  },
  {
    path: '/raybots',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'RaybotList',
        component: () => import('@/views/raybot/RaybotList.vue'),
        meta: {
          name: 'Raybots',
        },
      },
      {
        path: ':id',
        name: 'RaybotDetail',
        component: () => import('@/views/raybot/RaybotDetail.vue'),
        meta: {
          name: 'Raybot Detail',
        },
      },
    ],
  },
  {
    path: '/qr-locations',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'QRLocationList',
        component: () => import('@/views/qr-location/QRLocationList.vue'),
        meta: {
          name: 'QR Locations',
        },
      },
    ],
  },
  {
    path: '/workflows',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'Workflows',
        component: () => import('@/views/workflow/WorkflowList.vue'),
        meta: {
          name: 'Workflows',
        },
      },
      {
        path: 'new',
        name: 'NewWorkflow',
        component: () => import('@/views/workflow/WorkflowBuilder.vue'),
        meta: {
          name: 'Workflow Builder',
        },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})
const nprogress = useNProgress()

router.beforeEach((to, _, next) => {
  if (to.meta.name) {
    document.title = `${to.meta.name} | Roboflow`
  }
  nprogress.start()

  next()
})

router.afterEach(() => {
  nprogress.done()
})

export default router
