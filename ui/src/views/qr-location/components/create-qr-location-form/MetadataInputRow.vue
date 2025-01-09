<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { XIcon } from 'lucide-vue-next'

interface MetadataField {
  key: string
  value: string | number | boolean
}

const props = defineProps<{
  field: MetadataField
  index: number
}>()

const emit = defineEmits<{
  update: [index: number, field: MetadataField]
  remove: [index: number]
}>()

function updateField(type: 'key' | 'value', value: string) {
  emit('update', props.index, {
    ...props.field,
    [type]: value,
  })
}
</script>

<template>
  <div class="flex items-start gap-2">
    <Input
      placeholder="Key"
      :value="field.key"
      class="w-1/3"
      @update:model-value="(value) => updateField('key', value as string)"
    />
    <Input
      placeholder="Value"
      :value="field.value"
      class="flex-1"
      @update:model-value="(value) => updateField('value', value as string)"
    />
    <Button
      variant="ghost"
      size="icon"
      type="button"
      class="shrink-0"
      @click="() => emit('remove', index)"
    >
      <XIcon class="w-4 h-4" />
    </Button>
  </div>
</template>
