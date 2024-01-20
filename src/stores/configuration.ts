import { defineStore } from 'pinia'
import { inspect, summarize, validate } from '@/services/rest'
import { cloneDeep } from 'lodash'
import type { Field, CurrentType, ValueError, Path, BreadCrumb } from '@/types/app'
import router from '@/router'

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
              name: 'Leon K',
              semester: 5,
              active: true
            }
          ]
        }
      }
    } as any,
    rawCurrentType: 'complex' as CurrentType,
    rawFields: [] as Field[],
    rawFormat: 'json' as Format,
    rawErrors: [] as ValueError[]
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
        const result = await inspect(path, this.rawCurrent)
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
      const newCurrent = setValue(path, value, this.rawCurrent)
      const res = await validate(path, newCurrent)

      if (res.valid) {
        this.rawCurrent = newCurrent
        await this.summarize()
      }

      return res
    },
    async unset(path: string[]) {
      if (this.get(path) !== undefined) {
        this.rawCurrent = unsetValue(path, this.rawCurrent)
        this.summarize()
        this.rawFields = (await inspect(this.rawPath, this.rawCurrent)).properties
      }
    },
    async setToEmpty(path: string[], isArray = false) {
      this.rawCurrent = setValue(path, isArray ? [] : {}, this.rawCurrent)
      await this.summarize()
    },
    async addToArray() {
      this.rawCurrent = pushToArray(this.rawPath, this.rawCurrent)
      this.summarize()
      this.rawFields = (await inspect(this.rawPath, this.rawCurrent)).properties
    },
    async summarize() {
      const result = await summarize(this.rawCurrent)
      this.rawErrors = result.errors
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
