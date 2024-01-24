<template>
  <div class="list-container">
    <Preview v-for="field in fields" :key="field.path.join('.')" :path="field.path"></Preview>
    <button class="icon-button" @click="addToArray()">
      <i class="pi pi-plus-circle"></i>
    </button>
  </div>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import Preview from '@/components/fields/Preview.vue'
import { storeToRefs } from 'pinia'

const configuration = useConfigurationStore()

const { addToArray } = configuration

const { fields } = storeToRefs(configuration)
</script>

<style scoped>
.list-container {
  display: flex;
  flex-direction: column;
  gap: 1em;
}
</style>
