export interface EntityType {
  kind: string
  plural: string
  description: string
  category: string
}

export interface EntitySummary {
  kind: string
  apiVersion: string
  name: string
  namespace: string
  displayName?: string
  published: boolean
  createdAt: string
  labels?: Record<string, string>
}

export interface EntityListResponse {
  items: EntitySummary[]
  total: number
}

export interface EntityTypesResponse {
  entityTypes: EntityType[]
  categories: Record<string, EntityType[]>
}

export interface NamespaceSummary {
  name: string
  labels?: Record<string, string>
  status: string
}

export interface NamespaceListResponse {
  namespaces: NamespaceSummary[]
  total: number
}

export interface SchemaProperty {
  type: string
  description?: string
  default?: unknown
  enum?: string[]
  minimum?: number
  maximum?: number
  pattern?: string
  format?: string
  required?: string[]
  properties?: Record<string, SchemaProperty>
  items?: SchemaProperty
}

export interface EntitySchema {
  kind: string
  apiVersion: string
  description: string
  schema: SchemaProperty
}

export interface KubernetesResource {
  apiVersion: string
  kind: string
  metadata: {
    name: string
    namespace: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
    resourceVersion?: string
    creationTimestamp?: string
  }
  spec: Record<string, unknown>
  status?: Record<string, unknown>
}
