<script setup lang="ts">
import type { CreateRaybotCommandParams } from '@/api/raybot-command'
import type { RaybotCommandType } from '@/types/raybot-command'
import { Button } from '@/components/ui/button'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useCreateRaybotCommandMutation } from '@/composables/use-raybot-command'
import { RaybotCommandTypeValues } from '@/types/raybot-command'
import { toTypedSchema } from '@vee-validate/zod'
import { LoaderCircleIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import DirectionSelect from './DirectionSelect.vue'
import LocationSelect from './LocationSelect.vue'
import { createRaybotCommandSchema } from './schema'

const route = useRoute()
const raybotId = computed(() => route.params.id as string)

const selectedCommandType = defineModel<RaybotCommandType>()
const { handleSubmit, resetForm, meta } = useForm<CreateRaybotCommandParams<RaybotCommandType>>({
  validationSchema: toTypedSchema(createRaybotCommandSchema),
  initialValues: {
    type: undefined,
    input: undefined,
  },
})

const { mutate, isPending } = useCreateRaybotCommandMutation(raybotId)
const isFormValid = computed(() => meta.value.valid)

const onSubmit = handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      toast.success('Command created successfully')
      resetForm()
    },
    onError: (error) => {
      toast.error(error.message)
    },
  })
})

function getDefaultInputForType(type: RaybotCommandType) {
  switch (type) {
    case 'MOVE_TO_LOCATION':
      return { location: '', direction: 'FORWARD' }
    case 'LIFT_BOX':
    case 'DROP_BOX':
      return { distance: 50 }
    case 'CHECK_QR':
      return { qr_code: '' }
    default:
      return undefined
  }
}

function onCommandTypeChange(type: RaybotCommandType) {
  selectedCommandType.value = type
  resetForm({
    values: {
      type,
      input: getDefaultInputForType(type),
    },
  })
}
</script>

<template>
  <form @submit="onSubmit">
    <div class="grid gap-4">
      <div class="flex flex-col space-y-1.5">
        <!-- Type field -->
        <FormField v-slot="{ field, errorMessage }" name="type">
          <FormItem>
            <FormLabel for="type">
              Command Type
            </FormLabel>
            <Select
              v-model="field.value"
              @update:model-value="(type) => {
                field.onChange(type)
                onCommandTypeChange(type as RaybotCommandType)
              }"
            >
              <FormControl>
                <SelectTrigger id="type">
                  <SelectValue placeholder="Select a command" />
                </SelectTrigger>
              </FormControl>
              <SelectContent position="popper">
                <SelectItem
                  v-for="type in RaybotCommandTypeValues"
                  :key="type"
                  :value="type"
                >
                  {{ type }}
                </SelectItem>
              </SelectContent>
            </Select>
            <FormMessage>{{ errorMessage }}</FormMessage>
          </FormItem>
        </FormField>
      </div>

      <!-- Input field -->
      <!-- MOVE_TO_LOCATION input -->
      <template v-if="selectedCommandType === 'MOVE_TO_LOCATION'">
        <FormField v-slot="{ componentField, errorMessage }" name="input.location">
          <FormItem>
            <FormLabel>Location</FormLabel>
            <LocationSelect
              v-bind="componentField"
            />
            <FormMessage>{{ errorMessage }}</FormMessage>
          </FormItem>
        </FormField>
        <FormField v-slot="{ componentField, errorMessage }" name="input.direction">
          <FormItem>
            <FormLabel>Direction</FormLabel>
            <DirectionSelect
              v-bind="componentField"
            />
            <FormMessage>{{ errorMessage }}</FormMessage>
          </FormItem>
        </FormField>
      </template>

      <!-- LIFT_BOX and DROP_BOX input -->
      <template v-else-if="selectedCommandType === 'LIFT_BOX' || selectedCommandType === 'DROP_BOX'">
        <FormField v-slot="{ componentField, errorMessage }" name="input.distance">
          <FormItem>
            <FormLabel>Distance (30 - 80cm)</FormLabel>
            <Input
              type="number"
              min="30"
              max="80"
              step="5"
              placeholder="50"
              v-bind="componentField"
            />
            <FormMessage>{{ errorMessage }}</FormMessage>
          </FormItem>
        </FormField>
      </template>

      <!-- CHECK_QR input -->
      <template v-else-if="selectedCommandType === 'CHECK_QR'">
        <FormField v-slot="{ componentField, errorMessage }" name="input.qr_code">
          <FormItem>
            <FormLabel>QR Code</FormLabel>
            <Input
              placeholder="QR code"
              v-bind="componentField"
            />
            <FormMessage>{{ errorMessage }}</FormMessage>
          </FormItem>
        </FormField>
      </template>
    </div>

    <div class="flex justify-end mt-4">
      <Button
        type="submit"
        :disable="!isFormValid || isPending"
      >
        <LoaderCircleIcon v-if="isPending" class="w-4 h-4 animate-spin" />
        Create
      </Button>
    </div>
  </form>
</template>
