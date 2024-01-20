<template>
  <div class="form-input">
    <input
      class="form-control"
      v-bind:class="{ 'input-error': errors, 'input-success': success }"
      v-bind:type="type"
      v-model="curVal"
      @change="onChange()"
      @focus="onFocus()"
      @focusout="onFocusOut()"
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
import { useConfigurationStore } from '@/stores/configuration'
import { TrashIcon } from '@heroicons/vue/20/solid'
import { debounce } from 'lodash'
import { Mutex } from 'async-mutex'

const props = defineProps(['path', 'type', 'placeholder'])
const configuration = useConfigurationStore()
const { set, unset, get } = configuration

let curVal = ref(get(props.path))
const isUndefined = computed(() => get(props.path) === undefined)
const errors = ref(undefined as string[] | undefined)
const success = ref(false)
let focus = false

const mutex = new Mutex()

const onClear = async () => {
  await mutex.acquire()
  await unset(props.path)
  curVal.value = get(props.path)
  success.value = false
  errors.value = undefined
  mutex.release()
}

const onChange = debounce(async () => {
  console.log('change')
  await mutex.acquire()
  // This is needed because if the value is changed (but not commited) and then the clear button is pressed,
  // first the clear is triggered, letting the onChange (that is triggered now because the focus left the field)
  // wait for the mutex and when it gets it, it would try to set the value to curVal.val which now is undefined
  // because of the clear
  if (curVal.value != undefined) {
    const res = await set(props.path, curVal.value)
    if (!res.valid) {
      errors.value = res.errors
      success.value = false
    } else {
      errors.value = undefined
      if (focus) {
        success.value = true
      }
    }
  }
  mutex.release()
}, 150)

const onFocus = async () => {
  await mutex.acquire()
  focus = true
  console.log(focus);
  mutex.release()
}

const onFocusOut = async () => {
  await mutex.acquire()
  focus = false
  success.value = false
  console.log(focus)
  mutex.release()
}
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
