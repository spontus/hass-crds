import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  Lightbulb,
  ToggleLeft,
  Thermometer,
  Shield,
  Blinds,
  Bot,
  Camera,
  MapPin,
  Settings,
  ChevronDown,
  ChevronRight,
  Activity,
} from 'lucide-react'
import { useState } from 'react'
import { useEntityTypes } from '../../api/entities'

const categoryIcons: Record<string, typeof LayoutDashboard> = {
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

const categoryOrder = [
  'Controls',
  'Sensors',
  'Lighting',
  'Climate',
  'Security',
  'Covers',
  'Devices',
  'Media',
  'Tracking',
  'Utility',
]

export default function Sidebar() {
  const { data } = useEntityTypes()
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(new Set(['Controls']))

  const toggleCategory = (category: string) => {
    setExpandedCategories((prev) => {
      const next = new Set(prev)
      if (next.has(category)) {
        next.delete(category)
      } else {
        next.add(category)
      }
      return next
    })
  }

  const categories = data?.categories || {}

  return (
    <aside className="w-64 bg-slate-900 border-r border-slate-800 flex flex-col h-full">
      <div className="p-4 border-b border-slate-800">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-lg bg-gradient-to-br from-ha-blue to-ha-cyan flex items-center justify-center">
            <span className="text-white font-bold text-sm">HA</span>
          </div>
          <div>
            <h1 className="font-semibold text-slate-100">HASS CRDs</h1>
            <p className="text-xs text-slate-500">MQTT Entity Manager</p>
          </div>
        </div>
      </div>

      <nav className="flex-1 overflow-y-auto scrollbar-thin py-4">
        <NavLink
          to="/"
          className={({ isActive }) =>
            `flex items-center gap-3 px-4 py-2 mx-2 rounded-lg text-sm transition-colors ${
              isActive
                ? 'bg-slate-800 text-slate-100'
                : 'text-slate-400 hover:text-slate-100 hover:bg-slate-800/50'
            }`
          }
        >
          <LayoutDashboard className="w-4 h-4" />
          Dashboard
        </NavLink>

        <div className="mt-4">
          <div className="px-4 py-2 text-xs font-medium text-slate-500 uppercase tracking-wider">
            Entity Types
          </div>

          {categoryOrder.map((category) => {
            const types = categories[category] || []
            if (types.length === 0) return null

            const Icon = categoryIcons[category] || Settings
            const isExpanded = expandedCategories.has(category)

            return (
              <div key={category}>
                <button
                  onClick={() => toggleCategory(category)}
                  className="w-full flex items-center justify-between px-4 py-2 mx-2 pr-4 text-sm text-slate-400 hover:text-slate-100 transition-colors"
                  style={{ width: 'calc(100% - 16px)' }}
                >
                  <span className="flex items-center gap-3">
                    <Icon className="w-4 h-4" />
                    {category}
                  </span>
                  {isExpanded ? (
                    <ChevronDown className="w-4 h-4" />
                  ) : (
                    <ChevronRight className="w-4 h-4" />
                  )}
                </button>

                {isExpanded && (
                  <div className="ml-6 border-l border-slate-800 pl-3 mt-1 space-y-0.5">
                    {types.map((type) => (
                      <NavLink
                        key={type.kind}
                        to={`/entities/${type.kind}`}
                        className={({ isActive }) =>
                          `block px-3 py-1.5 rounded-md text-sm transition-colors ${
                            isActive
                              ? 'bg-slate-800 text-ha-blue'
                              : 'text-slate-500 hover:text-slate-300 hover:bg-slate-800/50'
                          }`
                        }
                      >
                        {type.kind.replace('MQTT', '')}
                      </NavLink>
                    ))}
                  </div>
                )}
              </div>
            )
          })}
        </div>
      </nav>

      <div className="p-4 border-t border-slate-800">
        <div className="text-xs text-slate-500">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-ha-green animate-pulse" />
            Connected to cluster
          </div>
        </div>
      </div>
    </aside>
  )
}
