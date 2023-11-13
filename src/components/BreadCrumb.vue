<template>
  <div class="breadcrumbs">
    <template v-for="(crumb, index) in breadcrumbs" :key="index">
      <router-link :to="{ query: { p: crumb.path.join('.') } }">
        {{ crumb.crumb }}
      </router-link>
      <template v-if="index != breadcrumbs.length - 1"> > </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import { storeToRefs } from 'pinia'

const configuration = useConfigurationStore()
const { breadcrumbs } = storeToRefs(configuration)
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
</style>
