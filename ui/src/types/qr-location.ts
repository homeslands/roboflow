export interface QRLocation {
  id: string
  name: string
  qrCode: string
  metadata: Record<string, string | number | boolean>
  createdAt: Date
  updatedAt: Date
}
