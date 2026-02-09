import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from './client'
import type {
  EntityListResponse,
  EntityTypesResponse,
  NamespaceListResponse,
  EntitySchema,
  KubernetesResource,
} from '../types/entity'

export function useEntityTypes() {
  return useQuery({
    queryKey: ['entity-types'],
    queryFn: () => api.get<EntityTypesResponse>('/entity-types'),
  })
}

export function useEntitySchema(kind: string | undefined) {
  return useQuery({
    queryKey: ['entity-schema', kind],
    queryFn: () => api.get<EntitySchema>(`/entity-types/${kind}/schema`),
    enabled: !!kind,
  })
}

export function useNamespaces() {
  return useQuery({
    queryKey: ['namespaces'],
    queryFn: () => api.get<NamespaceListResponse>('/namespaces'),
  })
}

export function useEntities(kind?: string, namespace?: string) {
  const params = new URLSearchParams()
  if (kind) params.set('kind', kind)
  if (namespace) params.set('namespace', namespace)
  const query = params.toString() ? `?${params.toString()}` : ''

  return useQuery({
    queryKey: ['entities', kind, namespace],
    queryFn: () => api.get<EntityListResponse>(`/entities${query}`),
  })
}

export function useEntity(kind: string, namespace: string, name: string) {
  return useQuery({
    queryKey: ['entity', kind, namespace, name],
    queryFn: () => api.get<KubernetesResource>(`/entities/${kind}/${namespace}/${name}`),
    enabled: !!kind && !!namespace && !!name,
  })
}

export function useCreateEntity() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ kind, namespace, data }: { kind: string; namespace: string; data: unknown }) =>
      api.post<KubernetesResource>(`/entities/${kind}/${namespace}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['entities'] })
    },
  })
}

export function useUpdateEntity() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      kind,
      namespace,
      name,
      data,
    }: {
      kind: string
      namespace: string
      name: string
      data: unknown
    }) => api.put<KubernetesResource>(`/entities/${kind}/${namespace}/${name}`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['entities'] })
      queryClient.invalidateQueries({
        queryKey: ['entity', variables.kind, variables.namespace, variables.name],
      })
    },
  })
}

export function useDeleteEntity() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ kind, namespace, name }: { kind: string; namespace: string; name: string }) =>
      api.delete(`/entities/${kind}/${namespace}/${name}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['entities'] })
    },
  })
}
