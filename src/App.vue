<template>
  <header>
    <div class="wrapper">
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/about">About</RouterLink>
      </nav>
    </div>
  </header>

  <div class="main">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import router from '@/router'
import { useConfigurationStore } from '@/stores/configuration'
import { watch } from 'vue'
import { RouterLink, RouterView } from 'vue-router'

const configuration = useConfigurationStore()
const { jumpTo } = configuration

jumpTo((router.currentRoute.value.query.p ?? 'universities').toString().split('.'))
configuration.summarize()

watch(router.currentRoute, () => {
  const path = (router.currentRoute.value.query.p ?? 'universities').toString()
  jumpTo(path.split('.'))
})
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
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}
</style>
