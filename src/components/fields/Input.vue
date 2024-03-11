<template>
  <div class="form-input">
    <input
      class="form-control"
      :class="{ 'input-error': errors, 'input-success': success }"
      :type="type"
      v-model="curVal"
      :step="step"
      ref="input"
      @change="onChange()"
      @focus="onFocus()"
      @focusout="onFocusOut()"
    />
    <button
      class="icon-button"
      :disabled="currentType !== 'list' && (isUndefined || isDerived)"
      @click="onClear()"
    >
      <i class="pi pi-eraser"></i>
    </button>
    <div class="error" v-if="errors">
      <div v-if="errors.self.length > 0">
        <template v-if="errors.self.length > 1">
          <p>Errors:</p>
          <ul>
            <li v-for="(err, i) in errors.self" :key="i">{{ err }}</li>
          </ul>
        </template>
        <p v-else>Error: {{ (errors.self ?? [])[0] }}</p>
      </div>
      <div v-if="Object.entries(errors.others).length > 0">
        <p>{{ errors.self.length > 0 ? 'And causes:' : 'Causes:' }}</p>
        <div class="other-errors">
          <template v-for="[key, value] of Object.entries(errors.others)" :key="key">
            <p>{{ key }}:</p>
            <ul>
              <li v-for="(msg, i) in value" :key="key + i">
                {{ msg }}
              </li>
            </ul>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import emitter from '@/event-bus'
import { useConfigurationStore } from '@/stores/configuration'
import { debounce, isArray, isEqual } from 'lodash'
import { computed, onMounted, ref } from 'vue'

const props = defineProps(['path', 'type', 'placeholder'])
const configuration = useConfigurationStore()
const { set, unset, get, isDerived: derived, currentType } = configuration

let curVal = ref(get(props.path))
const isUndefined = computed(() => get(props.path) === undefined)
const isDerived = computed(() => derived(props.path))

const errors = ref(undefined as ValidationError | undefined)
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

// Use ref to be able to set the focus to the input field
const input = ref()
onMounted(() => {
  // List to focus event and use reference to input field in case it was called with my path
  emitter.$on('focus', (args: string[]) => {
    if (args && isArray(args) && isEqual(args, props.path) && input.value) {
      input.value.focus()
    }
  })

  // Subsribe to all changes and fetch my value again in case it changed
  configuration.$onAction(({ name, after }) => {
    if (name === 'set') {
      after(({ valid }) => {
        if (valid) {
          const tmp = curVal.value
          if (tmp !== get(props.path)) {
            curVal.value = get(props.path)
            errors.value = undefined
            success.value = false
          }
        }
      })
    }
  })
})
</script>

<style scoped>
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
.other-errors {
  width: 100%;
  display: grid;
  grid-template-columns: min-content auto;
}
</style>
