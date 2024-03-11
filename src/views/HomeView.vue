<template>
  <SetDefinition v-if="!schemaSet" />

  <div v-else class="f-container">
    <BreadCrumb class="f-breadcrumb" />
    <DynamicFormContainer class="f-form" />

    <FormatSelector class="f-format" />
    <CodeSnippet class="f-current" />

    <Errors class="f-errors" />
  </div>
</template>

<script setup lang="ts">
import BreadCrumb from '@/components/BreadCrumb.vue'
import DynamicFormContainer from '@/components/DynamicFormContainer.vue'
import Errors from '@/components/Errors.vue'
import FormatSelector from '@/components/FormatSelector.vue'
import CodeSnippet from '@/components/code-snippet/CodeSnippet.vue'
import SetDefinition from '@/components/SetDefinition.vue'
import { useConfigurationStore } from '@/stores/configuration'
import { storeToRefs } from 'pinia'

const configuration = useConfigurationStore()
const { schemaSet } = storeToRefs(configuration)
</script>

<style>
.f-container {
  display: grid;
  grid-template-rows: min-content auto auto;
  grid-template-columns: calc((100% - 20px) / 2) calc((100% - 20px) / 2);
  gap: 5px 20px;
}

.f-breadcrumb {
  grid-row: 1;
  grid-column: 1;
}

.f-format {
  grid-row: 1;
  grid-column: 2;
}

.f-form {
  grid-row: 2;
  grid-column: 1;
}

.f-current {
  grid-row: 2;
  grid-column: 2;
}

.f-errors {
  margin-top: 1rem;
  grid-row: 3;
  grid-column: 1 / span 2;
}
</style>
