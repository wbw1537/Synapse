import { defineStore } from 'pinia'
import { ref } from 'vue'
import mqtt from 'mqtt'

export interface Component {
  id: string
  type: string
  label: string
  value: any
  unit?: string
  // Props
  copyable?: boolean
  mapping?: Record<string, any>
  min?: number
  max?: number
  thresholds?: Record<string, string>
  max_items?: number
  items?: any[]
  action_id?: string
  uri?: string
  text?: string
  icon?: string
  animate?: boolean
  style?: string
  confirm?: boolean
}

export interface LayoutSection {
  type: string
  title: string
  children: string[]
}

export interface LayoutSchema {
  type: string
  root: LayoutSection[]
}

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
  
  // Protocol v2
  api_version?: string
  layout?: LayoutSchema
  components?: Record<string, Component>
  
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

  const executeAction = async (serviceId: string, actionId: string) => {
    try {
      const response = await fetch(`/api/v1/services/${serviceId}/actions/${actionId}`, {
        method: 'POST'
      })
      if (!response.ok) {
        throw new Error(await response.text())
      }
      return true
    } catch (err) {
      console.error('Failed to execute action:', err)
      alert(`Failed to execute action: ${err}`)
      return false
    }
  }

  const initMQTT = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
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
        const newData = JSON.parse(payload.toString()) as Service
        
        if (newData.id) {
          const existing = services.value[newData.id]
          
          if (existing && existing.components && newData.components) {
            // Merge Components (specifically for log streams)
            Object.keys(newData.components).forEach(compId => {
              const newComp = newData.components![compId]
              const oldComp = existing.components![compId]

              if (!newComp) return

              if (oldComp && newComp.type === 'log_stream') {
                let logs: any[] = []
                
                if (Array.isArray(oldComp.value)) {
                  logs = [...oldComp.value]
                } else if (typeof oldComp.value === 'string') {
                  logs = [oldComp.value]
                }

                if (typeof newComp.value === 'string') {
                  logs.push(newComp.value)
                } else if (Array.isArray(newComp.value)) {
                    // If the payload sends a list, usually it's replacing, but let's assume append for safety unless explicit? 
                    // Actually, if it sends a list, it's likely the full history or a batch. 
                    // For now, if it's a list, assume it replaces or appends? 
                    // Let's assume replace if list, append if string.
                    // But wait, if backend sends full state on load, it sends a list.
                    // If Axon sends a list, it might be init.
                    // Let's stick to: String = Append, List = Replace.
                    logs = newComp.value
                }

                const maxItems = newComp.max_items || 10
                if (logs.length > maxItems) {
                  logs = logs.slice(logs.length - maxItems)
                }
                
                newComp.value = logs
              }
            })
            
            // Re-assign the merged components back to newData to ensure reactivity update uses the merged data
            // (Though we modified newComp in place, which is inside newData)
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
    executeAction,
    initMQTT
  }
})
