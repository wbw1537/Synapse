import assert from 'node:assert/strict'
import test from 'node:test'
import { createServer } from 'vite'
import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import { renderToString } from '@vue/server-renderer'

test('v1 discovery renders all declared capabilities with flat properties', async () => {
  const vite = await createServer({ server: { middlewareMode: true }, appType: 'custom' })
  try {
    const { discoveryFixture: service } = await vite.ssrLoadModule('/src/fixtures/discovery.ts')
    const { default: ServiceCard } = await vite.ssrLoadModule('/src/components/ServiceCard.vue')
    assert.equal(service.api_version, 'v1')
    assert.equal(service.status, 'online')
    assert.equal(service.components.controls.items[0].action_id, 'restart')
    assert.equal(service.components.cpu.max, 100)
    const app = createSSRApp(ServiceCard, { service })
    app.use(createPinia())
    const html = await renderToString(app)
    for (const text of ['Protocol Demo', 'Uptime', '42', 'CPU', '25%', 'Healthy', 'Application Logs', 'Ready', 'Restart', 'Open docs']) {
      assert.ok(html.includes(text), `missing rendered capability: ${text}`)
    }
    assert.ok(html.includes('width:25%'))
    assert.ok(html.includes('href="https://example.com/docs"'))
    for (const section of service.layout.root) {
      for (const id of section.children) assert.equal(service.components[id].id, id)
    }
  } finally {
    await vite.close()
  }
})
