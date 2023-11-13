export type FieldType = 'string' | 'number' | 'bool' | 'list' | 'complex'
export type Path = string[]
export type Error = string

export type Field = {
  path: Path
  type: FieldType
  index: number
}

export type CurrentType = 'list' | 'complex'

export type ValueError = {
  path: Path
  errors: Error[]
}

export type BreadCrumb = {
  crumb: string
  path: string[]
}