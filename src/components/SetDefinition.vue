<template>
  <div class="s-container">
    <h1>Webform for CUElang definitions</h1>
    <p>
      Please enter your CUE defintion in the editor below. Make sure to "export" the defition that
      should be used as the root by creating a definition called #export and setting it to the
      desired definition.
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
    <button class="next" @click="onClick()">Next</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { customTheme } from './code-snippet/editor-theme';
import { useConfigurationStore } from '@/stores/configuration'

const configuration = useConfigurationStore()
const { setSchema } = configuration

// TODO: remove default
const code = ref(`#student: {
	matNr:  string & =~"^[0-9]{8}$"
	name:   string
	active: *true | bool
    if active {
        semester: int
    }
}

#universities: {
	tuwien: {
		name: "Vienna University of Technology" | "University of Vienna",
		students: [...#student]
	},
	countryCode: string
}

#export: #universities`)
const extensions = [customTheme]

const onClick = () => {
  setSchema(code.value)
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
</style>
