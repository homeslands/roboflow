import * as z from 'zod'

export const createWorkflowTempParamsSchema = z.object({
  name: z.string()
    .min(1, { message: 'Name is required' })
    .max(100, { message: 'Name is too long' })
    .regex(/^[a-z0-9 ]+$/i, { message: 'Name can only contain letters, numbers, and spaces' }),
})

export type CreateWorkflowTempParams = z.infer<typeof createWorkflowTempParamsSchema>
