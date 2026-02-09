import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { Plus, Search, Filter, MoreVertical, Trash2, Edit, Check, Clock } from 'lucide-react'
import { useEntities, useEntityTypes, useNamespaces, useDeleteEntity } from '../api/entities'

export default function EntityListPage() {
  const { kind } = useParams<{ kind?: string }>()
  const navigate = useNavigate()
  const [namespaceFilter, setNamespaceFilter] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [menuOpen, setMenuOpen] = useState<string | null>(null)

  const { data: entitiesData, isLoading } = useEntities(kind, namespaceFilter || undefined)
  const { data: typesData } = useEntityTypes()
  const { data: namespacesData } = useNamespaces()
  const deleteEntity = useDeleteEntity()

  const entityType = typesData?.entityTypes?.find((t) => t.kind === kind)
  const entities = entitiesData?.items || []
  const namespaces = [...(namespacesData?.namespaces || [])].sort((a, b) => a.name.localeCompare(b.name))

  const filteredEntities = entities.filter((entity) => {
    if (!searchQuery) return true
    const query = searchQuery.toLowerCase()
    return (
      entity.name.toLowerCase().includes(query) ||
      entity.displayName?.toLowerCase().includes(query) ||
      entity.namespace.toLowerCase().includes(query)
    )
  })

  const handleDelete = async (entityKind: string, namespace: string, name: string) => {
    if (confirm(`Delete ${name}?`)) {
      await deleteEntity.mutateAsync({ kind: entityKind, namespace, name })
      setMenuOpen(null)
    }
  }

  const title = kind ? `${entityType?.description || kind.replace('MQTT', '')}` : 'All Entities'

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-semibold text-slate-100">{title}</h1>
          {kind && entityType && (
            <p className="text-slate-400 mt-1">{entityType.description}</p>
          )}
        </div>
        {kind && (
          <Link to={`/create/${kind}`} className="btn btn-primary">
            <Plus className="w-4 h-4" />
            Create {kind.replace('MQTT', '')}
          </Link>
        )}
      </div>

      <div className="card mb-6">
        <div className="p-4 flex items-center gap-4">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
            <input
              type="text"
              placeholder="Search entities..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input pl-9"
            />
          </div>
          <div className="flex items-center gap-2">
            <Filter className="w-4 h-4 text-slate-500" />
            <select
              value={namespaceFilter}
              onChange={(e) => setNamespaceFilter(e.target.value)}
              className="input w-48"
            >
              <option value="">All namespaces</option>
              {namespaces.map((ns) => (
                <option key={ns.name} value={ns.name}>
                  {ns.name}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-slate-800">
              <th className="text-left p-4 text-sm font-medium text-slate-400">Name</th>
              <th className="text-left p-4 text-sm font-medium text-slate-400">Namespace</th>
              {!kind && (
                <th className="text-left p-4 text-sm font-medium text-slate-400">Kind</th>
              )}
              <th className="text-left p-4 text-sm font-medium text-slate-400">Status</th>
              <th className="text-left p-4 text-sm font-medium text-slate-400">Created</th>
              <th className="w-12"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800">
            {isLoading ? (
              <tr>
                <td colSpan={kind ? 5 : 6} className="p-8 text-center text-slate-500">
                  Loading...
                </td>
              </tr>
            ) : filteredEntities.length === 0 ? (
              <tr>
                <td colSpan={kind ? 5 : 6} className="p-8 text-center text-slate-500">
                  {searchQuery
                    ? 'No matching entities found'
                    : kind
                      ? `No ${kind.replace('MQTT', '')} entities yet`
                      : 'No entities yet'}
                </td>
              </tr>
            ) : (
              filteredEntities.map((entity) => {
                const key = `${entity.kind}/${entity.namespace}/${entity.name}`
                return (
                  <tr key={key} className="hover:bg-slate-800/30 transition-colors">
                    <td className="p-4">
                      <div className="font-medium text-slate-100">
                        {entity.displayName || entity.name}
                      </div>
                      {entity.displayName && (
                        <div className="text-sm text-slate-500">{entity.name}</div>
                      )}
                    </td>
                    <td className="p-4 text-slate-400">{entity.namespace}</td>
                    {!kind && (
                      <td className="p-4">
                        <span className="text-xs px-2 py-1 rounded bg-slate-800 text-slate-300">
                          {entity.kind.replace('MQTT', '')}
                        </span>
                      </td>
                    )}
                    <td className="p-4">
                      <div className="flex items-center gap-2">
                        {entity.published ? (
                          <>
                            <Check className="w-4 h-4 text-ha-green" />
                            <span className="text-sm text-ha-green">Published</span>
                          </>
                        ) : (
                          <>
                            <Clock className="w-4 h-4 text-ha-yellow" />
                            <span className="text-sm text-ha-yellow">Pending</span>
                          </>
                        )}
                      </div>
                    </td>
                    <td className="p-4 text-slate-500 text-sm">
                      {new Date(entity.createdAt).toLocaleDateString()}
                    </td>
                    <td className="p-4 relative">
                      <button
                        onClick={() => setMenuOpen(menuOpen === key ? null : key)}
                        className="p-1 rounded hover:bg-slate-700 transition-colors"
                      >
                        <MoreVertical className="w-4 h-4 text-slate-400" />
                      </button>
                      {menuOpen === key && (
                        <>
                          <div
                            className="fixed inset-0 z-10"
                            onClick={() => setMenuOpen(null)}
                          />
                          <div className="absolute right-4 top-12 z-20 w-40 bg-slate-800 border border-slate-700 rounded-lg shadow-xl py-1">
                            <button
                              onClick={() => {
                                setMenuOpen(null)
                                navigate(`/edit/${entity.kind}/${entity.namespace}/${entity.name}`)
                              }}
                              className="w-full flex items-center gap-2 px-3 py-2 text-sm text-slate-300 hover:bg-slate-700 transition-colors"
                            >
                              <Edit className="w-4 h-4" />
                              Edit
                            </button>
                            <button
                              onClick={() =>
                                handleDelete(entity.kind, entity.namespace, entity.name)
                              }
                              className="w-full flex items-center gap-2 px-3 py-2 text-sm text-ha-red hover:bg-slate-700 transition-colors"
                            >
                              <Trash2 className="w-4 h-4" />
                              Delete
                            </button>
                          </div>
                        </>
                      )}
                    </td>
                  </tr>
                )
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
