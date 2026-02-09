import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { ArrowLeft, Save, Trash2 } from 'lucide-react'
import { useEntity, useEntitySchema, useUpdateEntity, useDeleteEntity, useEntityTypes } from '../api/entities'
import DynamicForm from '../components/forms/DynamicForm'

export default function EntityEditPage() {
  const { kind, namespace, name } = useParams<{ kind: string; namespace: string; name: string }>()
  const navigate = useNavigate()
  const [formData, setFormData] = useState<Record<string, unknown>>({})
  const [isInitialized, setIsInitialized] = useState(false)

  const { data: entity, isLoading: entityLoading } = useEntity(kind!, namespace!, name!)
  const { data: schema, isLoading: schemaLoading } = useEntitySchema(kind)
  const { data: typesData } = useEntityTypes()
  const updateEntity = useUpdateEntity()
  const deleteEntity = useDeleteEntity()

  const entityType = typesData?.entityTypes?.find((t) => t.kind === kind)

  useEffect(() => {
    if (entity?.spec && !isInitialized) {
      setFormData(entity.spec as Record<string, unknown>)
      setIsInitialized(true)
    }
  }, [entity, isInitialized])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!kind || !namespace || !name) return

    const cleanSpec = cleanEmptyValues(formData)

    await updateEntity.mutateAsync({
      kind,
      namespace,
      name,
      data: {
        metadata: {
          name,
          namespace,
        },
        spec: cleanSpec,
      },
    })

    navigate(`/entities/${kind}`)
  }

  const handleDelete = async () => {
    if (!kind || !namespace || !name) return
    if (confirm(`Delete ${name}?`)) {
      await deleteEntity.mutateAsync({ kind, namespace, name })
      navigate(`/entities/${kind}`)
    }
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

  if (entityLoading || schemaLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-slate-400">Loading...</div>
      </div>
    )
  }

  if (!entity) {
    return (
      <div className="flex flex-col items-center justify-center h-64">
        <div className="text-slate-400 mb-4">Entity not found</div>
        <Link to={`/entities/${kind}`} className="btn btn-secondary">
          Back to list
        </Link>
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-4">
          <Link
            to={`/entities/${kind}`}
            className="p-2 rounded-lg hover:bg-slate-800 transition-colors"
          >
            <ArrowLeft className="w-5 h-5 text-slate-400" />
          </Link>
          <div>
            <h1 className="text-2xl font-semibold text-slate-100">
              Edit {kind?.replace('MQTT', '')}
            </h1>
            <p className="text-slate-400 mt-1">
              {namespace}/{name}
            </p>
          </div>
        </div>
        <button
          onClick={handleDelete}
          disabled={deleteEntity.isPending}
          className="btn btn-ghost text-ha-red hover:bg-ha-red/10"
        >
          <Trash2 className="w-4 h-4" />
          Delete
        </button>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="grid grid-cols-3 gap-6">
          <div className="col-span-2 space-y-6">
            <div className="card p-6">
              <h2 className="text-lg font-medium text-slate-100 mb-4">Resource Info</h2>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label block mb-2">Name</label>
                  <input
                    type="text"
                    value={name}
                    disabled
                    className="input opacity-50"
                  />
                </div>
                <div>
                  <label className="label block mb-2">Namespace</label>
                  <input
                    type="text"
                    value={namespace}
                    disabled
                    className="input opacity-50"
                  />
                </div>
              </div>
              {entityType && (
                <p className="text-sm text-slate-500 mt-4">{entityType.description}</p>
              )}
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
              <h2 className="text-lg font-medium text-slate-100 mb-4">Status</h2>

              <div className="space-y-3 mb-6">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-slate-400">Published</span>
                  <span className="text-slate-100">
                    {(entity.status as Record<string, unknown>)?.lastPublished ? 'Yes' : 'No'}
                  </span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-slate-400">Created</span>
                  <span className="text-slate-100">
                    {entity.metadata.creationTimestamp
                      ? new Date(entity.metadata.creationTimestamp).toLocaleDateString()
                      : '-'}
                  </span>
                </div>
                {Boolean((entity.status as Record<string, unknown>)?.discoveryTopic) && (
                  <div className="text-sm">
                    <span className="text-slate-400 block mb-1">Discovery Topic</span>
                    <code className="text-xs text-slate-300 bg-slate-800 px-2 py-1 rounded block truncate">
                      {String((entity.status as Record<string, unknown>).discoveryTopic)}
                    </code>
                  </div>
                )}
              </div>

              <button
                type="submit"
                disabled={updateEntity.isPending}
                className="btn btn-primary w-full"
              >
                <Save className="w-4 h-4" />
                {updateEntity.isPending ? 'Saving...' : 'Save Changes'}
              </button>

              {updateEntity.isError && (
                <div className="mt-4 p-3 rounded-lg bg-ha-red/10 border border-ha-red/20 text-sm text-ha-red">
                  {updateEntity.error?.message || 'Failed to update entity'}
                </div>
              )}
            </div>
          </div>
        </div>
      </form>
    </div>
  )
}
