<template>
  <div class="form-input">
    <input
      class="form-control"
      v-bind:class="{ 'input-error': errors, 'input-success': success }"
      v-bind:type="type"
      v-model="curVal"
      @change="onChange()"
      @focusout="success = false"
    />
    <button class="icon-button" :disabled="isUndefined" @click="onClear()">
      <TrashIcon class="icon-20" />
    </button>
    <div class="error" v-if="errors">
      <template v-for="(err, i) in errors" :key="i">
        {{ err }}
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useConfigurationStore } from '@/stores/configuration'
import { TrashIcon } from '@heroicons/vue/20/solid'
import { debounce } from 'lodash'
import { Mutex } from 'async-mutex';

const props = defineProps(['path', 'type', 'placeholder'])
const configuration = useConfigurationStore()
const { get } = storeToRefs(configuration)
const { set, unset } = configuration

const curVal = ref(get.value(props.path))
const isUndefined = computed(() => curVal.value === undefined)
const errors = ref(undefined as string[] | undefined)
const success = ref(false)

const mutex = new Mutex();

const onClear = async () => {
  await mutex.acquire()
  await unset(props.path)
  curVal.value = undefined
  success.value = false
  errors.value = undefined
  mutex.release()
}

const onChange = debounce(async () => {
  await mutex.acquire();
  const res = await set(props.path, curVal.value);
  if (!res.valid) {
    errors.value = res.errors
    success.value = false
  } else {
    errors.value = undefined
    success.value = true
  }
  mutex.release();
}, 150)
</script>

<style>
.form-input {
  display: grid;
  grid-template-columns: minmax(80%, 100%) auto;
  width: 100%;
}
.error {
  padding: 0.5rem 0.7rem 0.5rem;
  line-height: 1rem;
  font-size: 0.8rem;
}
</style>
