<script setup lang="ts">
defineProps<{
  id: string
  name: string
}>()
</script>

<template>
  <RouterLink :to="`/raybots/${id}`" class="text-primary">
    {{ name }}
  </RouterLink>
</template>
