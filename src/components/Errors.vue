<template>
  <div class="errors-container" v-if="(errors ?? []).length > 0">
    <p class="errors-header">Resolve next</p>
    <template v-for="error in errors" v-bind:key="error.path.join('.')">
      <a @click="jumpTo(error.path)">{{ error.path.join('.') }}</a>
      <ul v-if="(error.errors ?? []).length > 1">
        <li v-for="err in error.errors" v-bind:key="err">
          {{ err }}
        </li>
      </ul>
      <div v-else>
        {{ (error.errors ?? [])[0] }}
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import { computed } from 'vue'

const configuration = useConfigurationStore()
const { jumpTo } = configuration;

const errors = computed(() => {
  return [...configuration.errors].sort((a, b) => a.path.length - b.path.length)
})
</script>

<style scoped>
.errors-container {
  width: 100%;
  display: grid;
  grid-template-columns: min-content auto;
  gap: 1rem;
}
.errors-header {
  grid-column: 1 / span 2;
  color: var(--vt-c-grey);
  font-size: 1rem;
  font-weight: bold;
}

a {
  cursor: pointer;
}
</style>
