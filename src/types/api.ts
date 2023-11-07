export type ValidationResult = {
  valid: boolean
  errors: string[]
}

export type ValidationBody = {
  path: string[]
  value: any
}
