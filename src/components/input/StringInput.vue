<template>
  <input class="form-control" type="text" placeholder="text here" v-model="curVal" @change="set(props.path, curVal)" />
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useConfigurationStore } from '../../stores/configuration'

const props = defineProps(['path'])
const configuration = useConfigurationStore()
const { get } = storeToRefs(configuration)
const { set } = configuration

const curVal = get.value(props.path)
</script>
