<template>
  <div class="form-input">
    <input
      class="form-control"
      :class="{ 'input-error': errors, 'input-success': success }"
      :type="type"
      v-model="curVal"
      :step="step"
      @change="onChange()"
      @focus="onFocus()"
      @focusout="onFocusOut()"
    />
    <button class="icon-button" :disabled="isUndefined" @click="onClear()">
      <i class="pi pi-eraser"></i>
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
import { debounce } from 'lodash'

const props = defineProps(['path', 'type', 'placeholder'])
const configuration = useConfigurationStore()
const { set, unset, get } = configuration

let curVal = ref(get(props.path))
const isUndefined = computed(() => get(props.path) === undefined)

const errors = ref(undefined as string[] | undefined)
const success = ref(false)
let focus = false

const onClear = async () => {
  await unset(props.path)
  curVal.value = get(props.path)
  success.value = false
  errors.value = undefined
}

const onChange = debounce(async () => {
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
}, 0)

const onFocus = async () => {
  focus = true
}

const onFocusOut = async () => {
  focus = false
  success.value = false
}

const step = computed(() => {
  const value = curVal.value
  if (!value) return 1
  if (Math.floor(value.valueOf()) === value.valueOf()) return 1

  var str = value.toString()
  if (str.indexOf('.') !== -1) {
    return Math.pow(10, -1 * (str.split('.')[1].length || 1))
  }
  return 1
})
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
