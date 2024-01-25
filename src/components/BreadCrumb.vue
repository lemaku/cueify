<template>
  <div class="breadcrumbs">
    <template v-if="breadcrumbs.length > 1">
      <template v-for="(crumb, index) in breadcrumbs" :key="index">
        <a
          @click="isEqual(path, crumb.path) ? () => {} : jumpTo(crumb.path)"
          :class="{ disabled: isEqual(path, crumb.path) }"
          >{{ crumb.crumb }}</a
        >
        <template v-if="index != breadcrumbs.length - 1"> > </template>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import { isEqual } from 'lodash'
import { storeToRefs } from 'pinia'

const configuration = useConfigurationStore()
const { jumpTo } = configuration
const { breadcrumbs, path } = storeToRefs(configuration)
</script>

<style scoped>
.breadcrumbs {
  display: flex;
  align-items: center;
  gap: 10px;
}

a {
  cursor: pointer;
}

.disabled {
  cursor: default;
  color: var(--color-text);
}
</style>
