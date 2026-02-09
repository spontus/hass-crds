import { Link } from 'react-router-dom'
import {
  Lightbulb,
  ToggleLeft,
  Thermometer,
  Shield,
  Blinds,
  Bot,
  Camera,
  MapPin,
  Settings,
  Activity,
  Plus,
  ArrowRight,
} from 'lucide-react'
import { useEntities, useEntityTypes } from '../api/entities'

const categoryIcons: Record<string, typeof ToggleLeft> = {
  Controls: ToggleLeft,
  Sensors: Activity,
  Lighting: Lightbulb,
  Climate: Thermometer,
  Security: Shield,
  Covers: Blinds,
  Devices: Bot,
  Media: Camera,
  Tracking: MapPin,
  Utility: Settings,
}

const categoryColors: Record<string, string> = {
  Controls: 'from-blue-500 to-blue-600',
  Sensors: 'from-green-500 to-green-600',
  Lighting: 'from-yellow-500 to-yellow-600',
  Climate: 'from-orange-500 to-orange-600',
  Security: 'from-red-500 to-red-600',
  Covers: 'from-purple-500 to-purple-600',
  Devices: 'from-indigo-500 to-indigo-600',
  Media: 'from-pink-500 to-pink-600',
  Tracking: 'from-cyan-500 to-cyan-600',
  Utility: 'from-slate-500 to-slate-600',
}

export default function Dashboard() {
  const { data: entitiesData, isLoading: entitiesLoading } = useEntities()
  const { data: typesData, isLoading: typesLoading } = useEntityTypes()

  const isLoading = entitiesLoading || typesLoading

  const entities = entitiesData?.items || []
  const categories = typesData?.categories || {}

  const entityCountByKind = entities.reduce(
    (acc, entity) => {
      acc[entity.kind] = (acc[entity.kind] || 0) + 1
      return acc
    },
    {} as Record<string, number>
  )

  const categoryStats = Object.entries(categories).map(([category, types]) => ({
    category,
    count: types.reduce((sum, type) => sum + (entityCountByKind[type.kind] || 0), 0),
    types: types.length,
  }))

  const recentEntities = [...entities]
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    .slice(0, 5)

  const publishedCount = entities.filter((e) => e.published).length
  const unpublishedCount = entities.length - publishedCount

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-semibold text-slate-100">Dashboard</h1>
          <p className="text-slate-400 mt-1">Manage your MQTT entities across the cluster</p>
        </div>
        <Link to="/entities" className="btn btn-primary">
          <Plus className="w-4 h-4" />
          Create Entity
        </Link>
      </div>

      <div className="grid grid-cols-4 gap-4 mb-8">
        <div className="card p-5">
          <div className="text-sm text-slate-400 mb-1">Total Entities</div>
          <div className="text-3xl font-semibold text-slate-100">
            {isLoading ? '-' : entities.length}
          </div>
        </div>
        <div className="card p-5">
          <div className="text-sm text-slate-400 mb-1">Published</div>
          <div className="text-3xl font-semibold text-ha-green">
            {isLoading ? '-' : publishedCount}
          </div>
        </div>
        <div className="card p-5">
          <div className="text-sm text-slate-400 mb-1">Pending</div>
          <div className="text-3xl font-semibold text-ha-yellow">
            {isLoading ? '-' : unpublishedCount}
          </div>
        </div>
        <div className="card p-5">
          <div className="text-sm text-slate-400 mb-1">Entity Types</div>
          <div className="text-3xl font-semibold text-slate-100">
            {isLoading ? '-' : typesData?.entityTypes?.length || 0}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="card">
          <div className="p-4 border-b border-slate-800">
            <h2 className="font-medium text-slate-100">Entities by Category</h2>
          </div>
          <div className="p-4 space-y-3">
            {categoryStats.map(({ category, count }) => {
              const Icon = categoryIcons[category] || Settings
              const colorClass = categoryColors[category] || 'from-slate-500 to-slate-600'

              return (
                <Link
                  key={category}
                  to={`/entities?category=${category}`}
                  className="flex items-center gap-3 p-3 rounded-lg hover:bg-slate-800/50 transition-colors group"
                >
                  <div
                    className={`w-10 h-10 rounded-lg bg-gradient-to-br ${colorClass} flex items-center justify-center`}
                  >
                    <Icon className="w-5 h-5 text-white" />
                  </div>
                  <div className="flex-1">
                    <div className="font-medium text-slate-100">{category}</div>
                    <div className="text-sm text-slate-500">{count} entities</div>
                  </div>
                  <ArrowRight className="w-4 h-4 text-slate-600 group-hover:text-slate-400 transition-colors" />
                </Link>
              )
            })}
          </div>
        </div>

        <div className="card">
          <div className="p-4 border-b border-slate-800 flex items-center justify-between">
            <h2 className="font-medium text-slate-100">Recent Entities</h2>
            <Link to="/entities" className="text-sm text-ha-blue hover:text-ha-cyan transition-colors">
              View all
            </Link>
          </div>
          <div className="divide-y divide-slate-800">
            {isLoading ? (
              <div className="p-8 text-center text-slate-500">Loading...</div>
            ) : recentEntities.length === 0 ? (
              <div className="p-8 text-center text-slate-500">No entities yet</div>
            ) : (
              recentEntities.map((entity) => (
                <Link
                  key={`${entity.namespace}/${entity.name}`}
                  to={`/edit/${entity.kind}/${entity.namespace}/${entity.name}`}
                  className="flex items-center gap-3 p-4 hover:bg-slate-800/50 transition-colors"
                >
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-medium text-slate-100 truncate">
                        {entity.displayName || entity.name}
                      </span>
                      <span className="text-xs px-2 py-0.5 rounded bg-slate-800 text-slate-400">
                        {entity.kind.replace('MQTT', '')}
                      </span>
                    </div>
                    <div className="text-sm text-slate-500 truncate">
                      {entity.namespace}/{entity.name}
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <div
                      className={`w-2 h-2 rounded-full ${entity.published ? 'bg-ha-green' : 'bg-ha-yellow'}`}
                    />
                    <span className="text-xs text-slate-500">
                      {entity.published ? 'Published' : 'Pending'}
                    </span>
                  </div>
                </Link>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
