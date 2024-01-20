<template>
  <ComplexForm v-if="currentType == 'complex'" />
  <ListForm v-else-if="currentType == 'list'" />
  <div v-else>Error</div>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import { storeToRefs } from 'pinia'
import ComplexForm from '@/components/dynamic-form/ComplexForm.vue'
import ListForm from '@/components/dynamic-form/ListForm.vue'

const configuration = useConfigurationStore()

const { currentType } = storeToRefs(configuration)
</script>
