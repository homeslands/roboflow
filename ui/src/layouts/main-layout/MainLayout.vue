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
import Logo from './Logo.vue'
</script>

<template>
  <SidebarProvider>
    <Sidebar collapsible="icon" class="bg-popover">
      <SidebarHeader class="flex-row items-center mt-3">
        <Logo />
      </SidebarHeader>
      <SidebarContent class="mx-2 mt-6">
        <SidebarMenu>
          <SidebarMenuItem v-for="route in routes" :key="route.name">
            <SidebarMenuButton as-child>
              <RouterLink
                :to="route.path"
                active-class="text-primary bg-muted"
                class="text-xs font-medium transition-colors rounded-lg text-muted-foreground hover:bg-accent hover:text-accent-foreground"
              >
                <component :is="route.icon" class="w-4 h-4" />
                {{ route.name }}
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
