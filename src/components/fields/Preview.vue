<template>
  <div class="preview-container" v-if="!isUndefined">
    <div class="preview">
      <div class="inner">
        <router-link class="link" :to="{ query: { p: props.path.join('.') } }">
          {{ content }}</router-link
        >
        <button class="icon-button" @click="collapsed = !collapsed">
          <i class="pi pi-chevron-right" style="font-size: 0.9rem;" v-if="collapsed"></i>
          <i class="pi pi-chevron-down" style="font-size: 0.9rem;" v-if="!collapsed"></i>
        </button>
      </div>
      <div v-if="!collapsed && !isUndefined">
        <CodeSnippet :path="props.path" class="snippet" />
      </div>
    </div>
    <button class="icon-button button-preview" @click="unset(path)">
      <i class="pi pi-trash"></i>
    </button>
  </div>
  <button v-else class="icon-button button-preview" @click="setToEmpty(path, isArray)">
    <i class="pi pi-plus-circle"></i>
  </button>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useConfigurationStore } from '@/stores/configuration'
import CodeSnippet from '@/components/code-snippet/CodeSnippet.vue'

const props = defineProps(['path'])
const configuration = useConfigurationStore()

const { unset, get, setToEmpty, fields } = configuration

const isArray = fields.find((f) => f.path === props.path)?.type === 'list'

const content = computed(() => JSON.stringify(configuration.get(props.path)))
const collapsed = ref(true)

const isUndefined = computed(() => get(props.path) === undefined)
</script>

<style scoped>
.preview-container {
  width: 100%;
  display: grid;
  grid-template-columns: minmax(80%, 100%) auto;
  align-items: start;
}
.button-preview {
  height: 2rem;
  margin-top: 0.1rem;
}

.preview {
  width: 100%;
  height: fit-content;
  border: 0.1rem solid var(--vt-c-grey);
  border-radius: 0.4rem;
  color: var(--color-text);
}

.inner {
  height: 2rem;
  display: flex;
  justify-content: space-between;
}

.link {
  flex-grow: 1;
  color: var(--color-text);
  cursor: pointer;
  line-height: 2rem;
  font-size: 0.8rem;
  padding-left: 0.5rem;
  padding-right: 0.5rem;
  /* https://lennartc.dk/en/css-text-overflow-ellipsis-not-working-or-pushing-flex-content-to-max-width/
  overflow: hidden;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  display: -webkit-box;
  white-space: normal; */
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}
.snippet {
  font-size: 0.8rem;
}
</style>
