<template>
  <ComplexForm v-if="currentType == 'complex'" />
  <ListForm v-if="currentType == 'list'" />
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import { storeToRefs } from 'pinia'
import ComplexForm from '@/components/dynamic-form/ComplexForm.vue'
import ListForm from '@/components/dynamic-form/ListForm.vue'

const configuration = useConfigurationStore()

const { currentType } = storeToRefs(configuration)
</script>
