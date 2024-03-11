<template>
  <div class="editor-container">
    <codemirror
      v-model="code"
      class="editor"
      placeholder="CUE definition goes here..."
      :autofocus="false"
      :indent-with-tab="true"
      :tab-size="2"
      :extensions="extensions"
      :disabled="true"
      :style="{ width: '100%' }"
    />
    <button v-if="isValid" class="icon-button copy-button" @click="onCopy()">
      <i class="pi pi-copy"></i>
    </button>
  </div>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { stringify } from 'yaml'
import { Codemirror } from 'vue-codemirror'
import { readOnlyTheme } from './editor-theme'
import copy from 'copy-to-clipboard'
import { useToast } from 'vue-toastification'

const props = defineProps(['path', 'depth'])

const configuration = useConfigurationStore()
const { get, current, isValid } = storeToRefs(configuration)

const extensions = [readOnlyTheme]

const code = computed(() => {
  const code = props.path ? get.value(props.path) : current.value
  switch (configuration.format) {
    case 'yaml':
      return stringify(code)
    default:
      return JSON.stringify(code, null, 2)
  }
})

const toast = useToast()
const onCopy = () => {
  copy(code.value, {
    message: 'Press #{key} to copy'
  })
  toast.success('Copied to clipboard', {
    closeOnClick: false,
    closeButton: false,
    pauseOnHover: false,
    draggable: false,
    showCloseButtonOnHover: false,
    hideProgressBar: true,
    timeout: 2000,
  })
}
</script>

<style scoped>
.editor-container {
  width: 100%;
  height: fit-content;
  position: relative;
}

.copy-button {
  position: absolute;
  /* Icon-button is 1.1rem + 2*0.3rem padding => 1.7rem / 2 => 0.85rem offset to place it in the corner */
  right: -0.85rem;
  bottom: -0.85rem;
  background: var(--vt-c-accent);
  color: var(--vt-c-white);
  border-radius: 50%;
  padding: 0.3rem;
}

.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
