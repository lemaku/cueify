/* eslint-disable @typescript-eslint/no-unused-vars */
import '@/types/app.d.ts'
import './wasm_exec.d.ts'
import emitter from '../event-bus'

export class WasmAPIStub implements WasmAPI {
  Validate(path: Path, payload: any, schema: string) {
    const result = this._validate(path, JSON.stringify(payload), schema)
    if (!result) {
      emitter.$emit('wasm-error')
      return undefined
    }
    return JSON.parse(result)
  }
  ValidateSchema(payload: any) {
    const result = this._validateSchema(payload)
    if (!result) {
      emitter.$emit('wasm-error')
      return undefined
    }
    return JSON.parse(result)
  }
  Inspect(path: Path, payload: any, schema: string) {
    const result = this._inspect(path, JSON.stringify(payload), schema)
    if (!result) {
      emitter.$emit('wasm-error')
      return undefined
    }
    return JSON.parse(result)
  }
  Summarize(payload: any, schema: string) {
    const result = this._summarize(JSON.stringify(payload), schema)
    if (!result) {
      emitter.$emit('wasm-error')
      return undefined
    }
    return JSON.parse(result)
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
