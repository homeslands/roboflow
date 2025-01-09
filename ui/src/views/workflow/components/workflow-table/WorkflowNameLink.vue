<script setup lang="ts">
defineProps<{
  id: string
  name: string
}>()
</script>

<template>
  <RouterLink
    :to="`/workflows/${id}`"
    class="max-w-[500px] truncate text-blue-400 hover:underline"
  >
    {{ name }}
  </RouterLink>
</template>
