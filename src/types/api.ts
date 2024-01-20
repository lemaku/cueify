import type { CurrentType, Field, Path, ValueError } from './app'

export type ValidationResult = {
  valid: boolean
  errors: Path
}

export type ValidationBody = {
  path: Path
  value: any
}

export type InspectBody = {
  path: Path
  value: any
}

export type InspectResult = {
  type: CurrentType
  properties: Field[]
}

export type SummarizeBody = {
  value: any
}

export type SummarizeResult = {
  value: any
  valid: boolean
  errors: ValueError[]
}
