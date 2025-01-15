import * as z from 'zod'

const inputConfigSchema = z.object({
  key: z.string().min(1, 'Key must be at least 1 character long'),
  inputType: z.enum(['text', 'number'], {
    message: 'Input type must be either "text", "number"',
  }),
  defaultValue: z.string().min(1, 'Default value must be at least 1 character long'),
  required: z.boolean(),
})

export const inputConfigsSchema = z.object({
  configs: z.array(inputConfigSchema),
})
