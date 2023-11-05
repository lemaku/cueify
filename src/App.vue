<template>
  <header>
    <div class="wrapper">
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/about">About</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView class="main" />
</template>

<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { useConfigurationStore } from './stores/configuration'

const configuration = useConfigurationStore()
configuration.load()
</script>

<style scoped>
header {
  font-size: 1.5rem;
  padding: 0.5rem;
  text-align: center;
}

main {
  width: 100%;
  max-width: 1500px;
  min-width: 800px;
  padding: 20px;
  box-sizing: border-box;
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
