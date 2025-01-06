import * as z from 'zod'

export const createQRLocationParamsSchema = z.object({
  name: z.string().min(1, { message: 'Name is required' }),
  qrCode: z.string().min(1, { message: 'QR Code is required' }),
  metadata: z.record(z.union([z.string(), z.number(), z.boolean()])),
})

export const metadataSchema = z.record(z.union([z.string(), z.number(), z.boolean()]))
