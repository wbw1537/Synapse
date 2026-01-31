<script setup lang="ts">
import { computed } from 'vue'
import { Activity, AlertTriangle, XCircle, Info, ExternalLink } from 'lucide-vue-next'
import { useServiceStore } from '../stores/services'
import type { Service } from '../stores/services'
import StatWidget from './widgets/StatWidget.vue'
import StatusWidget from './widgets/StatusWidget.vue'
import GaugeWidget from './widgets/GaugeWidget.vue'
import LogStreamWidget from './widgets/LogStreamWidget.vue'
import ActionGroupWidget from './widgets/ActionGroupWidget.vue'
import LinkWidget from './widgets/LinkWidget.vue'

const props = defineProps<{
  service: Service
}>()

const select = () => {
  useServiceStore().selectService(props.service.id)
}

const widgetMap: Record<string, any> = {
  stat: StatWidget,
  status_indicator: StatusWidget,
  gauge: GaugeWidget,
  log_stream: LogStreamWidget,
  action_group: ActionGroupWidget,
  link: LinkWidget
}

const statusColor = computed(() => {
  switch (props.service.status) {
    case 'online': return 'text-emerald-400 bg-emerald-400/10 border-emerald-400/20'
    case 'warning': return 'text-amber-400 bg-amber-400/10 border-amber-400/20'
    case 'error': return 'text-rose-400 bg-rose-400/10 border-rose-400/20'
    default: return 'text-zinc-500 bg-zinc-500/10 border-zinc-500/20'
  }
})

const statusIcon = computed(() => {
  switch (props.service.status) {
    case 'online': return Activity
    case 'warning': return AlertTriangle
    case 'error': return XCircle
    default: return Info
  }
})

// Helper to get component definition safely
const getComponent = (id: string) => {
  return props.service.components?.[id]
}
</script>

<template>
  <div 
    @click="select"
    class="group p-5 rounded-xl bg-zinc-900 border border-zinc-800 hover:border-zinc-700 transition-all shadow-sm flex flex-col gap-4 cursor-pointer"
  >
    <!-- Header -->
    <div class="flex items-start justify-between">
      <div class="flex items-center gap-3">
        <div :class="['p-2 rounded-lg border', statusColor]">
          <component :is="statusIcon" class="w-5 h-5" />
        </div>
        <div>
          <h3 class="font-semibold text-zinc-100">{{ service.name }}</h3>
          <p class="text-xs text-zinc-500 uppercase tracking-wider font-medium">{{ service.group }}</p>
        </div>
      </div>
      
      <a 
        v-if="service.url" 
        :href="service.url" 
        target="_blank" 
        @click.stop
        class="p-1.5 text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800 rounded-md transition-colors"
      >
        <ExternalLink class="w-4 h-4" />
      </a>
    </div>

    <div v-if="service.message" class="text-sm text-zinc-400 line-clamp-2">
      {{ service.message }}
    </div>

    <!-- Layout Rendering -->
    <div v-if="service.layout && service.layout.root && service.components" class="flex flex-col gap-4 pt-2 border-t border-zinc-800/50">
      <div v-for="(section, idx) in service.layout.root" :key="idx" class="flex flex-col gap-2">
        
        <!-- Section Title (Optional) -->
        <h4 v-if="section.title" class="text-[10px] font-bold text-zinc-600 uppercase tracking-wider">{{ section.title }}</h4>
        
        <!-- Components in Section -->
        <template v-for="compId in section.children" :key="compId">
          <component 
            v-if="getComponent(compId)"
            :is="widgetMap[getComponent(compId)!.type] || StatWidget" 
            :widget="getComponent(compId)"
            :service-id="service.id"
          />
        </template>
      </div>
    </div>

    <!-- Footer Tags -->
    <div class="mt-auto pt-4 flex items-center gap-2">
      <span v-for="tag in service.tags" :key="tag" class="px-2 py-0.5 rounded text-[10px] font-bold bg-zinc-800 text-zinc-500 uppercase tracking-tighter">
        {{ tag }}
      </span>
    </div>
  </div>
</template>