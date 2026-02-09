import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { ArrowLeft, Save, Code } from 'lucide-react'
import { useEntitySchema, useNamespaces, useCreateEntity, useEntityTypes } from '../api/entities'
import DynamicForm from '../components/forms/DynamicForm'

export default function EntityCreatePage() {
  const { kind } = useParams<{ kind: string }>()
  const navigate = useNavigate()
  const [namespace, setNamespace] = useState('default')
  const [name, setName] = useState('')
  const [formData, setFormData] = useState<Record<string, unknown>>({})
  const [showYaml, setShowYaml] = useState(false)

  const { data: schema, isLoading: schemaLoading } = useEntitySchema(kind)
  const { data: namespacesData } = useNamespaces()
  const { data: typesData } = useEntityTypes()
  const createEntity = useCreateEntity()

  const entityType = typesData?.entityTypes?.find((t) => t.kind === kind)
  const namespaces = [...(namespacesData?.namespaces || [])].sort((a, b) => a.name.localeCompare(b.name))

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!kind || !name) return

    const cleanSpec = cleanEmptyValues(formData)

    await createEntity.mutateAsync({
      kind,
      namespace,
      data: {
        metadata: { name },
        spec: cleanSpec,
      },
    })

    navigate(`/entities/${kind}`)
  }

  const cleanEmptyValues = (obj: Record<string, unknown>): Record<string, unknown> => {
    const result: Record<string, unknown> = {}
    for (const [key, value] of Object.entries(obj)) {
      if (value === '' || value === null || value === undefined) continue
      if (Array.isArray(value) && value.length === 0) continue
      if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
        const cleaned = cleanEmptyValues(value as Record<string, unknown>)
        if (Object.keys(cleaned).length > 0) {
          result[key] = cleaned
        }
      } else {
        result[key] = value
      }
    }
    return result
  }

  const yamlPreview = `apiVersion: mqtt.home-assistant.io/v1alpha1
kind: ${kind}
metadata:
  name: ${name || '<name>'}
  namespace: ${namespace}
spec:
${formatYamlSpec(cleanEmptyValues(formData), '  ')}`

  if (schemaLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-slate-400">Loading schema...</div>
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <Link
          to={`/entities/${kind}`}
          className="p-2 rounded-lg hover:bg-slate-800 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-slate-400" />
        </Link>
        <div>
          <h1 className="text-2xl font-semibold text-slate-100">
            Create {kind?.replace('MQTT', '')}
          </h1>
          {entityType && (
            <p className="text-slate-400 mt-1">{entityType.description}</p>
          )}
        </div>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="grid grid-cols-3 gap-6">
          <div className="col-span-2 space-y-6">
            <div className="card p-6">
              <h2 className="text-lg font-medium text-slate-100 mb-4">Basic Information</h2>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label block mb-2">Resource Name *</label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, '-'))}
                    placeholder="my-entity"
                    className="input"
                    required
                  />
                  <p className="text-xs text-slate-500 mt-1">
                    Lowercase, alphanumeric, hyphens only
                  </p>
                </div>
                <div>
                  <label className="label block mb-2">Namespace *</label>
                  <select
                    value={namespace}
                    onChange={(e) => setNamespace(e.target.value)}
                    className="input"
                    required
                  >
                    {namespaces.map((ns) => (
                      <option key={ns.name} value={ns.name}>
                        {ns.name}
                      </option>
                    ))}
                  </select>
                </div>
              </div>
            </div>

            {schema && (
              <DynamicForm
                schema={schema.schema}
                value={formData}
                onChange={setFormData}
              />
            )}
          </div>

          <div className="space-y-6">
            <div className="card p-6 sticky top-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-medium text-slate-100">Preview</h2>
                <button
                  type="button"
                  onClick={() => setShowYaml(!showYaml)}
                  className="btn btn-ghost text-xs"
                >
                  <Code className="w-3 h-3" />
                  {showYaml ? 'Hide' : 'Show'} YAML
                </button>
              </div>

              {showYaml && (
                <pre className="bg-slate-950 rounded-lg p-4 text-sm font-mono text-slate-300 overflow-x-auto mb-4 max-h-96 overflow-y-auto scrollbar-thin">
                  {yamlPreview}
                </pre>
              )}

              <button
                type="submit"
                disabled={!name || createEntity.isPending}
                className="btn btn-primary w-full"
              >
                <Save className="w-4 h-4" />
                {createEntity.isPending ? 'Creating...' : 'Create Entity'}
              </button>

              {createEntity.isError && (
                <div className="mt-4 p-3 rounded-lg bg-ha-red/10 border border-ha-red/20 text-sm text-ha-red">
                  {createEntity.error?.message || 'Failed to create entity'}
                </div>
              )}
            </div>
          </div>
        </div>
      </form>
    </div>
  )
}

function formatYamlSpec(obj: Record<string, unknown>, indent: string): string {
  const lines: string[] = []

  for (const [key, value] of Object.entries(obj)) {
    if (value === null || value === undefined) continue

    if (typeof value === 'object' && !Array.isArray(value)) {
      lines.push(`${indent}${key}:`)
      lines.push(formatYamlSpec(value as Record<string, unknown>, indent + '  '))
    } else if (Array.isArray(value)) {
      lines.push(`${indent}${key}:`)
      for (const item of value) {
        if (typeof item === 'object') {
          lines.push(`${indent}- `)
          lines.push(formatYamlSpec(item as Record<string, unknown>, indent + '  '))
        } else {
          lines.push(`${indent}- ${JSON.stringify(item)}`)
        }
      }
    } else if (typeof value === 'string') {
      lines.push(`${indent}${key}: "${value}"`)
    } else {
      lines.push(`${indent}${key}: ${value}`)
    }
  }

  return lines.filter((l) => l.trim()).join('\n')
}
