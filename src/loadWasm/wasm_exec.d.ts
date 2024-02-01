declare class Go {
  constructor()
  importObject: any
  run(v: any): void
}

declare interface WasmAPI {
  Validate: (path: Path, payload: any, schema: string) => { valid: boolean; errors: string[] }
  ValidateSchema: (payload: any) => { valid: boolean; error: string }
  Inspect: (path: Path, payload: any, schema: string) => { type: FieldType[]; properties: Field[] }
  Summarize: (payload: any, schema: string) => { value: any; valid: boolean; errors: ValueError[] }
  _validate: (path: Path, json: string, schema: string) => string
  _validateSchema: (raw: string) => string
  _inspect: (path: Path, json: string, schema: string) => string
  _summarize: (json: string, schema: string) => string
}

declare interface Window {
  WasmAPI: WasmAPI
}
