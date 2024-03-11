import '@/types/app.d.ts'
import type { MutexInterface } from 'async-mutex'
import { Mutex } from 'async-mutex'
import { cloneDeep, isArray, isEmpty, isEqual, isObject } from 'lodash'
import { defineStore } from 'pinia'
import emitter from '@/event-bus'
import { useGlobalStore } from './global'
import { watch } from 'vue'

export const supportedFormats = ['json', 'yaml'] as const
export type Format = (typeof supportedFormats)[number]

export const useConfigurationStore = defineStore({
  id: 'configuration',
  state: () => ({
    rawSchema: '' as string,
    rawPath: [] as string[],
    rawCurrent: {} as any,
    rawCurrentType: 'struct' as CurrentType,
    rawCurrentOf: undefined as undefined | FieldType[],
    rawFields: [] as Field[],
    rawFormat: 'json' as Format,
    rawErrors: [] as ValueError[],
    lock: new Mutex() as MutexInterface
  }),
  persist: {
    key: 'cueify.appdata',
    paths: ['rawSchema', 'rawPath', 'rawCurrent', 'rawFormat'],
    storage: sessionStorage,
    afterRestore: (ctx) => {
      const { rawSchema, rawPath, rawCurrent, rawFormat } = ctx.store.$state
      const stop = watch(useGlobalStore(), (state) => {
        if (state.wasmInitialized) {
          if (rawSchema && isArray(rawPath) && isObject(rawCurrent) && rawFormat) {
            try {
              ctx.store.summarize(rawCurrent)
              ctx.store.jumpTo(rawPath)
            } catch (_) {
              ctx.store.$reset()
            }
          } else {
            ctx.store.$reset()
          }
          stop()
        }
      })
    }
  },
  getters: {
    schemaSet: (state): boolean => {
      return state.rawSchema ? true : false
    },
    fields: (state): Field[] => {
      return state.rawFields
    },
    currentType: (state): CurrentType => {
      return state.rawCurrentType
    },
    currentOf: (state): FieldType[] | undefined => {
      return state.rawCurrentOf
    },
    current: (state): any => {
      return state.rawCurrent
    },
    path: (state): Path => {
      return state.rawPath
    },
    errors: (state): ValueError[] => {
      return state.rawErrors
    },
    isValid: (state): boolean => {
      return isEmpty(state.rawErrors)
    },
    breadcrumbs: (state): BreadCrumb[] => {
      const path: string[] = []
      const resultArray: BreadCrumb[] = [{ crumb: '~', path: [] }]

      for (const crumb of state.rawPath) {
        path.push(crumb)
        resultArray.push({ crumb, path: [...path] })
      }

      return resultArray
    },
    format: (state): Format => {
      return state.rawFormat
    },
    get: (state) => {
      return (path: string[]) => {
        return getValue(path, state.rawCurrent)
      }
    },
    isDerived: (state) => {
      // Calculate diff between current value and the value that would result off unsetting the current value
      // => If there is no difference, it means that the current value is derived and unsetting has no effect
      return (path: string[]) => {
        const currentValue = getValue(path, state.rawCurrent)
        const newRaw = unsetValue(path, state.rawCurrent)
        const result = window.WasmAPI.Summarize(newRaw, state.rawSchema)
        const newValue = getValue(path, result.value)
        return isEqual(currentValue, newValue)
      }
    }
  },
  actions: {
    setSchema(schema: string) {
      const res = window.WasmAPI.ValidateSchema(schema)

      if (res.valid) {
        this.rawSchema = schema

        this.rawCurrent = {}

        this.summarize(this.rawCurrent)
        this.jumpTo([])
      }
      return res
    },
    jumpTo(path: string[]) {
      const result = window.WasmAPI.Inspect(path, this.rawCurrent, this.rawSchema)

      if (!isEqual(result.type, ['struct']) && !isEqual(result.type, ['list'])) {
        this.jumpTo(path.slice(0, path.length - 1))
        // Give components time to be rendered and then trigger focus event
        setTimeout(() => emitter.$emit('focus', path), 50)
      } else {
        this.rawPath = path
        this.rawCurrentType = result.type[0] as CurrentType
        this.rawCurrentOf = result.of as CurrentType[] | undefined
        this.rawFields = result.properties
      }
    },
    async set(path: string[], value: any) {
      await this.lock.acquire()
      const newRaw = setValue(path, value, this.rawCurrent)
      const res = window.WasmAPI.Validate(path, newRaw, this.rawSchema)

      if (res.valid) {
        this.summarize(newRaw)
        this.jumpTo(this.rawPath)
      }

      this.lock.release()
      return res
    },
    async unset(path: string[]) {
      await this.lock.acquire()
      if (this.get(path) !== undefined) {
        const newRaw = unsetValue(path, this.rawCurrent)
        this.summarize(newRaw)
        this.jumpTo(this.rawPath)
      }
      this.lock.release()
    },
    async setToEmpty(path: string[], isArray = false) {
      await this.lock.acquire()
      const newRaw = setValue(path, isArray ? [] : {}, this.rawCurrent)
      this.summarize(newRaw)
      this.jumpTo(this.rawPath)
      this.lock.release()
    },
    async addToArray(toPush: any) {
      await this.lock.acquire()
      const newRaw = pushToArray(this.rawPath, this.rawCurrent, toPush)
      this.summarize(newRaw)
      this.jumpTo(this.rawPath)
      this.lock.release()
    },
    summarize(raw: any) {
      const result = window.WasmAPI.Summarize(raw, this.rawSchema)
      this.rawErrors = result.errors
      this.rawCurrent = result.value
    },
    changeFormat(format: Format) {
      this.rawFormat = format
    }
  }
})

function getValue(path: string[], value: any) {
  let obj = value
  let i = 0
  for (i = 0; i < path.length - 1; i++) {
    obj = obj[path[i]]
  }
  return obj[path[i]]
}

function setValue(path: string[], value: any, object: any): any {
  const ref = cloneDeep(object ?? {})
  let obj = ref
  let i = 0
  for (i = 0; i < path.length - 1; i++) {
    obj = obj[path[i]]
  }

  obj[path[i]] = value
  return ref
}

function unsetValue(path: string[], object: any): any {
  const ref = cloneDeep(object ?? {})
  let obj = ref
  let i = 0

  for (i = 0; i < path.length - 1; i++) {
    obj = obj[path[i]]
  }

  if (Array.isArray(obj)) {
    obj.splice(path[i] as unknown as number, 1)
  } else {
    delete obj[path[i]]
  }

  return ref
}

function pushToArray(path: string[], object: any, toPush: any): any {
  const ref = cloneDeep(object ?? {})
  let obj = ref
  let i = 0

  for (i = 0; i < path.length; i++) {
    obj = obj[path[i]]
  }

  if (Array.isArray(obj)) {
    obj.push(toPush)
  }

  return ref
}
