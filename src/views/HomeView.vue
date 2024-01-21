<template>
  <div class="container">
    <BreadCrumb class="breadcrumb"></BreadCrumb>
    <DynamicFormContainer class="form"></DynamicFormContainer>

    <FormatSelector class="format"></FormatSelector>
    <CodeSnippet class="current"></CodeSnippet>

    <Errors class="errors" />
  </div>
</template>

<script setup lang="ts">
import BreadCrumb from '@/components/BreadCrumb.vue'
import DynamicFormContainer from '@/components/DynamicFormContainer.vue'
import Errors from '@/components/Errors.vue'
import FormatSelector from '@/components/FormatSelector.vue'
import CodeSnippet from '@/components/code-snippet/CodeSnippet.vue'
</script>

<style>
.container {
  display: grid;
  grid-template-rows: min-content auto auto;
  grid-template-columns: 50% 50%;
  gap: 5px 20px;
}

.breadcrumb {
  grid-row: 1;
  grid-column: 1;
}

.format {
  grid-row: 1;
  grid-column: 2;
}

.form {
  grid-row: 2;
  grid-column: 1;
}

.current {
  grid-row: 2;
  grid-column: 2;
}

.errors {
  margin-top: 1rem;
  grid-row: 3;
  grid-column: 1 / span 2;
}
</style>
