<template>
  <div class="list-container">
    <template v-if="shouldDisplayChoice">
      <Choice :path="path" :type="currentOf" :pushToArray="true" />
    </template>

    <template v-else>
      <template v-if="isEqual(currentOf, ['struct']) || isEqual(currentOf, ['list'])">
        <Preview v-for="field in fields" :key="field.path.join('.')" :path="field.path" />
      </template>
      <template v-else>
        <Input
          v-for="field in fields"
          :key="field.path.join('.')"
          :path="field.path"
          :type="field.type"
        />
      </template>

      <button class="icon-button" @click="handleAdd()">
        <i class="pi pi-plus-circle"></i>
      </button>
    </template>
  </div>
</template>

<script setup lang="ts">
import Choice from '@/components/fields/Choice.vue'
import Input from '@/components/fields/Input.vue'
import Preview from '@/components/fields/Preview.vue'
import { useConfigurationStore } from '@/stores/configuration'
import { isEqual, isNull } from 'lodash'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'

const configuration = useConfigurationStore()
const { fields, currentOf, path } = storeToRefs(configuration)
const { addToArray } = configuration

const handleAdd = () => {
  let toPush
  switch (currentOf.value![0]) {
    case 'null':
      toPush = null
      break
    case 'bool':
      toPush = false
      break
    case 'int':
      toPush = 0
      break
    case 'float':
      toPush = 0
      break
    case 'string':
      toPush = ''
      break
    case 'struct':
      toPush = {}
      break
    case 'list':
      toPush = []
      break
    default:
      break
  }
  addToArray(toPush)
}

const shouldDisplayChoice = computed(
  () =>
    (isNull(fields.value) || isEqual(fields.value?.length, 0)) &&
    currentOf.value != null &&
    currentOf.value?.length > 1
)
</script>

<style scoped>
.list-container {
  display: flex;
  flex-direction: column;
  gap: 1em;
}
</style>
