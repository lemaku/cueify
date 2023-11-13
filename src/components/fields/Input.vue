<template>
  <div class="form-input">
    <input class="form-control" 
      v-bind:class=" { 'input-error': errors, 'input-success': success }" 
      v-bind:type="props.type"
      v-model="curVal" 
      @change="onChange(props.path, curVal)"
      @focusout="success=false" />
    <div class="error" v-if="errors">
      <template v-for="(err, i) in errors" :key="i">
        {{ err }}
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useConfigurationStore } from '@/stores/configuration'

const props = defineProps(['path', 'type', 'placeholder'])
const configuration = useConfigurationStore()
const { get } = storeToRefs(configuration)
const { set } = configuration

const curVal = ref(get.value(props.path))

let errors = ref(undefined as string[] | undefined);
// let timeouts: number[] = [];
let success = ref(false);

const onChange = async (path: string[], val: string) => {
  // timeouts.forEach(clearTimeout);
  const res = await set(path, val)
  if (!res.valid) {
    errors.value = res.errors
    success.value = false;

    // timeouts.push(setTimeout(() => {
    //   curVal.value = get.value(props.path)
    //   errors.value = undefined
    // }, 5000));
  } else {
    errors.value = undefined;
    success.value = true;
  }
}
</script>

<style>
.form-input {
  display: flex;
  flex-direction: column;
  width: 100%;
}
.form-input > * {
  width: 100%;
}
.error {
  padding: .5rem .7rem .5rem;
  line-height: 1rem;
  font-size: 0.8rem;
}
</style>