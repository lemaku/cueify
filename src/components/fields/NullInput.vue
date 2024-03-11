<template>
  <template v-if="isUndefined">
    <button class="icon-button button-null-input" @click="set(path, null)">
      <i class="pi pi-plus-circle"></i>
    </button>
  </template>
  <template v-else>
    <div class="form-input">
      <input class="form-control read-only" v-model="curVal" disabled />
      <button
        class="icon-button"
        :disabled="currentType !== 'list' && (isUndefined || isDerived)"
        @click="onClear()"
      >
        <i class="pi pi-eraser"></i>
      </button>
    </div>
  </template>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'
import { computed } from 'vue'

const props = defineProps(['path', 'type', 'placeholder'])
const configuration = useConfigurationStore()
const { set, unset, get, currentType, isDerived: derived } = configuration

const isDerived = computed(() => derived(props.path))
const isUndefined = computed(() => get(props.path) === undefined)
const curVal = computed(() => (isUndefined.value ? '' : 'null'))

const onClear = async () => {
  await unset(props.path)
}
</script>

<style scoped>
.button-null-input {
  height: 2rem;
  margin-top: 0.1rem;
}
.form-input {
  display: grid;
  grid-template-columns: minmax(80%, 100%) auto;
  width: 100%;
}

.read-only {
  opacity: 60%;
  cursor: not-allowed;
}
</style>
