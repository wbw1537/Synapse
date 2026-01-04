<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'

const props = defineProps<{
  widget: any,
  autoHeight?: boolean
}>()

const logRef = ref<HTMLElement | null>(null)

const logs = computed(() => {
  if (Array.isArray(props.widget.value)) {
    return props.widget.value
  }
  return [String(props.widget.value)]
})

// Auto-scroll to bottom on update
watch(logs, () => {
  nextTick(() => {
    if (logRef.value) {
      logRef.value.scrollTop = logRef.value.scrollHeight
    }
  })
}, { deep: true })
</script>

<template>
  <div class="flex flex-col gap-1 py-1">
    <span class="text-xs text-zinc-500">{{ widget.label }}</span>
    <div 
      ref="logRef"
      class="bg-zinc-950 rounded border border-zinc-800 p-2 text-[10px] font-mono text-zinc-400 overflow-y-auto whitespace-pre-wrap transition-all"
      :class="autoHeight ? 'min-h-24' : 'h-24'"
    >
      <div v-for="(line, idx) in logs" :key="idx" class="border-b border-zinc-800/50 last:border-0 pb-0.5 mb-0.5 last:pb-0 last:mb-0">
        {{ line }}
      </div>
    </div>
  </div>
</template>
