<template>
  <div class="choice-container">
    <button v-if="type.includes('null')" @click="onClick(null)">null</button>
    <button v-if="type.includes('bool')" @click="onClick(false)">bool</button>
    <button v-if="type.includes('int') || type.includes('float')" @click="onClick(0)">number</button>
    <button v-if="type.includes('string')" @click="onClick('')">string</button>
    <button v-if="type.includes('bytes')" @click="onClick('')">bytes</button>
    <button v-if="type.includes('struct')" @click="onClick({})">struct</button>
    <button v-if="type.includes('list')" @click="onClick([])">struct</button>
  </div>
</template>

<script setup lang="ts">
import { useConfigurationStore } from '@/stores/configuration';

const props = defineProps(['path', 'type'])
const configuration = useConfigurationStore()
const { set, jumpTo, path: parent } = configuration

const onClick = async (value: any) => {
  await set(props.path, value)
  // Do this as a way to reload after choice for type has been made and can be rendered properly
  jumpTo(parent)
}
</script>

<style scoped>
.choice-container {
  width: 100%;
  display: flex;
  justify-content: center;
  gap: 1em;
}

.choice-container > button {
  width: 80px
}
</style>
