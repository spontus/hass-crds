import { describe, it, expect } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactNode } from 'react'
import {
  useEntityTypes,
  useEntitySchema,
  useNamespaces,
  useEntities,
  useEntity,
} from './entities'

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  })
  return ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}

describe('useEntityTypes', () => {
  it('fetches entity types', async () => {
    const { result } = renderHook(() => useEntityTypes(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(result.current.data?.entityTypes).toBeDefined()
    expect(result.current.data?.entityTypes.length).toBeGreaterThan(0)
    expect(result.current.data?.categories).toBeDefined()
  })

  it('includes expected entity kinds', async () => {
    const { result } = renderHook(() => useEntityTypes(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    const kinds = result.current.data?.entityTypes.map((et) => et.kind)
    expect(kinds).toContain('MQTTButton')
    expect(kinds).toContain('MQTTSensor')
  })
})

describe('useEntitySchema', () => {
  it('fetches schema for valid kind', async () => {
    const { result } = renderHook(() => useEntitySchema('MQTTButton'), {
      wrapper: createWrapper(),
    })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(result.current.data?.kind).toBe('MQTTButton')
    expect(result.current.data?.schema).toBeDefined()
    expect(result.current.data?.schema.properties).toBeDefined()
  })

  it('does not fetch when kind is undefined', async () => {
    const { result } = renderHook(() => useEntitySchema(undefined), {
      wrapper: createWrapper(),
    })

    expect(result.current.isFetching).toBe(false)
    expect(result.current.data).toBeUndefined()
  })

  it('returns error for unknown kind', async () => {
    const { result } = renderHook(() => useEntitySchema('UnknownKind'), {
      wrapper: createWrapper(),
    })

    await waitFor(() => expect(result.current.isError).toBe(true))
  })
})

describe('useNamespaces', () => {
  it('fetches namespaces', async () => {
    const { result } = renderHook(() => useNamespaces(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(result.current.data?.namespaces).toBeDefined()
    expect(result.current.data?.namespaces.length).toBeGreaterThan(0)
    expect(result.current.data?.namespaces[0].name).toBeDefined()
    expect(result.current.data?.namespaces[0].status).toBeDefined()
  })
})

describe('useEntities', () => {
  it('fetches all entities without filters', async () => {
    const { result } = renderHook(() => useEntities(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(result.current.data?.items).toBeDefined()
    expect(result.current.data?.total).toBeGreaterThan(0)
  })

  it('filters by kind', async () => {
    const { result } = renderHook(() => useEntities('MQTTButton'), {
      wrapper: createWrapper(),
    })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(result.current.data?.items.every((e) => e.kind === 'MQTTButton')).toBe(true)
  })

  it('filters by namespace', async () => {
    const { result } = renderHook(() => useEntities(undefined, 'default'), {
      wrapper: createWrapper(),
    })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(result.current.data?.items.every((e) => e.namespace === 'default')).toBe(true)
  })

  it('filters by both kind and namespace', async () => {
    const { result } = renderHook(() => useEntities('MQTTButton', 'default'), {
      wrapper: createWrapper(),
    })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(
      result.current.data?.items.every(
        (e) => e.kind === 'MQTTButton' && e.namespace === 'default'
      )
    ).toBe(true)
  })
})

describe('useEntity', () => {
  it('fetches single entity', async () => {
    const { result } = renderHook(
      () => useEntity('MQTTButton', 'default', 'test-button'),
      { wrapper: createWrapper() }
    )

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(result.current.data?.kind).toBe('MQTTButton')
    expect(result.current.data?.metadata.name).toBe('test-button')
    expect(result.current.data?.metadata.namespace).toBe('default')
    expect(result.current.data?.spec).toBeDefined()
  })

  it('does not fetch with missing parameters', async () => {
    const { result } = renderHook(() => useEntity('', 'default', 'test'), {
      wrapper: createWrapper(),
    })

    expect(result.current.isFetching).toBe(false)
    expect(result.current.data).toBeUndefined()
  })
})
