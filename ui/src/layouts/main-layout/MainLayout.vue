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
import UserDropdown from './UserDropdown.vue'
</script>

<template>
  <SidebarProvider>
    <Sidebar collapsible="icon">
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
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
    <SidebarInset>
      <header
        class="flex justify-between h-16 shrink-0 items-center gap-2 px-4 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12"
      >
        <SidebarTrigger class="-ml-1" />
        <UserDropdown />
      </header>
      <!-- Main content -->
      <div class="flex flex-col flex-1 gap-4 p-4 pt-0">
        <RouterView />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
