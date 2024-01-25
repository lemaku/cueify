/* eslint-disable @typescript-eslint/no-unused-vars */
import '@/types/app.d.ts'
import './wasm_exec.d.ts'

export class WasmAPIStub implements WasmAPI {
  Validate(path: Path, payload: any, schema: string) {
    return JSON.parse(this._validate(path, JSON.stringify(payload), schema))
  }
  ValidateSchema(payload: any) {
    return JSON.parse(this._validateSchema(payload))
  }
  Inspect(path: Path, payload: any, schema: string) {
    return JSON.parse(this._inspect(path, JSON.stringify(payload), schema))
  }
  Summarize(payload: any, schema: string) {
    return JSON.parse(this._summarize(JSON.stringify(payload), schema))
  }
  _validate(path: Path, json: string, schema: string) {
    return '{}'
  }
  _validateSchema(raw: string) {
    return '{}'
  }
  _inspect(path: Path, json: string, schema: string) {
    return '{}'
  }
  _summarize(json: string, schema: string) {
    return '{}'
  }
}
