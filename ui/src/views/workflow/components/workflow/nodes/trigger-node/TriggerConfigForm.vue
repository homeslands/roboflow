<script setup lang="ts">
import type { InputConfig } from '@/types/workflow/node/definition/trigger-node-definition'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { toTypedSchema } from '@vee-validate/zod'

import { LoaderCircleIcon, PlusIcon, Trash2Icon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { inputConfigsSchema } from './schema'

const inputConfigs = defineModel<InputConfig[]>('inputConfigs', { required: true })
const isLoading = ref<boolean>(false)

const { handleSubmit, resetForm, meta, setFieldValue } = useForm<{
  configs: InputConfig[]
}>({
  validationSchema: toTypedSchema(inputConfigsSchema),
  initialValues: {
    configs: inputConfigs.value,
  },
})
const onSubmit = handleSubmit((values) => {
  isLoading.value = true

  setTimeout(() => {
    isLoading.value = false
    inputConfigs.value = values.configs
    resetForm()
    inputConfigs.value = []
  }, 2000)
})

const isFormValid = computed(() => meta.value.valid)

function updateInputConfigs() {
  setFieldValue('configs', inputConfigs.value)
}

function addInput(inputType: InputConfig['inputType']) {
  inputConfigs.value.push({
    key: '',
    inputType,
    defaultValue: '',
    required: false,
  })
  updateInputConfigs()
}

function removeInput(index: number) {
  inputConfigs.value.splice(index, 1)
  updateInputConfigs()
}

function getPlaceHolderValueText(type: InputConfig['inputType']) {
  switch (type) {
    case 'text':
      return 'Add text'
    case 'number':
      return 'Add number'
    default:
      return 'Add text'
  }
}

function getDefaultValueType(type: InputConfig['inputType']) {
  switch (type) {
    case 'text':
      return 'string'
    case 'number':
      return 'number'
    default:
      return 'string'
  }
}
</script>

<template>
  <form @submit="onSubmit">
    <FormField v-slot="{ errorMessage }" name="configs">
      <FormItem>
        <FormControl>
          <div class="space-y-2">
            <span class="text-sm font-medium">
              Environment variables
            </span>
            <div class="flex items-center gap-2">
              <p class="w-1/3 text-sm text-muted-foreground">
                Key
              </p>
              <p class="flex-1 text-sm text-muted-foreground">
                Value
              </p>
              <p class="mr-12 text-sm text-muted-foreground">
                *
              </p>
            </div>
            <h3 />
            <div v-for="(config, index) in inputConfigs" :key="index" class="flex items-center gap-2">
              <Input
                v-model="config.key"
                class="w-1/3"
                placeholder="Add key"
                @input="updateInputConfigs"
              />
              <Input
                v-model="config.defaultValue"
                class="flex-1"
                :placeholder="getPlaceHolderValueText(config.inputType)"
                :type="getDefaultValueType(config.inputType)"
                @input="updateInputConfigs"
              />
              <Checkbox
                v-model="config.required"
                @input="updateInputConfigs"
              />
              <Button
                type="button"
                variant="ghost"
                size="icon"
                @click="removeInput(index)"
              >
                <Trash2Icon class="w-4 h-4" />
              </Button>
            </div>

            <DropdownMenu>
              <DropdownMenuTrigger>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                >
                  <PlusIcon class="w-4 h-4" />
                  Add input
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuLabel>Select input type</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem @click="addInput('text')">
                  Text
                </DropdownMenuItem>
                <DropdownMenuItem @click="addInput('number')">
                  Number
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </FormControl>
        <FormMessage>{{ errorMessage }}</FormMessage>
      </FormItem>
    </FormField>
    <div class="flex justify-end pr-3 mt-4">
      <Button
        type="submit"
        :disabled="!isFormValid || isLoading"
      >
        <LoaderCircleIcon v-if="isLoading" class="w-4 h-4 animate-spin" />
        Save
      </Button>
    </div>
  </form>
</template>
