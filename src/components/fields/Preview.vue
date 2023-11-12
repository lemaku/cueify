<template>
  <div class="preview">
    <div class="inner">
      <a class="link" @click="jumpTo(props.path)">
        <p>{{ content }}</p></a
      >
      <button class="extend" @click="collapsed = !collapsed">
        <ChevronRightIcon class="icon" v-if="collapsed" />
        <ChevronDownIcon class="icon" v-if="!collapsed" />
      </button>
    </div>
    <div v-if="!collapsed">
      <CurrentConfig />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useConfigurationStore } from '@/stores/configuration'
import CurrentConfig from '@/components/CurrentConfig.vue'
import { ChevronRightIcon, ChevronDownIcon } from '@heroicons/vue/20/solid'

const props = defineProps(['path'])
const configuration = useConfigurationStore()
const { jumpTo } = configuration

const content = JSON.stringify(configuration.get(props.path))
const collapsed = ref(true)
</script>

<style scoped>
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
  /* https://lennartc.dk/en/css-text-overflow-ellipsis-not-working-or-pushing-flex-content-to-max-width/ */
  overflow: hidden;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  display: -webkit-box;
  white-space: normal;
}
.extend {
  width: fit-content;
  margin: 0;
  padding: 0;
  background: transparent;
  border: none;
  display: flex;
  padding-left: 0.3rem;
  padding-right: 0.3rem;
}

.icon {
  width: 20px;
  color: var(--color-text);
}
</style>
