<template>
  <template v-for="(field, index) in fields" :key="index">
    <Preview :path="field.path"></Preview>
  </template>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import Preview from '@/components/fields/Preview.vue'
import { storeToRefs } from 'pinia'

const configuration = useConfigurationStore()

const { fields } = storeToRefs(configuration)
</script>

<style scoped></style>
