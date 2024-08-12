<template>
  <div class="form-container">
    <codemirror
      v-model="code"
      class="editor"
      placeholder="Please enter valid JSON"
      :autofocus="true"
      :indent-with-tab="true"
      :tab-size="2"
      :extensions="extensions"
      :style="{ width: '100%' }"
    />
    <div class="error" v-if="errors">
      <div v-if="errors.self.length > 0">
        <template v-if="errors.self.length > 1">
          <p>Errors:</p>
          <ul>
            <li v-for="(err, i) in errors.self" :key="i">{{ err }}</li>
          </ul>
        </template>
        <p v-else>Error: {{ (errors.self ?? [])[0] }}</p>
      </div>
      <div v-if="Object.entries(errors.others).length > 0">
        <p>{{ errors.self.length > 0 ? 'And causes:' : 'Causes:' }}</p>
        <div class="other-errors">
          <template v-for="[key, value] of Object.entries(errors.others)" :key="key">
            <p>{{ key }}:</p>
            <ul>
              <li v-for="(msg, i) in value" :key="key + i">
                {{ msg }}
              </li>
            </ul>
          </template>
        </div>
      </div>
    </div>
    <button class="next" @click="onClick()">Next</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { editTheme } from '../code-snippet/editor-theme'
import { useConfigurationStore } from '@/stores/configuration'

const props = defineProps(['path'])

const configuration = useConfigurationStore()
const { set } = configuration

const code = ref('')
const extensions = [editTheme]
const errors = ref(undefined as ValidationError | undefined)

const onClick = async () => {
  const tmp = JSON.parse(code.value)
  const res = await set(props.path, tmp)
  if (!res.valid) {
    errors.value = res.errors
  } else {
    errors.value = undefined
  }
}
</script>

<style scoped></style>
