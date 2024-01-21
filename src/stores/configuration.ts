import router from '@/router'
import '@/types/app.d.ts'
import type { MutexInterface } from 'async-mutex'
import { Mutex } from 'async-mutex'
import { cloneDeep } from 'lodash'
import { defineStore } from 'pinia'

export const supportedFormats = ['json', 'yaml'] as const
export type Format = (typeof supportedFormats)[number]

export const useConfigurationStore = defineStore({
  id: 'configuration',
  state: () => ({
    rawPath: [] as string[],
    rawCurrent: {
      universities: {
        tuwien: {
          name: 'Vienna University of Technology',
          students: [
            {
              matNr: '12119877',
              name: 'Leon K'
            }
          ]
        }
      }
    } as any,
    rawCurrentType: 'complex' as CurrentType,
    rawFields: [] as Field[],
    rawFormat: 'json' as Format,
    rawErrors: [] as ValueError[],
    lock: new Mutex() as MutexInterface
  }),
  getters: {
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
      const resultArray = []
      const path = []

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
    async jumpTo(path: string[]) {
      if (path && JSON.stringify(path) != JSON.stringify(this.path)) {
        const result = window.WasmAPI.Inspect(path, this.rawCurrent)

        if (result.type != 'complex' && result.type != 'list') {
          const parent = path.slice(0, path.length - 1)
          const next = parent.length <= 0 ? ['universities'] : parent
          router.push({ query: { p: next.join('.') } })
        } else {
          this.rawPath = path
          this.rawCurrentType = result.type
          this.rawFields = result.properties
        }
      }
    },
    async set(path: string[], value: any) {
      await this.lock.acquire()
      const newRaw = setValue(path, value, this.rawCurrent)
      const res = window.WasmAPI.Validate(path, newRaw)

      if (res.valid) {
        await this.summarize(newRaw)
      }

      this.lock.release()
      return res
    },
    async unset(path: string[]) {
      await this.lock.acquire()
      if (this.get(path) !== undefined) {
        const newRaw = unsetValue(path, this.rawCurrent)
        await this.summarize(newRaw)
        this.rawFields = window.WasmAPI.Inspect(this.rawPath, this.rawCurrent).properties
      }
      this.lock.release()
    },
    async setToEmpty(path: string[], isArray = false) {
      await this.lock.acquire()
      const newRaw = setValue(path, isArray ? [] : {}, this.rawCurrent)
      await this.summarize(newRaw)
      this.lock.release()
    },
    async addToArray() {
      await this.lock.acquire()
      const newRaw = pushToArray(this.rawPath, this.rawCurrent)
      this.summarize(newRaw)
      this.rawFields = window.WasmAPI.Inspect(this.rawPath, this.rawCurrent).properties
      this.lock.release()
    },
    async summarize(raw: any) {
      const result = window.WasmAPI.Summarize(raw)
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

  console.log(Array.isArray(obj))

  if (Array.isArray(obj)) {
    obj.push({})
  }

  return ref
}
