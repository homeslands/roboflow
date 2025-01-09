<script setup lang="ts">
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import CommandControl from './CommandControl.vue'
</script>

<template>
  <Tabs default-value="control">
    <TabsList>
      <TabsTrigger value="info">
        Raybot Info
      </TabsTrigger>
      <TabsTrigger value="control">
        Raybot Control Panel
      </TabsTrigger>
    </TabsList>
    <TabsContent value="info">
      Make changes to your account here.
    </TabsContent>
    <TabsContent value="control">
      <CommandControl />
    </TabsContent>
  </Tabs>
</template>
