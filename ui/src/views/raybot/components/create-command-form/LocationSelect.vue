<script setup lang="ts">
import type { ListQRLocationParams } from '@/api/qr-location'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useQRLocationQuery } from '@/composables/use-qr-location'

const selectedQRLocation = defineModel<string>()

const params: ListQRLocationParams = {
  pageSize: 10000,
}
const { data, isPending } = useQRLocationQuery(params)
</script>

<template>
  <template v-if="isPending">
    Loading...
  </template>
  <template v-else-if="data?.items.length === 0 || !data">
    No locations found.
  </template>
  <template v-else>
    <Select
      v-model="selectedQRLocation"
    >
      <SelectTrigger>
        <SelectValue placeholder="Select a location (name - qr code)" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem
          v-for="location in data.items"
          :key="location.id"
          :value="location.qrCode"
        >
          {{ location.name }} - {{ location.qrCode }}
        </SelectItem>
      </SelectContent>
    </Select>
  </template>
</template>
