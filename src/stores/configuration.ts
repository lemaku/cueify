import { defineStore } from 'pinia'
import { validate } from '@/services/rest'
import { cloneDeep } from 'lodash';

export const supportedFormats = ['json', 'yaml'] as const
export type Format = (typeof supportedFormats)[number]

export type BreadCrumb = {
  crumb: string
  path: string[]
}

export type Input = {
  type: 'string' | 'number' | 'bool'
  path: string[]
}

export const useConfigurationStore = defineStore({
  id: 'configuration',
  state: () => ({
    rawPath: ['universities', 'tuwien', 'students', '0'],
    rawCurrent: {
      universities: {
        tuwien: {
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
    rawFields: [] as Input[],
    rawFormat: 'json' as Format
  }),
  getters: {
    fields: (state): Input[] => {
      return state.rawFields
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
    jumpTo(path: string[]) {
      if (path && JSON.stringify(path) != JSON.stringify(this.path)) {
        console.log('jumping to ', path)
        this.rawPath = path
      }
    },
    async load() {
      console.log('syncing with server')
      // send to api
      this.rawFields = [
        {
          type: 'string',
          path: ['universities', 'tuwien', 'students', '0', 'matNr']
        },
        {
          type: 'string',
          path: ['universities', 'tuwien', 'students', '0', 'name']
        },
        {
          type: 'number',
          path: ['universities', 'tuwien', 'students', '0', 'semester']
        },
        {
          type: 'bool',
          path: ['universities', 'tuwien', 'students', '0', 'active']
        }
      ]
    },
    async set(path: string[], value: any) {
      const newCurrent = setValue(path, value, this.rawCurrent)
      const res = await validate(path, newCurrent);

      if (res.valid) {
        this.rawCurrent = newCurrent;
      }

      return res;
    },
    changeFormat(format: Format) {
      this.rawFormat = format
    }
  }
})

function setValue(path: string[], value: any, object: any) : any {
  const ref = cloneDeep(object ?? {});
  let obj = ref;
  let i = 0
  for (i = 0; i < path.length - 1; i++) {
    obj = obj[path[i]]
  }

  obj[path[i]] = value;
  return ref;
}