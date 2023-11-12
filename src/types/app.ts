export type FieldType = 'string' | 'number' | 'bool' | 'list' | 'complex';
export type Path = string[]

export type Field = {
    path: Path
    type: FieldType
    index: number
}

export type CurrentType = 'list' | 'complex';