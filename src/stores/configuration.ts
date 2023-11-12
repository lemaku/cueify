import { defineStore } from 'pinia'
import { inspect, validate } from '@/services/rest'
import { cloneDeep } from 'lodash'
import type { Field, CurrentType } from '@/types/app'

export const supportedFormats = ['json', 'yaml'] as const
export type Format = (typeof supportedFormats)[number]

export type BreadCrumb = {
  crumb: string
  path: string[]
}

export const useConfigurationStore = defineStore({
  id: 'configuration',
  state: () => ({
    rawPath: ['universities', 'tuwien', 'students', '0'],
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
    rawFormat: 'json' as Format
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
    path: (state): any => {
      return state.rawPath
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
      return (path: string) => {
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
        this.rawPath = path
        this.load()
      }
    },
    async load() {
      const result = await inspect(this.rawPath, this.rawCurrent)
      this.rawCurrentType = result.type
      this.rawFields = result.properties
    },
    async set(path: string[], value: any) {
      const newCurrent = setValue(path, value, this.rawCurrent)
      const res = await validate(path, newCurrent)

      if (res.valid) {
        this.rawCurrent = newCurrent
      }

      return res
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
