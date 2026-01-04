<script setup lang="ts">
import { computed } from 'vue'
import { 
  CheckCircle2, AlertTriangle, XCircle, Info, Wifi, WifiOff, Gavel, Activity, Server, Shield
} from 'lucide-vue-next'

const props = defineProps<{
  widget: any
}>()

const iconMap: Record<string, any> = {
  'check': CheckCircle2,
  'alert': AlertTriangle,
  'x': XCircle,
  'info': Info,
  'wifi': Wifi,
  'wifi-off': WifiOff,
  'gavel': Gavel,
  'activity': Activity,
  'server': Server,
  'shield': Shield
}

const currentState = computed(() => {
  const val = String(props.widget.value)
  return props.widget.mapping?.[val] || { text: val, color: 'neutral' }
})

const colorClass = computed(() => {
  switch (currentState.value.color) {
    case 'success': return 'text-emerald-400'
    case 'warning': return 'text-amber-400'
    case 'danger': return 'text-rose-400'
    case 'primary': return 'text-blue-400'
    default: return 'text-zinc-400'
  }
})

const iconComponent = computed(() => {
  if (!currentState.value.icon) return null
  // clean up icon name (e.g., mdi-check -> check) if needed, but assuming simple names
  return iconMap[currentState.value.icon] || Info
})
</script>

<template>
  <div class="flex items-center justify-between text-xs py-1">
    <span class="text-zinc-500">{{ widget.label }}</span>
    <div class="flex items-center gap-1.5" :class="colorClass">
      <component 
        :is="iconComponent" 
        v-if="iconComponent" 
        class="w-3.5 h-3.5" 
        :class="{'animate-pulse': currentState.animate}"
      />
      <span class="font-medium">{{ currentState.text }}</span>
    </div>
  </div>
</template>
