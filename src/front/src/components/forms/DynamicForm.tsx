import { useState } from 'react'
import { ChevronDown, ChevronRight, Plus, Trash2, HelpCircle } from 'lucide-react'
import type { SchemaProperty } from '../../types/entity'

interface DynamicFormProps {
  schema: SchemaProperty
  value: Record<string, unknown>
  onChange: (value: Record<string, unknown>) => void
}

const commonFields = ['name', 'uniqueId', 'icon', 'entityCategory', 'enabledByDefault', 'objectId']
const deviceFields = ['device', 'deviceRef']
const availabilityFields = ['availability', 'availabilityMode']
const mqttFields = ['qos', 'retain', 'encoding', 'jsonAttributesTopic', 'jsonAttributesTemplate', 'rediscoverInterval']

function categorizeFields(properties: Record<string, SchemaProperty>, required: string[] = []) {
  const categorized = {
    required: [] as string[],
    common: [] as string[],
    entitySpecific: [] as string[],
    device: [] as string[],
    availability: [] as string[],
    advanced: [] as string[],
  }

  for (const fieldName of Object.keys(properties)) {
    if (required.includes(fieldName)) {
      categorized.required.push(fieldName)
    } else if (commonFields.includes(fieldName)) {
      categorized.common.push(fieldName)
    } else if (deviceFields.includes(fieldName)) {
      categorized.device.push(fieldName)
    } else if (availabilityFields.includes(fieldName)) {
      categorized.availability.push(fieldName)
    } else if (mqttFields.includes(fieldName)) {
      categorized.advanced.push(fieldName)
    } else {
      categorized.entitySpecific.push(fieldName)
    }
  }

  return categorized
}

export default function DynamicForm({ schema, value, onChange }: DynamicFormProps) {
  const properties = schema.properties || {}
  const required = schema.required || []
  const categorized = categorizeFields(properties, required)

  const updateField = (fieldName: string, fieldValue: unknown) => {
    onChange({ ...value, [fieldName]: fieldValue })
  }

  return (
    <div className="space-y-6">
      {(categorized.required.length > 0 || categorized.entitySpecific.length > 0) && (
        <FormSection title="Entity Configuration" defaultOpen>
          <div className="space-y-4">
            {[...categorized.required, ...categorized.entitySpecific].map((fieldName) => (
              <FieldRenderer
                key={fieldName}
                name={fieldName}
                schema={properties[fieldName]}
                value={value[fieldName]}
                onChange={(v) => updateField(fieldName, v)}
                required={required.includes(fieldName)}
              />
            ))}
          </div>
        </FormSection>
      )}

      {categorized.common.length > 0 && (
        <FormSection title="Common Settings" defaultOpen>
          <div className="space-y-4">
            {categorized.common.map((fieldName) => (
              <FieldRenderer
                key={fieldName}
                name={fieldName}
                schema={properties[fieldName]}
                value={value[fieldName]}
                onChange={(v) => updateField(fieldName, v)}
                required={required.includes(fieldName)}
              />
            ))}
          </div>
        </FormSection>
      )}

      {categorized.device.length > 0 && (
        <FormSection title="Device Configuration">
          <div className="space-y-4">
            {categorized.device.map((fieldName) => (
              <FieldRenderer
                key={fieldName}
                name={fieldName}
                schema={properties[fieldName]}
                value={value[fieldName]}
                onChange={(v) => updateField(fieldName, v)}
                required={required.includes(fieldName)}
              />
            ))}
          </div>
        </FormSection>
      )}

      {categorized.availability.length > 0 && (
        <FormSection title="Availability">
          <div className="space-y-4">
            {categorized.availability.map((fieldName) => (
              <FieldRenderer
                key={fieldName}
                name={fieldName}
                schema={properties[fieldName]}
                value={value[fieldName]}
                onChange={(v) => updateField(fieldName, v)}
                required={required.includes(fieldName)}
              />
            ))}
          </div>
        </FormSection>
      )}

      {categorized.advanced.length > 0 && (
        <FormSection title="Advanced MQTT Settings">
          <div className="space-y-4">
            {categorized.advanced.map((fieldName) => (
              <FieldRenderer
                key={fieldName}
                name={fieldName}
                schema={properties[fieldName]}
                value={value[fieldName]}
                onChange={(v) => updateField(fieldName, v)}
                required={required.includes(fieldName)}
              />
            ))}
          </div>
        </FormSection>
      )}
    </div>
  )
}

interface FormSectionProps {
  title: string
  defaultOpen?: boolean
  children: React.ReactNode
}

function FormSection({ title, defaultOpen = false, children }: FormSectionProps) {
  const [isOpen, setIsOpen] = useState(defaultOpen)

  return (
    <div className="card">
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="w-full flex items-center justify-between p-4 text-left hover:bg-slate-800/50 transition-colors"
      >
        <h2 className="text-lg font-medium text-slate-100">{title}</h2>
        {isOpen ? (
          <ChevronDown className="w-5 h-5 text-slate-400" />
        ) : (
          <ChevronRight className="w-5 h-5 text-slate-400" />
        )}
      </button>
      {isOpen && <div className="p-4 pt-0 border-t border-slate-800">{children}</div>}
    </div>
  )
}

interface FieldRendererProps {
  name: string
  schema: SchemaProperty
  value: unknown
  onChange: (value: unknown) => void
  required?: boolean
}

function FieldRenderer({ name, schema, value, onChange, required }: FieldRendererProps) {
  const label = camelToTitle(name)
  const description = schema.description

  if (schema.type === 'object' && schema.properties) {
    return (
      <ObjectField
        name={name}
        schema={schema}
        value={(value as Record<string, unknown>) || {}}
        onChange={onChange}
      />
    )
  }

  if (schema.type === 'array') {
    return (
      <ArrayField
        name={name}
        schema={schema}
        value={(value as unknown[]) || []}
        onChange={onChange}
      />
    )
  }

  if (schema.type === 'boolean') {
    return (
      <div className="flex items-start gap-3">
        <input
          type="checkbox"
          id={name}
          checked={Boolean(value)}
          onChange={(e) => onChange(e.target.checked)}
          className="mt-1 h-4 w-4 rounded border-slate-700 bg-slate-900 text-ha-blue focus:ring-ha-blue/50"
        />
        <div>
          <label htmlFor={name} className="label cursor-pointer">
            {label}
            {required && <span className="text-ha-red ml-1">*</span>}
          </label>
          {description && <p className="text-xs text-slate-500 mt-0.5">{description}</p>}
        </div>
      </div>
    )
  }

  if (schema.enum) {
    return (
      <div>
        <label htmlFor={name} className="label block mb-2">
          {label}
          {required && <span className="text-ha-red ml-1">*</span>}
          {description && (
            <span className="ml-2 text-slate-500 font-normal" title={description}>
              <HelpCircle className="inline w-3 h-3" />
            </span>
          )}
        </label>
        <select
          id={name}
          value={String(value || '')}
          onChange={(e) => onChange(e.target.value || undefined)}
          className="input"
          required={required}
        >
          <option value="">Select...</option>
          {schema.enum.map((opt) => (
            <option key={opt} value={opt}>
              {opt}
            </option>
          ))}
        </select>
      </div>
    )
  }

  if (schema.type === 'integer' || schema.type === 'number') {
    return (
      <div>
        <label htmlFor={name} className="label block mb-2">
          {label}
          {required && <span className="text-ha-red ml-1">*</span>}
          {description && (
            <span className="ml-2 text-slate-500 font-normal" title={description}>
              <HelpCircle className="inline w-3 h-3" />
            </span>
          )}
        </label>
        <input
          type="number"
          id={name}
          value={value !== undefined && value !== null ? Number(value) : ''}
          onChange={(e) => onChange(e.target.value ? Number(e.target.value) : undefined)}
          min={schema.minimum}
          max={schema.maximum}
          className="input"
          required={required}
        />
        {(schema.minimum !== undefined || schema.maximum !== undefined) && (
          <p className="text-xs text-slate-500 mt-1">
            {schema.minimum !== undefined && `Min: ${schema.minimum}`}
            {schema.minimum !== undefined && schema.maximum !== undefined && ' | '}
            {schema.maximum !== undefined && `Max: ${schema.maximum}`}
          </p>
        )}
      </div>
    )
  }

  return (
    <div>
      <label htmlFor={name} className="label block mb-2">
        {label}
        {required && <span className="text-ha-red ml-1">*</span>}
        {description && (
          <span className="ml-2 text-slate-500 font-normal" title={description}>
            <HelpCircle className="inline w-3 h-3" />
          </span>
        )}
      </label>
      <input
        type="text"
        id={name}
        value={String(value || '')}
        onChange={(e) => onChange(e.target.value || undefined)}
        placeholder={schema.default !== undefined ? String(schema.default) : undefined}
        className="input"
        required={required}
      />
    </div>
  )
}

interface ObjectFieldProps {
  name: string
  schema: SchemaProperty
  value: Record<string, unknown>
  onChange: (value: unknown) => void
}

function ObjectField({ name, schema, value, onChange }: ObjectFieldProps) {
  const [isExpanded, setIsExpanded] = useState(Object.keys(value).length > 0)
  const properties = schema.properties || {}
  const required = schema.required || []

  const updateNestedField = (fieldName: string, fieldValue: unknown) => {
    onChange({ ...value, [fieldName]: fieldValue })
  }

  return (
    <div className="border border-slate-800 rounded-lg">
      <button
        type="button"
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full flex items-center justify-between p-3 text-left hover:bg-slate-800/50 transition-colors rounded-lg"
      >
        <span className="label">{camelToTitle(name)}</span>
        {isExpanded ? (
          <ChevronDown className="w-4 h-4 text-slate-400" />
        ) : (
          <ChevronRight className="w-4 h-4 text-slate-400" />
        )}
      </button>
      {isExpanded && (
        <div className="p-3 pt-0 space-y-3">
          {Object.entries(properties).map(([fieldName, fieldSchema]) => (
            <FieldRenderer
              key={fieldName}
              name={fieldName}
              schema={fieldSchema}
              value={value[fieldName]}
              onChange={(v) => updateNestedField(fieldName, v)}
              required={required.includes(fieldName)}
            />
          ))}
        </div>
      )}
    </div>
  )
}

interface ArrayFieldProps {
  name: string
  schema: SchemaProperty
  value: unknown[]
  onChange: (value: unknown) => void
}

function ArrayField({ name, schema, value, onChange }: ArrayFieldProps) {
  const itemSchema = schema.items

  const addItem = () => {
    if (itemSchema?.type === 'object') {
      onChange([...value, {}])
    } else if (itemSchema?.type === 'array') {
      onChange([...value, []])
    } else {
      onChange([...value, ''])
    }
  }

  const updateItem = (index: number, itemValue: unknown) => {
    const newValue = [...value]
    newValue[index] = itemValue
    onChange(newValue)
  }

  const removeItem = (index: number) => {
    onChange(value.filter((_, i) => i !== index))
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-2">
        <label className="label">{camelToTitle(name)}</label>
        <button
          type="button"
          onClick={addItem}
          className="btn btn-ghost text-xs"
        >
          <Plus className="w-3 h-3" />
          Add
        </button>
      </div>
      <div className="space-y-2">
        {value.map((item, index) => (
          <div key={index} className="flex items-start gap-2">
            <div className="flex-1">
              {itemSchema?.type === 'object' && itemSchema.properties ? (
                <div className="border border-slate-800 rounded-lg p-3 space-y-3">
                  {Object.entries(itemSchema.properties).map(([fieldName, fieldSchema]) => (
                    <FieldRenderer
                      key={fieldName}
                      name={fieldName}
                      schema={fieldSchema}
                      value={(item as Record<string, unknown>)[fieldName]}
                      onChange={(v) =>
                        updateItem(index, { ...(item as Record<string, unknown>), [fieldName]: v })
                      }
                      required={itemSchema.required?.includes(fieldName)}
                    />
                  ))}
                </div>
              ) : (
                <input
                  type="text"
                  value={String(item || '')}
                  onChange={(e) => updateItem(index, e.target.value)}
                  className="input"
                  placeholder={`Item ${index + 1}`}
                />
              )}
            </div>
            <button
              type="button"
              onClick={() => removeItem(index)}
              className="p-2 text-slate-500 hover:text-ha-red transition-colors"
            >
              <Trash2 className="w-4 h-4" />
            </button>
          </div>
        ))}
        {value.length === 0 && (
          <p className="text-sm text-slate-500 py-2">No items. Click "Add" to add one.</p>
        )}
      </div>
    </div>
  )
}

function camelToTitle(str: string): string {
  return str
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, (s) => s.toUpperCase())
    .trim()
}
