import type { Path, Field, CurrentType } from "./app"

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