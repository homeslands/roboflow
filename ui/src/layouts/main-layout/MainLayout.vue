<script setup lang="ts">
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarRail,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { routes } from '@/router/routes'
</script>

<template>
  <SidebarProvider>
    <Sidebar collapsible="icon" class="bg-popover">
      <SidebarHeader class="flex-row items-center mt-3">
        <img alt="Vue logo" src="@/assets/logo.svg" width="32" height="32">
        <span class="ml-2 font-bold">Roboflow</span>
      </SidebarHeader>
      <SidebarContent class="mx-2 mt-6">
        <SidebarMenu>
          <SidebarMenuItem v-for="route in routes" :key="route.name">
            <SidebarMenuButton as-child>
              <RouterLink
                :to="route.path" class="rounded-lg"
                active-class="text-primary bg-muted"
              >
                <component :is="route.icon" />
                <span>{{ route.name }}</span>
              </RouterLink>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
        <SidebarTrigger class="absolute -right-4 bottom-8" />
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
    <SidebarInset>
      <div class="flex-1 w-full h-full p-4">
        <RouterView />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
