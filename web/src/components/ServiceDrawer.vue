<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useServiceStore } from '../stores/services'
import { 
  X, Settings, BookOpen, ScrollText, 
  Activity, Clock, MapPin 
} from 'lucide-vue-next'
import MarkdownIt from 'markdown-it'
import StatWidget from './widgets/StatWidget.vue'
import StatusWidget from './widgets/StatusWidget.vue'
import GaugeWidget from './widgets/GaugeWidget.vue'
import LogStreamWidget from './widgets/LogStreamWidget.vue'
import ActionGroupWidget from './widgets/ActionGroupWidget.vue'
import LinkWidget from './widgets/LinkWidget.vue'

const store = useServiceStore()
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

const activeTab = ref<'runbook' | 'health' | 'logs'>('runbook')

const service = computed(() => {
  if (!store.selectedServiceId) return null
  return store.services[store.selectedServiceId]
})

const isOpen = computed(() => !!store.selectedServiceId)

const close = () => {
  store.selectService(null)
}

const renderedDocs = computed(() => {
  if (!service.value?.markdown_docs) return '<p class="text-zinc-500 italic">No documentation provided.</p>'
  return md.render(service.value.markdown_docs)
})

const widgetMap: Record<string, any> = {
  stat: StatWidget,
  status_indicator: StatusWidget,
  gauge: GaugeWidget,
  log_stream: LogStreamWidget,
  action_group: ActionGroupWidget,
  link: LinkWidget
}

// Separate widgets by type for tabs
const logWidgets = computed(() => {
  return service.value?.widgets?.filter(w => w.type === 'log_stream') || []
})

const healthWidgets = computed(() => {
  return service.value?.widgets?.filter(w => w.type !== 'log_stream') || []
})

// Handle "Critical" actions for the footer
// For now, we define critical as style === 'danger' or confirm === true
const criticalActions = computed(() => {
  return service.value?.actions?.filter(a => a.style === 'danger' || a.require_confirmation) || []
})

// Reset tab when service changes
watch(() => store.selectedServiceId, () => {
  activeTab.value = 'runbook'
})

const handleAction = (actionId: string) => {
  alert(`Triggering action: ${actionId} (Mock)`)
}
</script>

<template>
  <div 
    v-if="isOpen"
    class="fixed inset-0 z-50 flex justify-end"
  >
    <!-- Backdrop -->
    <div 
      class="absolute inset-0 bg-black/60 backdrop-blur-sm transition-opacity"
      @click="close"
    ></div>

    <!-- Panel -->
    <div 
      class="relative w-full max-w-2xl bg-zinc-950 border-l border-zinc-800 shadow-2xl flex flex-col transform transition-transform duration-300 ease-out"
      :class="isOpen ? 'translate-x-0' : 'translate-x-full'"
    >
      <template v-if="service">
        <!-- Header -->
        <header class="p-6 border-b border-zinc-900">
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-zinc-900 border border-zinc-800 text-emerald-400">
                <Activity class="w-6 h-6" />
              </div>
              <div>
                <h2 class="text-2xl font-black text-white uppercase tracking-tight italic">{{ service.name }}</h2>
                <div class="flex items-center gap-3 mt-1 text-xs font-bold uppercase tracking-widest text-zinc-500">
                  <span class="flex items-center gap-1"><MapPin class="w-3 h-3" /> {{ service.group }}</span>
                  <span class="flex items-center gap-1"><Clock class="w-3 h-3" /> Seen 2m ago</span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <button class="p-2 hover:bg-zinc-900 rounded-lg text-zinc-500 transition-colors">
                <Settings class="w-5 h-5" />
              </button>
              <button 
                @click="close"
                class="p-2 hover:bg-zinc-900 rounded-lg text-zinc-500 transition-colors"
              >
                <X class="w-5 h-5" />
              </button>
            </div>
          </div>
        </header>

        <!-- Tabs Navigation -->
        <nav class="flex px-6 border-b border-zinc-900 bg-zinc-950/50 sticky top-0 z-10">
          <button 
            v-for="tab in (['runbook', 'health', 'logs'] as const)" 
            :key="tab"
            @click="activeTab = tab"
            :class="[
              'px-4 py-3 text-xs font-black uppercase tracking-widest border-b-2 transition-all',
              activeTab === tab 
                ? 'border-emerald-500 text-emerald-500 bg-emerald-500/5' 
                : 'border-transparent text-zinc-500 hover:text-zinc-300'
            ]"
          >
            <span class="flex items-center gap-2">
              <BookOpen v-if="tab === 'runbook'" class="w-3.5 h-3.5" />
              <Activity v-if="tab === 'health'" class="w-3.5 h-3.5" />
              <ScrollText v-if="tab === 'logs'" class="w-3.5 h-3.5" />
              {{ tab }}
            </span>
          </button>
        </nav>

        <!-- Content Area -->
        <div class="flex-1 overflow-y-auto p-6">
          <!-- Runbook Tab -->
          <div v-if="activeTab === 'runbook'" class="prose prose-invert prose-zinc max-w-none">
            <div v-html="renderedDocs" class="markdown-body"></div>
          </div>

          <!-- Health Tab -->
          <div v-if="activeTab === 'health'" class="grid grid-cols-1 gap-4">
             <div v-for="widget in healthWidgets" :key="widget.id" class="p-4 rounded-xl bg-zinc-900/50 border border-zinc-800/50">
                <component 
                  :is="widgetMap[widget.type] || StatWidget" 
                  :widget="widget" 
                />
             </div>
          </div>

          <!-- Logs Tab -->
          <div v-if="activeTab === 'logs'" class="flex flex-col gap-6">
            <div v-if="logWidgets.length === 0" class="text-center py-20 border-2 border-dashed border-zinc-900 rounded-2xl">
              <p class="text-zinc-600 text-sm font-bold uppercase tracking-widest">No log streams active</p>
            </div>
            <div v-for="widget in logWidgets" :key="widget.id">
              <LogStreamWidget :widget="widget" autoHeight />
            </div>
          </div>
        </div>

        <!-- Sticky Footer for Actions -->
        <footer v-if="criticalActions.length > 0" class="p-6 border-t border-zinc-900 bg-zinc-950/80 backdrop-blur-md">
          <div class="flex flex-col gap-3">
            <h4 class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-500">Critical Operations</h4>
            <div class="flex gap-3">
              <button 
                v-for="action in criticalActions" 
                :key="action.id"
                @click="handleAction(action.id)"
                :class="[
                  'flex-1 px-4 py-3 rounded-xl text-sm font-bold transition-all border shadow-lg',
                  action.style === 'danger' 
                    ? 'bg-rose-600 border-rose-500 text-white hover:bg-rose-500' 
                    : 'bg-zinc-800 border-zinc-700 text-zinc-100 hover:bg-zinc-700'
                ]"
              >
                {{ action.label }}
              </button>
            </div>
          </div>
        </footer>
      </template>
    </div>
  </div>
</template>

<style scoped>
@reference "tailwindcss";

.markdown-body :deep(h1) { @apply text-2xl font-black uppercase tracking-tight italic text-white mb-4 border-b border-zinc-900 pb-2; }
.markdown-body :deep(h2) { @apply text-xl font-bold text-zinc-200 mt-6 mb-3; }
.markdown-body :deep(p) { @apply text-zinc-400 leading-relaxed mb-4; }
.markdown-body :deep(code) { @apply bg-zinc-900 text-emerald-400 px-1.5 py-0.5 rounded font-mono text-sm; }
.markdown-body :deep(pre) { @apply bg-zinc-900 border border-zinc-800 p-4 rounded-xl overflow-x-auto mb-4; }
.markdown-body :deep(pre code) { @apply bg-transparent p-0 text-zinc-300; }
.markdown-body :deep(ul) { @apply list-disc list-inside text-zinc-400 mb-4; }
.markdown-body :deep(li) { @apply mb-1; }
</style>
