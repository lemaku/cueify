<template>
  <ComplexForm v-if="currentType == 'struct'" />
  <ListForm v-else-if="currentType == 'list'" />
  <FreeForm v-else-if="currentType == 'freeform'" :path="path" />
  <div v-else>Error</div>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import { storeToRefs } from 'pinia'
import ComplexForm from '@/components/dynamic-form/ComplexForm.vue'
import ListForm from '@/components/dynamic-form/ListForm.vue'
import FreeForm from '@/components/dynamic-form/FreeForm.vue'

const configuration = useConfigurationStore()

const { currentType, path } = storeToRefs(configuration)
</script>
