<template>
  <div class="format-selector">
    <label for="format">Format</label>
    <select v-model="selected" @change="changeFormat(selected)" id="format">
      <option v-for="opt in values" :value="opt" v-bind:key="opt">{{ opt }}</option>
    </select>
  </div>
</template>

<script setup lang="ts">
import type { Ref } from 'vue'
import { ref } from 'vue'
import { supportedFormats, useConfigurationStore, type Format } from '../stores/configuration'

const configuration = useConfigurationStore()

const values = supportedFormats
const selected: Ref<Format> = ref('json')
const { changeFormat } = configuration
</script>

<style scoped>
.format-selector {
  display: flex;
  justify-content: end;
  align-items: center;
  gap: 10px
}

select {
  width: fit-content;
}
</style>
