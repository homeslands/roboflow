<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { toTypedSchema } from '@vee-validate/zod'
import { LoaderCircleIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { type CreateWorkflowTempParams, createWorkflowTempParamsSchema } from './schema'

const router = useRouter()

const { handleSubmit, meta } = useForm<CreateWorkflowTempParams>({
  validationSchema: toTypedSchema(createWorkflowTempParamsSchema),
  initialValues: {
    name: '',
  },
})

const isLoading = ref<boolean>(false)
const isFormValid = computed(() => meta.value.valid)

const onSubmit = handleSubmit((values) => {
  isLoading.value = true
  setTimeout(() => {
    router.push({ name: 'NewWorkflow', query: { name: values.name } })
    isLoading.value = false
  }, 500)
})
</script>

<template>
  <form @submit="onSubmit">
    <div class="grid gap-4 pr-3">
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
    </div>

    <div class="flex justify-end pr-3 mt-4">
      <Button
        type="submit"
        :disabled="!isFormValid || isLoading"
      >
        <LoaderCircleIcon v-if="isLoading" class="w-4 h-4 animate-spin" />
        Create
      </Button>
    </div>
  </form>
</template>
