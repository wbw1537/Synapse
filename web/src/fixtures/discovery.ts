import payload from '../../../testdata/discovery-v1.json'
import type { Service } from '../stores/services'

// The transport fixture is also compiled against the UI's flat data model.
if (payload.api_version !== 'v1' || payload.status !== 'online') {
  throw new Error('Discovery fixture must use v1 and online status')
}
export const discoveryFixture: Service = {
  ...payload,
  api_version: payload.api_version,
  status: payload.status,
  last_seen: '',
}
