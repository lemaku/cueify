import './assets/main.scss'

import { WasmAPIStub } from './loadWasm/wasm_api'
import './loadWasm/wasm_exec.d.ts'
import './loadWasm/wasm_exec.js'

import { createPinia } from 'pinia'
import { createApp } from 'vue'
import PrimeVue from 'primevue/config'

import App from './App.vue'
import router from './router'

import { useGlobalStore } from './stores/global'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(PrimeVue, { styled: false })

window.WasmAPI = new WasmAPIStub()
const go = new Go()
const sourceFile = 'main.wasm'
WebAssembly.instantiateStreaming(fetch(sourceFile), go.importObject).then((result) => {
  go.run(result.instance)
  useGlobalStore().setWasmInitialized()
})

app.mount('#app')
