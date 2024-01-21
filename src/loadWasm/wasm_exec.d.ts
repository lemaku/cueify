declare class Go {
  constructor()
  importObject: any
  run(v: any): void
}

declare interface WasmAPI {
  Validate: (path: Path, payload: any) => { valid: boolean; errors: string[] }
  Inspect: (path: Path, payload: any) => { type: CurrentType; properties: Field[] }
  Summarize: (payload: any) => { value: any; valid: boolean; errors: ValueError[] }
  _validate: (path: Path, json: string) => string
  _inspect: (path: Path, json: string) => string
  _summarize: (json: string) => string
}

declare interface Window {
  WasmAPI: WasmAPI
}
