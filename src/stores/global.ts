import { defineStore } from 'pinia'
import '../loadWasm/wasm_exec.d.ts'

export const useGlobalStore = defineStore({
  id: 'global',
  state: () => ({
    rawWasmInitialized: false as boolean
  }),
  getters: {
    wasmInitialized: (state): boolean => {
      return state.rawWasmInitialized
    }
  },
  actions: {
    setWasmInitialized() {
      this.rawWasmInitialized = true
    }
  }
})
