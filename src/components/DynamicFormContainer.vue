<template>
  <div id="form-container">
    <template v-for="(field, index) in components" :key="index">
      <label :for="field.path.join('.')">{{ field.path[field.path.length - 1] }}</label>
      <component
        :is="field.type"
        :type="field.inputType"
        :path="field.path"
        :id="field.path.join('.')"
      ></component>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useConfigurationStore } from '@/stores/configuration'
import Input from './input/Input.vue'

const configuration = useConfigurationStore()

const components = computed(() =>
  configuration.fields.map((f) => {
    let inputType = 'text'
    switch (f.type) {
      case 'bool':
        inputType = 'checkbox'
        break
      case 'number':
        inputType = 'number'
        break
    }

    return {
      type: Input,
      path: f.path,
      inputType: inputType
    }
  })
)
</script>

<style scoped>
.form-container {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 1em;
  align-items: start;
  justify-items: start;
  height: fit-content;
}
</style>
