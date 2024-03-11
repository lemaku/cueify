<template>
  <header>
    <nav>
      <template v-if="schemaSet">
        <button class="icon-button back-button" @click="onBack()">
          <i class="pi pi-arrow-left"></i>
        </button>
      </template>
      <RouterLink to="/">Home</RouterLink>
      <RouterLink to="/about">About</RouterLink>
    </nav>
  </header>

  <loading
    v-if="!wasmInitialized"
    :active="true"
    :is-full-page="true"
    :loader="'dots'"
    :color="'#4a86e8'"
  ></loading>
  <div v-else class="main">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { useGlobalStore } from '@/stores/global'
import { useConfigurationStore } from '@/stores/configuration'
import { storeToRefs } from 'pinia'
import Loading from 'vue-loading-overlay'
import 'vue-loading-overlay/dist/css/index.css'
import { RouterLink, RouterView } from 'vue-router'
import { useToast } from 'vue-toastification'
import emitter from './event-bus'

const global = useGlobalStore()
const { wasmInitialized } = storeToRefs(global)

const configuration = useConfigurationStore()
const { schemaSet } = storeToRefs(configuration)
const onBack = () => {
  if (confirm('Do you really want to restart? You will loose your progress...')) {
    configuration.$reset()
  }
}

const toast = useToast()
emitter.$on('wasm-error', () => {
  toast.error(
    'Ups...something went wrong. This is either a bug or you used an unsupported feature. Sorry!',
    {
      closeOnClick: false,
      closeButton: false,
      pauseOnHover: false,
      draggable: false,
      showCloseButtonOnHover: false,
      hideProgressBar: true,
      timeout: false
    }
  )
})
</script>

<style scoped>
header {
  font-size: 1.5rem;
  padding: 0.5rem;
  text-align: center;
}

.main {
  width: 100%;
  max-width: 70vw;
  min-width: 800px;
  padding: 20px;
}

nav {
  text-align: center;
  display: flex;
  align-items: center;
}

nav a {
  padding: 0 1rem;
  border-left: 1px solid var(--color-text);
}

nav a:first-of-type {
  border: 0;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

.back-button {
  padding: 0;
  margin-left: -1.1rem;
}
</style>
