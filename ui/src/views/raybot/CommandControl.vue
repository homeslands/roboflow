<script setup lang="ts">
import type { ListRaybotCommandParams } from '@/api/raybot-command'
import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import { useRaybotCommandQuery } from '@/composables/use-raybot-command'
import { PlusIcon, RefreshCcwIcon } from 'lucide-vue-next'
import { columns, CommandTable } from './components/command-table'
import { CreateCommandForm } from './components/create-command-form'

const route = useRoute()
const raybotId = computed(() => route.params.id as string)
const params = ref<ListRaybotCommandParams>({
  sort: ['-created_at'],
})
const { data, isPending, isRefetching, refetch } = useRaybotCommandQuery(raybotId, params)
</script>

<template>
  <div class="flex flex-col items-end space-y-2">
    <div class="space-x-2">
      <Sheet>
        <SheetTrigger>
          <Button variant="outline">
            <PlusIcon />
            Create
          </Button>
        </SheetTrigger>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Create new command</SheetTitle>
            <SheetDescription>
              Create a new command to send to raybot.
              Raybot can only process <strong>one</strong> command at a time.
            </SheetDescription>
          </SheetHeader>
          <CreateCommandForm class="mt-4" />
        </SheetContent>
      </Sheet>

      <Button
        variant="outline"
        @click="refetch"
      >
        <RefreshCcwIcon />
        Refresh
      </Button>
    </div>

    <!-- Data table -->
    <CommandTable
      class="w-full"
      :columns="columns"
      :data="data?.items"
      :is-loading="isPending || isRefetching"
    />
  </div>
</template>
