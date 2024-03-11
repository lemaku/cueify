type FieldType = 'bottom' | 'null' | 'string' | 'bytes' | 'int' | 'float' | 'bool' | 'list' | 'struct';
type Path = string[]

type Field = {
  path: Path
  type: FieldType[]
  optional: boolean
}

type CurrentType = 'list' | 'struct'

type ValueError = {
  path: Path
  errors: string[]
}

type ValidationError = {
  self: string[]
  others: {
    [path: string]: string[]
  }
}

type BreadCrumb = {
  crumb: string
  path: string[]
}