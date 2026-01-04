import { defineStore } from 'pinia'
import { ref } from 'vue'
import mqtt from 'mqtt'

export interface Service {
  id: string
  name: string
  group: string
  tags: string[]
  icon: string
  url: string
  status: 'online' | 'warning' | 'error' | 'offline'
  message: string
  ttl: number
  description: string
  markdown_docs: string
  actions: any[]
  widgets: any[]
  last_seen: string
}

export const useServiceStore = defineStore('services', () => {
  const services = ref<Record<string, Service>>({})
  const loading = ref(false)
  const connected = ref(false)
  const selectedServiceId = ref<string | null>(null)

  const fetchServices = async () => {
    loading.value = true
    try {
      const response = await fetch('/api/v1/services')
      const data: Service[] = await response.json()
      data.forEach(s => {
        services.value[s.id] = s
      })
    } catch (err) {
      console.error('Failed to fetch services:', err)
    } finally {
      loading.value = false
    }
  }

  const selectService = (id: string | null) => {
    selectedServiceId.value = id
  }

  const initMQTT = () => {
    // Determine WS URL based on current origin
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    // For local dev, if vite is on 5173, we might need to point to 8083 specifically
    // but in production/embedded mode, it might be the same host.
    // For now, let's try to be smart.
    const host = window.location.hostname
    const port = '8083' // Synapse MQTT WS Port
    const wsUrl = `${protocol}//${host}:${port}`

    console.log('Connecting to MQTT:', wsUrl)
    const client = mqtt.connect(wsUrl)

    client.on('connect', () => {
      connected.value = true
      client.subscribe('synapse/v1/discovery/#')
    })

    client.on('message', (_topic, payload) => {
      try {
        const newData = JSON.parse(payload.toString())
        
        if (newData.id) {
          const existing = services.value[newData.id]
          
          if (existing) {
            // Merge Widgets
            if (newData.widgets && existing.widgets) {
              newData.widgets.forEach((newWidget: any) => {
                const oldWidget = existing.widgets.find((w: any) => w.id === newWidget.id)
                
                if (oldWidget && newWidget.type === 'log_stream') {
                  // Merge logs
                  let logs: any[] = []
                  
                  if (Array.isArray(oldWidget.value)) {
                    logs = [...oldWidget.value]
                  } else if (typeof oldWidget.value === 'string') {
                    logs = [oldWidget.value]
                  }

                  if (typeof newWidget.value === 'string') {
                    logs.push(newWidget.value)
                  }

                  // Max Items
                  const maxItems = newWidget.max_items || 10
                  if (logs.length > maxItems) {
                    logs = logs.slice(logs.length - maxItems)
                  }
                  
                  newWidget.value = logs
                }
              })
            }
          }
          
          services.value[newData.id] = newData
        }
      } catch (err) {
        console.error('Failed to parse MQTT message:', err)
      }
    })

    client.on('close', () => {
      connected.value = false
    })
  }

  return {
    services,
    loading,
    connected,
    selectedServiceId,
    fetchServices,
    selectService,
    initMQTT
  }
})
