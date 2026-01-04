<script setup lang="ts">
import { ExternalLink } from 'lucide-vue-next'

defineProps<{
  widget: any
}>()
</script>

<template>
  <div class="flex items-center justify-between text-xs py-1">
    <span class="text-zinc-500">{{ widget.label }}</span>
    <a 
      :href="widget.uri" 
      target="_blank"
      class="flex items-center gap-1.5 text-blue-400 hover:text-blue-300 hover:underline transition-colors"
    >
      <span class="font-medium">{{ widget.text || 'Open' }}</span>
      <ExternalLink class="w-3 h-3" />
    </a>
  </div>
</template>
