<template>
  <div id="form-container">
    <template v-for="(field, index) in components" :key="index">
      <label :for="field.path.join('.')">{{ field.path[field.path.length - 1] }}</label>
      <component :is="field.type" :path="field.path" :id="field.path.join('.')"></component>
    </template>

    <!-- <StringInput path="prop"></StringInput> -->
    <!-- <p>this is where input fields are going to be dynamically rendered to</p> -->
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useConfigurationStore } from '../stores/configuration'
import StringInput from './input/StringInput.vue'
import BoolInput from './input/BoolInput.vue'
import NumberInput from './input/NumberInput.vue'

const configuration = useConfigurationStore()

const components = computed(() =>
  configuration.fields.map((f) => {
    let type = StringInput
    switch (f.type) {
      case 'bool':
        type = BoolInput
        break
      case 'number':
        type = NumberInput
    }
    
    return {
      path: f.path,
      type: type
    }
  })
)
</script>

<style scoped>
.form-container {
  display: grid;
  grid-template-columns: 1fr 5fr;
  gap: 1em;
  align-items: center;
  justify-items: start;
  height: fit-content;
}
</style>
