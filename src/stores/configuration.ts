import '@/types/app.d.ts'
import type { MutexInterface } from 'async-mutex'
import { Mutex } from 'async-mutex'
import { cloneDeep, isEqual } from 'lodash'
import { defineStore } from 'pinia'

export const supportedFormats = ['json', 'yaml'] as const
export type Format = (typeof supportedFormats)[number]

export const useConfigurationStore = defineStore({
  id: 'configuration',
  state: () => ({
    rawSchema: '' as string,
    rawPath: [] as string[],
    rawCurrent: {} as any,
    rawCurrentType: 'struct' as CurrentType,
    rawFields: [] as Field[],
    rawFormat: 'json' as Format,
    rawErrors: [] as ValueError[],
    lock: new Mutex() as MutexInterface
  }),
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
    current: (state): any => {
      return state.rawCurrent
    },
    path: (state): Path => {
      return state.rawPath
    },
    errors: (state): ValueError[] => {
      return state.rawErrors
    },
    isLoading: (state): boolean => {
      return state.lock.isLocked()
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
        let obj = state.rawCurrent
        let i = 0
        for (i = 0; i < path.length - 1; i++) {
          obj = obj[path[i]]
        }
        return obj[path[i]]
      }
    }
  },
  actions: {
    setSchema(schema: string) {
      const res = window.WasmAPI.ValidateSchema(schema)

      if (res.valid) {
        this.rawSchema = schema

        // TODO reevalute this:
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
        setTimeout(() => this.focus(path), 25);
      } else {
        this.rawPath = path
        this.rawCurrentType = result.type[0] as CurrentType;
        this.rawFields = result.properties
      }
    },
    // Fix this as it should just be an event triggered
    focus(path: string[]) {
    },
    async set(path: string[], value: any) {
      await this.lock.acquire()
      const newRaw = setValue(path, value, this.rawCurrent)
      const res = window.WasmAPI.Validate(path, newRaw, this.rawSchema)

      if (res.valid) {
        this.summarize(newRaw)
      }

      this.lock.release()
      return res
    },
    async unset(path: string[]) {
      await this.lock.acquire()
      if (this.get(path) !== undefined) {
        const newRaw = unsetValue(path, this.rawCurrent)
        this.summarize(newRaw)
        this.rawFields = window.WasmAPI.Inspect(
          this.rawPath,
          this.rawCurrent,
          this.rawSchema
        ).properties
      }
      this.lock.release()
    },
    async setToEmpty(path: string[], isArray = false) {
      await this.lock.acquire()
      const newRaw = setValue(path, isArray ? [] : {}, this.rawCurrent)
      this.summarize(newRaw)
      this.lock.release()
    },
    async addToArray() {
      await this.lock.acquire()
      const newRaw = pushToArray(this.rawPath, this.rawCurrent)
      this.summarize(newRaw)
      this.rawFields = window.WasmAPI.Inspect(
        this.rawPath,
        this.rawCurrent,
        this.rawSchema
      ).properties
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

function pushToArray(path: string[], object: any): any {
  const ref = cloneDeep(object ?? {})
  let obj = ref
  let i = 0

  for (i = 0; i < path.length; i++) {
    obj = obj[path[i]]
  }

  if (Array.isArray(obj)) {
    obj.push({}) // TODO: don't just push an object because array might be of different type
  }

  return ref
}
