import 'vue-toastification/dist/index.css'
import './assets/main.scss'

import { WasmAPIStub } from './loadWasm/wasm_api'
import './loadWasm/wasm_exec.d.ts'
import './loadWasm/wasm_exec.js'

import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { createApp } from 'vue'

import App from './App.vue'
import router from './router'

import Toast from 'vue-toastification'
import { useGlobalStore } from './stores/global'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

const app = createApp(App)
app.use(pinia)
app.use(router)
app.use(Toast, {
  maxToasts: 1,
  filterBeforeCreate: (toast: any, toasts: any) => {
    if (toasts.filter((t: any) => t.type === toast.type).length !== 0) {
      return false
    }
    return toast
  }
})

window.WasmAPI = new WasmAPIStub()
const go = new Go()
const sourceFile = 'main.wasm'
WebAssembly.instantiateStreaming(fetch(sourceFile), go.importObject).then((result) => {
  go.run(result.instance)
  useGlobalStore().setWasmInitialized()
})

app.mount('#app')
