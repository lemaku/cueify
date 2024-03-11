<template>
  <div class="choice-container">
    <button v-if="type.includes('null')" @click="onClick(null)">null</button>
    <button v-if="type.includes('bool')" @click="onClick(false)">bool</button>
    <button v-if="type.includes('int') || type.includes('float')" @click="onClick(0)">
      number
    </button>
    <button v-if="type.includes('string')" @click="onClick('')">string</button>
    <button v-if="type.includes('struct')" @click="onClick({})">struct</button>
    <button v-if="type.includes('list')" @click="onClick([])">list</button>
  </div>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration'

const props = defineProps(['path', 'type', 'pushToArray'])
const configuration = useConfigurationStore()
const { set, addToArray } = configuration

const onClick = async (value: any) => {
  if (props.pushToArray) {
    await addToArray(value)
  } else {
    await set(props.path, value)
  }
}
</script>

<style scoped>
.choice-container {
  width: 100%;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1em;
}

.choice-container > button {
  width: 80px;
}
</style>
