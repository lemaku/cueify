<template>
  <header>
    <div class="wrapper">
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/about">About</RouterLink>
      </nav>
    </div>
  </header>

  <loading
    v-if="!wasmInitialized"
    :active="true"
    :is-full-page="true"
    :loader="'dots'"
    :color="'#4a86e8'"
  ></loading>
  <div v-else class="main">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { useGlobalStore } from '@/stores/global'
import { storeToRefs } from 'pinia'
import Loading from 'vue-loading-overlay'
import 'vue-loading-overlay/dist/css/index.css'
import { RouterLink, RouterView } from 'vue-router'

const global = useGlobalStore();
const { wasmInitialized } = storeToRefs(global)
</script>

<style scoped>
header {
  font-size: 1.5rem;
  padding: 0.5rem;
  text-align: center;
}

.main {
  width: 100%;
  max-width: 70vw;
  min-width: 800px;
  padding: 20px;
}

nav {
  text-align: center;
}

nav a {
  padding: 0 1rem;
  border-left: 1px solid var(--color-text);
}

nav a:first-of-type {
  border: 0;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}
</style>
