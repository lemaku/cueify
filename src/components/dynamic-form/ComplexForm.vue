<template>
  <div class="form-container">
    <template v-for="field in components" :key="field.path.join('.')">
      <label class="form-label" :title="field.path[field.path.length - 1]+(field.optional ? '?' : '')" :for="field.path.join('.')"
        >{{ field.path[field.path.length - 1] }}{{ field.optional ? '?' : '' }}</label
      >
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
import Choice from '@/components/fields/Choice.vue'
import Input from '@/components/fields/Input.vue'
import NullInput from '@/components/fields/NullInput.vue'
import Preview from '@/components/fields/Preview.vue'
import { useConfigurationStore } from '@/stores/configuration'
import { computed } from 'vue'

const configuration = useConfigurationStore()

const components = computed(() =>
  configuration.fields.map((f) => {
    let inputType
    let component
    if (f.type.length === 1) {
      switch (f.type[0]) {
        case 'null':
          component = NullInput
          break
        case 'bool':
          inputType = 'checkbox'
          component = Input
          break
        case 'string':
          inputType = 'text'
          component = Input
          break
        case 'bytes':
          inputType = 'text'
          component = Input
          break
        case 'int':
          inputType = 'number'
          component = Input
          break
        case 'float':
          inputType = 'number'
          component = Input
          break
        case 'struct':
          component = Preview
          break
        case 'list':
          component = Preview
          break
        default:
          console.error('Input for ', f.type[0], 'not implemented')
          break
      }
    } else {
      inputType = f.type
      component = Choice
    }

    return {
      type: component,
      path: f.path,
      inputType: inputType,
      optional: f.optional
    }
  })
)
</script>

<style scoped>
.form-container {
  display: grid;
  grid-template-columns: 30% 70%;
  grid-auto-rows: max-content;
  grid-gap: 1em;
  align-items: start;
  justify-items: start;
}
.form-label {
  width: 100%;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}
</style>
