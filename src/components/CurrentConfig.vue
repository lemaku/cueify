<template>
  <div v-html="compiled" class="markdown" />
</template>

<script setup lang="ts">
import hljs from 'highlight.js/lib/core'
import jsonHighlight from 'highlight.js/lib/languages/json'
import textHightlight from 'highlight.js/lib/languages/plaintext'
import yamlHighlight from 'highlight.js/lib/languages/yaml'
import 'highlight.js/styles/github.css'
import { Marked } from 'marked'
import { markedHighlight } from 'marked-highlight'
import { computed } from 'vue'
import { stringify } from 'yaml'
import { useConfigurationStore } from '../stores/configuration'

hljs.registerLanguage('json', jsonHighlight)
hljs.registerLanguage('plaintext', textHightlight)
hljs.registerLanguage('yaml', yamlHighlight)

const configuration = useConfigurationStore()

const marked = new Marked(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext'
      return hljs.highlight(code, { language }).value
    }
  })
)

const compiled = computed(() => {
  let snippet = ''
  switch (configuration.format) {
    case 'json':
      snippet = JSON.stringify(configuration.current, null, 2)
      break
    case 'yaml':
      snippet = stringify(configuration.current)
      break
    default:
      snippet = snippet = JSON.stringify(configuration.current, null, 2)
      break
  }

  //TODO inject a copy button into code html
  return marked.parse('```' + configuration.format + '\n' + snippet + '\n```')
})
</script>

<style>
code.hljs {
  border-radius: 10px;
  padding-right: 10em !important;
}
</style>
