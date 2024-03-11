<template>
  <div class="s-container">
    <h1>Web form for CUE definitions</h1>
    <p>
      Please enter your CUE defintion in the editor below. Make sure to name your definition "export".
    </p>
    <codemirror
      v-model="code"
      class="editor"
      placeholder="CUE definition goes here..."
      :autofocus="true"
      :indent-with-tab="true"
      :tab-size="2"
      :extensions="extensions"
      :style="{ width: '100%' }"
    />
    <div v-if="error" class="error">Error: {{ error }}</div>
    <button class="next" @click="onClick()">Next</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { editTheme } from './code-snippet/editor-theme';
import { useConfigurationStore } from '@/stores/configuration'

const configuration = useConfigurationStore()
const { setSchema } = configuration

// TODO: remove default
const code = ref(`#export: {
  x: int
  y: x
  z: y
  p?: { type: "type1", eitherThis: string } | { type: "type2", orThat: string }
  a: [...string]
  b: [...string] | [...int]
  c: [...{a: string}]
  d: [...[...string]]
}`)
const extensions = [editTheme]
const error = ref(undefined as string | undefined);

const onClick = () => {
  const res = setSchema(code.value);
  if (!res.valid) {
    error.value = res.error;
  }
}
</script>

<style scoped>
.s-container {
  display: flex;
  flex-direction: column;
  gap: 0.5em;
}
.editor {
  align-self: center;
}
.next {
  align-self: flex-end;
}
.error {
  color: var(--vt-c-error);
}
</style>
