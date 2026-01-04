<script setup lang="ts">
import { useServiceStore } from '../../stores/services'

const props = defineProps<{
  widget: any,
  serviceId?: string
}>()

const store = useServiceStore()

const handleAction = async (item: any) => {
  if (item.confirm) {
    if (!confirm(`Are you sure you want to ${item.label}?`)) return
  }
  
  if (!props.serviceId && !store.selectedServiceId) {
    console.error("No service ID available for action")
    return
  }

  const sid = props.serviceId || store.selectedServiceId
  if (sid) {
    await store.executeAction(sid, item.action_id)
  }
}

const getButtonStyle = (style: string) => {
  switch (style) {
    case 'primary': return 'bg-blue-500/20 text-blue-400 border-blue-500/30 hover:bg-blue-500/30'
    case 'danger': return 'bg-rose-500/20 text-rose-400 border-rose-500/30 hover:bg-rose-500/30'
    case 'success': return 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30 hover:bg-emerald-500/30'
    default: return 'bg-zinc-800 text-zinc-300 border-zinc-700 hover:bg-zinc-700'
  }
}
</script>

<template>
  <div class="flex flex-col gap-1 py-1">
    <span class="text-xs text-zinc-500">{{ widget.label }}</span>
    <div class="flex flex-wrap gap-2">
      <button
        v-for="(item, idx) in widget.items"
        :key="idx"
        class="px-3 py-1.5 rounded text-xs font-medium border transition-colors flex items-center gap-1.5"
        :class="getButtonStyle(item.style)"
        @click.stop="handleAction(item)"
      >
        <!-- Simple logic for icons based on common keywords if desired, or just text -->
        {{ item.label }}
      </button>
    </div>
  </div>
</template>
