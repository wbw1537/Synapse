<script setup lang="ts">
import { Copy } from 'lucide-vue-next'
import { ref } from 'vue'

const props = defineProps<{
  widget: any
}>()

const copied = ref(false)

const copyToClipboard = async () => {
  if (!props.widget.copyable) return
  try {
    await navigator.clipboard.writeText(String(props.widget.value))
    copied.value = true
    setTimeout(() => copied.value = false, 2000)
  } catch (err) {
    console.error('Failed to copy', err)
  }
}
</script>

<template>
  <div class="flex items-center justify-between text-xs py-1">
    <span class="text-zinc-500">{{ widget.label }}</span>
    <div 
      class="flex items-center gap-2" 
      :class="{'cursor-pointer hover:text-zinc-100': widget.copyable}"
      @click="copyToClipboard"
    >
      <span class="font-mono text-zinc-300">
        {{ widget.value }}<span v-if="widget.unit" class="text-zinc-500 ml-0.5">{{ widget.unit }}</span>
      </span>
      <Copy v-if="widget.copyable && !copied" class="w-3 h-3 text-zinc-600" />
      <span v-if="copied" class="text-[10px] text-emerald-500 font-bold">COPIED</span>
    </div>
  </div>
</template>
