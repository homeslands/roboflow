<script setup lang="ts">
import type { InputConfig, TriggerNodeDefinition } from '@/types/workflow/node/definition/trigger-node-definition'
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
import { useVueFlow } from '@vue-flow/core'
import { LoaderCircleIcon, PlusIcon, Trash2Icon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { inputConfigsSchema } from './schema'

interface Props {
  nodeId: string
  definition: TriggerNodeDefinition<'ON_DEMAND'>
}

const props = defineProps<Props>()
const definition = toRef(props, 'definition')
const isLoading = ref<boolean>(false)

const { handleSubmit, meta, setFieldValue } = useForm<{ configs: InputConfig[] }>({
  validationSchema: toTypedSchema(inputConfigsSchema),
  initialValues: {
    configs: definition.value.configs,
  },
})
const isFormValid = computed(() => meta.value.valid)

const { updateNodeData } = useVueFlow()

function updateInputConfigs() {
  setFieldValue('configs', definition.value.configs)
}

function addInput(inputType: InputConfig['inputType']) {
  definition.value.configs.push({
    key: '',
    inputType,
    defaultValue: '',
    required: false,
  })
  updateInputConfigs()
}

function removeInput(index: number) {
  definition.value.configs.splice(index, 1)
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

const onSubmit = handleSubmit((values) => {
  isLoading.value = true
  setTimeout(() => {
    updateNodeData(props.nodeId, {
      definition: {
        ...definition.value,
        configs: values.configs,
      },
    })
    isLoading.value = false
  }, 500)
})
</script>

<template>
  <form class="space-y-4" @submit="onSubmit">
    <FormField v-slot="{ errorMessage }" name="configs">
      <FormItem>
        <FormControl>
          <div class="space-y-4">
            <span class="text-sm font-medium">
              Environment variables
            </span>
            <!-- <div class="flex items-center gap-2">
              <p class="w-1/3 text-sm text-muted-foreground">
                Key
              </p>
              <p class="flex-1 text-sm text-muted-foreground">
                Value
              </p>
              <p class="mr-12 text-sm text-muted-foreground">
                *
              </p>
            </div> -->
            <div class="grid grid-cols-12 gap-2 text-sm text-muted-foreground">
              <div class="col-span-4">
                Key
              </div>
              <div class="col-span-6">
                Value
              </div>
              <div class="col-span-2 text-center">
                Required
              </div>
            </div>

            <div
              v-for="(config, index) in definition.configs "
              :key="index"
              class="grid items-center grid-cols-12 gap-2"
            >
              <Input
                :id="`key-${index}`"
                v-model="config.key"
                class="col-span-4"
                placeholder="Add key"
                @input="updateInputConfigs"
              />
              <Input
                :id="`value-${index}`"
                v-model="config.defaultValue"
                class="col-span-6"
                :placeholder="getPlaceHolderValueText(config.inputType)"
                :type="getDefaultValueType(config.inputType)"
                @input="updateInputConfigs"
              />
              <div class="flex items-center justify-center col-span-2 gap-2">
                <Checkbox
                  v-model:checked="config.required"
                  @update:checked="updateInputConfigs"
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  class="w-8 h-8"
                  @click="removeInput(index)"
                >
                  <Trash2Icon class="w-4 h-4" />
                </Button>
              </div>
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

    <div class="flex justify-end">
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
