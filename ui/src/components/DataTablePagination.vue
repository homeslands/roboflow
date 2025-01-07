<script setup lang="ts">
import {
  Pagination,
  PaginationFirst,
  PaginationLast,
  PaginationList,
  PaginationNext,
  PaginationPrev,
} from '@/components/ui/pagination'

interface Props {
  currentPage: number
  pageSize: number
  totalItems: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'pageChange', page: number): void
}>()

const totalPages = computed(() => Math.ceil(props.totalItems / props.pageSize))

function handlePageChange(page: number) {
  emit('pageChange', page)
}
</script>

<template>
  <div class="flex items-center space-x-4">
    <span class="text-sm font-medium text-muted-foreground">
      Page {{ currentPage }} of {{ totalPages }}
    </span>

    <Pagination>
      <PaginationList class="flex items-center gap-1">
        <PaginationFirst
          :disabled="currentPage === 1"
          @click="handlePageChange(1)"
        />
        <PaginationPrev
          :disabled="currentPage === 1"
          @click="handlePageChange(currentPage - 1)"
        />
        <PaginationNext
          :disabled="currentPage === totalPages"
          @click="handlePageChange(currentPage + 1)"
        />
        <PaginationLast
          :disabled="currentPage === totalPages"
          @click="handlePageChange(totalPages)"
        />
      </PaginationList>
    </Pagination>
  </div>
</template>
