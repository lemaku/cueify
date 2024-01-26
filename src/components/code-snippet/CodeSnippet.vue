<template>
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
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { stringify } from 'yaml'
import { Codemirror } from 'vue-codemirror'
import { readOnlyTheme } from './editor-theme'

const props = defineProps(['path', 'depth'])

const configuration = useConfigurationStore()
const { get, current } = storeToRefs(configuration)

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
</script>
