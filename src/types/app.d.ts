type FieldType = 'string' | 'number' | 'bool' | 'list' | 'complex'
type Path = string[]

type Field = {
  path: Path
  type: FieldType
  index: number
}

type CurrentType = 'list' | 'complex'

type ValueError = {
  path: Path
  errors: string[]
}

type BreadCrumb = {
  crumb: string
  path: string[]
}