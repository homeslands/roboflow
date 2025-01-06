export interface ErrorDetail {
  Field: string
  Message: string
}

export class RoboflowError extends Error {
  status: number
  errorCode: string
  details?: ErrorDetail[]

  constructor(message: string, status: number, errorCode: string, details?: ErrorDetail[]) {
    super(message)
    this.status = status
    this.errorCode = errorCode
    this.details = details
  }
}

export class HTTPError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}
