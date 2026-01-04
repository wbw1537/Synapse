<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  widget: any
}>()

const percentage = computed(() => {
  const min = props.widget.min || 0
  const max = props.widget.max || 100
  const val = props.widget.value || 0
  const p = ((val - min) / (max - min)) * 100
  return Math.min(Math.max(p, 0), 100)
})

const colorClass = computed(() => {
  if (!props.widget.thresholds) return 'bg-blue-500'
  
  const val = props.widget.value || 0
  let color = 'bg-emerald-500' // default good
  
  // Sort thresholds
  const sortedKeys = Object.keys(props.widget.thresholds)
    .map(Number)
    .sort((a, b) => a - b)
    
  for (const k of sortedKeys) {
    if (val >= k) {
      const c = props.widget.thresholds[String(k)]
      switch(c) {
        case 'warning': color = 'bg-amber-500'; break;
        case 'danger': color = 'bg-rose-500'; break;
        case 'success': color = 'bg-emerald-500'; break;
        case 'primary': color = 'bg-blue-500'; break;
        default: color = 'bg-zinc-500';
      }
    }
  }
  return color
})
</script>

<template>
  <div class="flex flex-col gap-1 py-1">
    <div class="flex justify-between text-xs">
      <span class="text-zinc-500">{{ widget.label }}</span>
      <span class="font-mono text-zinc-300">
        {{ widget.value }}{{ widget.unit }}
      </span>
    </div>
    <div class="h-1.5 w-full bg-zinc-800 rounded-full overflow-hidden">
      <div 
        class="h-full transition-all duration-500 ease-out rounded-full" 
        :class="colorClass" 
        :style="{ width: `${percentage}%` }"
      ></div>
    </div>
  </div>
</template>
