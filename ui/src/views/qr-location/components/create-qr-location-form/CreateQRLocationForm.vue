<script setup lang="ts">
import type { CreateQRLocationParams } from '@/api/qr-location'
import { Button } from '@/components/ui/button'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useQRLocationMutation } from '@/composables/use-qr-location'
import { toTypedSchema } from '@vee-validate/zod'
import { LoaderCircleIcon, PlusIcon, XIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { createQRLocationParamsSchema } from './schema'

interface MetadataPair {
  key: string
  value: string
}

const metadataPairs = ref<MetadataPair[]>([{ key: '', value: '' }])
const { handleSubmit, resetForm, setFieldValue } = useForm<CreateQRLocationParams>({
  validationSchema: toTypedSchema(createQRLocationParamsSchema),
  initialValues: {
    name: '',
    qrCode: '',
    metadata: {},
  },
})

const { mutate, isPending } = useQRLocationMutation()
const onSubmit = handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      toast.success('QR Location created successfully')
      resetForm()
      metadataPairs.value = [{ key: '', value: '' }]
    },
    onError: (error) => {
      toast.error(error.message)
    },
  })
})

function updateMetadata() {
  const metadata: Record<string, string | number | boolean> = {}
  metadataPairs.value.forEach((pair) => {
    if (!pair.key)
      return
    const value = pair.value.trim()
    if (!value)
      return

    // Cast to number or boolean if possible
    if (!Number.isNaN(Number(value)))
      metadata[pair.key] = Number(value)
    else if (value === 'true' || value === 'false')
      metadata[pair.key] = value === 'true'
    else
      metadata[pair.key] = value
  })

  setFieldValue('metadata', metadata)
}

function addPair() {
  metadataPairs.value.push({ key: '', value: '' })
}

function removePair(index: number) {
  metadataPairs.value.splice(index, 1)
  updateMetadata()
}
</script>

<template>
  <form @submit="onSubmit">
    <ScrollArea class="h-[calc(100vh-200px)]">
      <div class="grid gap-4 pr-3">
        <div class="flex flex-col space-y-1.5">
          <!-- Name field -->
          <FormField v-slot="{ componentField, errorMessage }" name="name">
            <FormItem>
              <FormLabel>Name</FormLabel>
              <Input
                type="text"
                placeholder="Enter name"
                v-bind="componentField"
              />
              <FormMessage>{{ errorMessage }}</FormMessage>
            </FormItem>
          </FormField>

          <!-- QR code field -->
          <FormField v-slot="{ componentField, errorMessage }" name="qrCode">
            <FormItem>
              <FormLabel>QR code</FormLabel>
              <Input
                type="text"
                placeholder="Enter QR code"
                v-bind="componentField"
              />
              <FormMessage>{{ errorMessage }}</FormMessage>
            </FormItem>
          </FormField>

          <!-- Metadata field -->
          <FormField v-slot="{ errorMessage }" name="metadata">
            <FormItem>
              <FormLabel>Metadata</FormLabel>
              <FormControl>
                <div class="space-y-2">
                  <div class="flex items-center gap-2">
                    <p class="w-1/3 text-sm text-muted-foreground">
                      Key
                    </p>
                    <p class="flex-1 text-sm text-muted-foreground">
                      Value
                    </p>
                  </div>

                  <!-- Loop through metadata fields -->
                  <div v-for="(pair, index) in metadataPairs" :key="index" class="flex items-start gap-2">
                    <Input
                      v-model="pair.key"
                      placeholder="e.g. key"
                      class="w-1/3"
                      @input="updateMetadata"
                    />
                    <Input
                      v-model="pair.value"
                      class="flex-1"
                      @input="updateMetadata"
                    />
                    <Button
                      variant="outline"
                      size="icon"
                      type="button"
                      @click="removePair(index)"
                    >
                      <XIcon class="w-4 h-4" />
                    </Button>
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    class="mt-2"
                    @click="addPair"
                  >
                    <PlusIcon class="w-4 h-4" />
                    Add Field
                  </Button>
                </div>
              </FormControl>
              <FormMessage>{{ errorMessage }}</FormMessage>
            </FormItem>
          </FormField>
        </div>
      </div>
    </ScrollArea>

    <div class="flex justify-end pr-3 mt-4">
      <Button
        type="submit"
        :disable="isPending"
      >
        <LoaderCircleIcon v-if="isPending" class="w-4 h-4 mr-2 animate-spin" />
        Create
      </Button>
    </div>
  </form>
</template>
