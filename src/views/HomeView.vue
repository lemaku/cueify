<script setup lang="ts">
import DynamicFormContainer from '@/components/DynamicFormContainer.vue'
import CurrentConfig from '@/components/CurrentConfig.vue'
import BreadCrumb from '@/components/BreadCrumb.vue'
import FormatSelector from '@/components/FormatSelector.vue'
</script>

<template>
  <div class="container">
    <BreadCrumb class="breadcrumb"></BreadCrumb>
    <DynamicFormContainer class="form-container"></DynamicFormContainer>

    <FormatSelector class="format"></FormatSelector>
    <CurrentConfig class="current"></CurrentConfig>
  </div>
</template>

<style>
.container {
  display: grid;
  grid-template-rows: min-content auto;
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

.form-container {
  grid-row: 2;
  grid-column: 1;
}

.current {
  grid-row: 2;
  grid-column: 2;
}
</style>
