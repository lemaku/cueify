/* eslint-disable @typescript-eslint/no-unused-vars */
import '@/types/app.d.ts'
import './wasm_exec.d.ts'

export class WasmAPIStub implements WasmAPI {
  Validate(path: Path, payload: any) {
    return JSON.parse(this._validate(path, JSON.stringify(payload)))
  }
  Inspect(path: Path, payload: any) {
    return JSON.parse(this._inspect(path, JSON.stringify(payload)))
  }
  Summarize(payload: any) {
    return JSON.parse(this._summarize(JSON.stringify(payload)))
  }
  _validate(path: Path, json: string) {
    return '{}'
  }
  _inspect(path: Path, json: string) {
    return '{}'
  }
  _summarize(json: string) {
    return '{}'
  }
}
