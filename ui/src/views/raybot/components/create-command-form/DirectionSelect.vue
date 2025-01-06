<script setup lang="ts">
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
</script>

<template>
  <Select
    default-value="FORWARD"
  >
    <SelectTrigger id="direction">
      <SelectValue placeholder="Select a direction" />
    </SelectTrigger>
    <SelectContent position="popper">
      <SelectItem value="FORWARD">
        FORWARD
      </SelectItem>
      <SelectItem value="BACKWARD">
        BACKWARD
      </SelectItem>
    </SelectContent>
  </Select>
</template>
