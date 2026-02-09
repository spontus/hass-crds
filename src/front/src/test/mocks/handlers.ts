import { http, HttpResponse } from 'msw'
import {
  mockEntityTypes,
  mockCategories,
  mockEntities,
  mockNamespaces,
  mockButtonSchema,
  mockKubernetesButton,
} from './data'

export const handlers = [
  http.get('/api/v1/entity-types', () => {
    return HttpResponse.json({
      entityTypes: mockEntityTypes,
      categories: mockCategories,
    })
  }),

  http.get('/api/v1/entity-types/:kind/schema', ({ params }) => {
    const { kind } = params
    const entityType = mockEntityTypes.find((et) => et.kind === kind)

    if (!entityType) {
      return HttpResponse.json(
        { error: `unknown entity type: ${kind}` },
        { status: 404 }
      )
    }

    return HttpResponse.json({
      kind,
      apiVersion: 'mqtt.home-assistant.io/v1alpha1',
      description: entityType.description,
      schema: mockButtonSchema,
    })
  }),

  http.get('/api/v1/namespaces', () => {
    return HttpResponse.json({
      namespaces: mockNamespaces,
      total: mockNamespaces.length,
    })
  }),

  http.get('/api/v1/entities', ({ request }) => {
    const url = new URL(request.url)
    const kind = url.searchParams.get('kind')
    const namespace = url.searchParams.get('namespace')

    let filtered = [...mockEntities]

    if (kind) {
      filtered = filtered.filter((e) => e.kind === kind)
    }
    if (namespace) {
      filtered = filtered.filter((e) => e.namespace === namespace)
    }

    return HttpResponse.json({
      items: filtered,
      total: filtered.length,
    })
  }),

  http.get('/api/v1/entities/:kind/:namespace/:name', ({ params }) => {
    const { kind, namespace, name } = params

    if (kind === 'MQTTButton' && namespace === 'default' && name === 'test-button') {
      return HttpResponse.json(mockKubernetesButton)
    }

    return HttpResponse.json(
      { error: 'entity not found' },
      { status: 404 }
    )
  }),

  http.post('/api/v1/entities/:kind/:namespace', async ({ params, request }) => {
    const { kind, namespace } = params
    const body = await request.json() as Record<string, unknown>
    const metadata = body.metadata as Record<string, unknown>

    return HttpResponse.json(
      {
        apiVersion: 'mqtt.home-assistant.io/v1alpha1',
        kind,
        metadata: {
          name: metadata.name,
          namespace,
          resourceVersion: '1',
          creationTimestamp: new Date().toISOString(),
        },
        spec: body.spec,
      },
      { status: 201 }
    )
  }),

  http.put('/api/v1/entities/:kind/:namespace/:name', async ({ params, request }) => {
    const { kind, namespace, name } = params
    const body = await request.json() as Record<string, unknown>

    return HttpResponse.json({
      apiVersion: 'mqtt.home-assistant.io/v1alpha1',
      kind,
      metadata: {
        name,
        namespace,
        resourceVersion: '2',
        creationTimestamp: '2024-01-15T10:00:00Z',
      },
      spec: body.spec,
    })
  }),

  http.delete('/api/v1/entities/:kind/:namespace/:name', () => {
    return new HttpResponse(null, { status: 204 })
  }),
]
