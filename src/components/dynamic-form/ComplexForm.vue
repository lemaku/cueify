<template>
  <div class="form-container">
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
import Input from '@/components/fields/Input.vue'
import Preview from '@/components/fields/Preview.vue'

const configuration = useConfigurationStore()

const components = computed(() =>
  configuration.fields
    .map((f) => {
      let inputType = 'text'
      let component = Input
      switch (f.type) {
        case 'bool':
          inputType = 'checkbox'
          break
        case 'number':
          inputType = 'number'
          break
        case 'complex':
          component = Preview
          break
        case 'list':
          component = Preview
          break
      }

      return {
        type: component,
        path: f.path,
        inputType: inputType,
        index: f.index
      }
    })
    .sort((a, b) => a.index - b.index)
)
</script>

<style scoped>
.form-container {
  display: grid;
  grid-template-columns: min-content auto;
  gap: 1em;
  align-items: start;
  justify-items: start;
  height: fit-content;
}
</style>
