import * as z from 'zod'

const inputSchemas = {
  STOP: z.undefined(),
  MOVE_FORWARD: z.undefined(),
  MOVE_BACKWARD: z.undefined(),
  MOVE_TO_LOCATION: z.object({
    location: z.string(),
    direction: z.enum(['FORWARD', 'BACKWARD']),
  }),
  OPEN_BOX: z.undefined(),
  CLOSE_BOX: z.undefined(),
  LIFT_BOX: z
    .object({
      distance: z.number().min(30).max(80).optional(),
    })
    .optional(),
  DROP_BOX: z
    .object({
      distance: z.number().min(30).max(80).optional(),
    })
    .optional(),
  CHECK_QR: z.object({
    qr_code: z.string(),
  }),
  WAIT_GET_ITEM: z.undefined(),
  SCAN_QR_LOCATION: z.undefined(),
}

export const createRaybotCommandSchema = z.discriminatedUnion('type', [
  z.object({ type: z.literal('STOP'), input: inputSchemas.STOP }),
  z.object({ type: z.literal('MOVE_FORWARD'), input: inputSchemas.MOVE_FORWARD }),
  z.object({ type: z.literal('MOVE_BACKWARD'), input: inputSchemas.MOVE_BACKWARD }),
  z.object({ type: z.literal('MOVE_TO_LOCATION'), input: inputSchemas.MOVE_TO_LOCATION }),
  z.object({ type: z.literal('OPEN_BOX'), input: inputSchemas.OPEN_BOX }),
  z.object({ type: z.literal('CLOSE_BOX'), input: inputSchemas.CLOSE_BOX }),
  z.object({ type: z.literal('LIFT_BOX'), input: inputSchemas.LIFT_BOX }),
  z.object({ type: z.literal('DROP_BOX'), input: inputSchemas.DROP_BOX }),
  z.object({ type: z.literal('CHECK_QR'), input: inputSchemas.CHECK_QR }),
  z.object({ type: z.literal('WAIT_GET_ITEM'), input: inputSchemas.WAIT_GET_ITEM }),
  z.object({ type: z.literal('SCAN_QR_LOCATION'), input: inputSchemas.SCAN_QR_LOCATION }),
])
