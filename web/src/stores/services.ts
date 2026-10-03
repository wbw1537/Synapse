import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Component {
  id: string
  type: string
  label: string
  value?: any
  unit?: string
  // Props
  copyable?: boolean
  mapping?: Record<string, any>
  min?: number
  max?: number
  thresholds?: Record<string, string>
  max_items?: number
  items?: ActionGroupItem[]
  action_id?: string
  uri?: string
  text?: string
  icon?: string
  animate?: boolean
  style?: string
  confirm?: boolean
}

export interface ActionGroupItem {
  action_id: string
  label: string
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
  
  // v1 wire contract; optional layout fields tolerate pre-alignment DB rows.
  api_version?: 'v1'
  layout?: LayoutSchema
  components?: Record<string, Component>
  
  last_seen: string
}

export const useServiceStore = defineStore('services', () => {
  const services = ref<Record<string, Service>>({})
  const loading = ref(false)
  const authenticated = ref(false)
  const connected = ref(false)
  const selectedServiceId = ref<string | null>(null)

  const fetchServices = async () => {
    loading.value = true
    try {
      const response = await fetch('/api/v1/services')
      const data: Service[] = await response.json()
      applySnapshot(data)
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

  let events: EventSource | null = null
  const applySnapshot = (data: Service[]) => {
    services.value = Object.fromEntries(data.map(service => [service.id, service]))
  }

  const checkSession = async () => {
    authenticated.value = (await fetch('/api/v1/session')).ok
    if (authenticated.value) initEvents()
  }

  const login = async (token: string) => {
    const response = await fetch('/api/v1/session', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token })
    })
    if (!response.ok) throw new Error('Invalid operator token')
    authenticated.value = true
    initEvents()
  }

  const logout = async () => {
    await fetch('/api/v1/session', { method: 'DELETE' })
    stopEvents()
    authenticated.value = false
    services.value = {}
  }

  const initEvents = () => {
    events?.close()
    loading.value = true
    const stream = new EventSource('/api/v1/events')
    events = stream
    stream.addEventListener('services', event => {
      if (events !== stream) return
      try {
        applySnapshot(JSON.parse((event as MessageEvent).data))
        connected.value = true
        loading.value = false
      } catch (err) {
        connected.value = false
        console.error('Failed to parse service snapshot:', err)
      }
    })
    stream.onerror = () => {
      if (events !== stream) return
      connected.value = false
      loading.value = false
      void fetch('/api/v1/session').then(response => {
        if (events === stream && response.status === 401) {
          stopEvents()
          authenticated.value = false
          services.value = {}
        }
      }).catch(() => {})
    }
  }

  const stopEvents = () => {
    events?.close()
    events = null
    connected.value = false
  }

  return {
    authenticated,
    checkSession,
    login,
    logout,
    services,
    loading,
    connected,
    selectedServiceId,
    fetchServices,
    selectService,
    executeAction,
    initEvents,
    stopEvents,
    applySnapshot
  }
})
