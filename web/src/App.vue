<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useServiceStore } from './stores/services'
import ServiceCard from './components/ServiceCard.vue'
import ServiceDrawer from './components/ServiceDrawer.vue'
import { Wifi, WifiOff } from 'lucide-vue-next'

const store = useServiceStore()

onMounted(() => {
  store.fetchServices()
  store.initMQTT()
})

const sortedServices = computed(() => {
  return Object.values(store.services).sort((a, b) => {
    // Sort by group then by name
    if (a.group !== b.group) return a.group.localeCompare(b.group)
    return a.name.localeCompare(b.name)
  })
})
</script>

<template>
  <div class="min-h-screen bg-zinc-950 text-zinc-100 p-6 md:p-10 font-sans selection:bg-emerald-500/30">
    <header class="max-w-7xl mx-auto mb-10 flex items-center justify-between">
      <div>
        <div class="flex items-center gap-3">
          <h1 class="text-4xl font-black tracking-tight text-white uppercase italic">Synapse</h1>
          <div :class="['flex items-center gap-1.5 px-2 py-1 rounded-full text-[10px] font-bold uppercase tracking-widest', store.connected ? 'bg-emerald-500/10 text-emerald-500' : 'bg-rose-500/10 text-rose-500']">
            <component :is="store.connected ? Wifi : WifiOff" class="w-3 h-3" />
            {{ store.connected ? 'Live' : 'Offline' }}
          </div>
        </div>
        <p class="text-zinc-500 mt-1 font-medium">The nervous system of your homelab.</p>
      </div>

      <div class="hidden md:flex items-center gap-4 text-xs font-bold text-zinc-600 uppercase tracking-widest">
        <span>{{ sortedServices.length }} Services</span>
      </div>
    </header>
    
    <main class="max-w-7xl mx-auto">
      <div v-if="store.loading && sortedServices.length === 0" class="flex items-center justify-center py-20">
        <div class="w-8 h-8 border-2 border-zinc-800 border-t-emerald-500 rounded-full animate-spin"></div>
      </div>

      <div v-else-if="sortedServices.length === 0" class="flex flex-col items-center justify-center py-20 border-2 border-dashed border-zinc-900 rounded-3xl">
        <p class="text-zinc-500 font-medium">No services found.</p>
        <p class="text-zinc-700 text-sm mt-1">Start an agent to see them appear here.</p>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        <ServiceCard v-for="svc in sortedServices" :key="svc.id" :service="svc" />
      </div>
    </main>

    <footer class="max-w-7xl mx-auto mt-20 pt-8 border-t border-zinc-900 text-[10px] text-zinc-700 font-bold uppercase tracking-[0.2em] flex justify-between">
      <span>Synapse v1.0</span>
      <span>Experimental EDA Dashboard</span>
    </footer>

    <ServiceDrawer />
  </div>
</template>

<style>
@import "tailwindcss";

body {
  margin: 0;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* Custom scrollbar for dark mode */
::-webkit-scrollbar {
  width: 8px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: #18181b;
  border-radius: 10px;
}
::-webkit-scrollbar-thumb:hover {
  background: #27272a;
}
</style>
