import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '../../test/utils/render'
import userEvent from '@testing-library/user-event'
import DynamicForm from './DynamicForm'
import type { SchemaProperty } from '../../types/entity'

const baseSchema: SchemaProperty = {
  type: 'object',
  required: ['commandTopic'],
  properties: {
    commandTopic: {
      type: 'string',
      description: 'MQTT topic',
    },
    name: {
      type: 'string',
      description: 'Display name',
    },
  },
}

describe('DynamicForm', () => {
  describe('field rendering', () => {
    it('renders text input for string fields', () => {
      const onChange = vi.fn()
      render(
        <DynamicForm
          schema={baseSchema}
          value={{}}
          onChange={onChange}
        />
      )

      expect(screen.getByLabelText(/Command Topic/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/Command Topic/i)).toHaveAttribute('type', 'text')
    })

    it('renders number input with min/max constraints', () => {
      const schema: SchemaProperty = {
        type: 'object',
        required: ['count'],
        properties: {
          count: {
            type: 'integer',
            description: 'Count value',
            minimum: 0,
            maximum: 10,
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      const input = screen.getByLabelText(/Count/i)
      expect(input).toHaveAttribute('type', 'number')
      expect(input).toHaveAttribute('min', '0')
      expect(input).toHaveAttribute('max', '10')
    })

    it('renders select for enum fields', () => {
      const schema: SchemaProperty = {
        type: 'object',
        required: ['status'],
        properties: {
          status: {
            type: 'string',
            enum: ['active', 'inactive', 'pending'],
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByRole('combobox')).toBeInTheDocument()
      expect(screen.getByText('active')).toBeInTheDocument()
      expect(screen.getByText('inactive')).toBeInTheDocument()
      expect(screen.getByText('pending')).toBeInTheDocument()
    })

    it('renders checkbox for boolean fields', () => {
      const schema: SchemaProperty = {
        type: 'object',
        required: ['enabled'],
        properties: {
          enabled: {
            type: 'boolean',
            description: 'Enable this feature',
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByRole('checkbox')).toBeInTheDocument()
    })

    it('shows required indicator for required fields', () => {
      const onChange = vi.fn()
      render(
        <DynamicForm
          schema={baseSchema}
          value={{}}
          onChange={onChange}
        />
      )

      const label = screen.getByText('Command Topic')
      expect(label.parentElement?.innerHTML).toContain('*')
    })
  })

  describe('field categorization', () => {
    it('groups fields into Entity Configuration section', () => {
      render(<DynamicForm schema={baseSchema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByText('Entity Configuration')).toBeInTheDocument()
    })

    it('groups common fields into Common Settings section', () => {
      const schema: SchemaProperty = {
        type: 'object',
        properties: {
          icon: { type: 'string' },
          uniqueId: { type: 'string' },
          entityCategory: { type: 'string' },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByText('Common Settings')).toBeInTheDocument()
    })

    it('groups device fields into Device Configuration section', () => {
      const schema: SchemaProperty = {
        type: 'object',
        properties: {
          device: {
            type: 'object',
            properties: {
              name: { type: 'string' },
            },
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByText('Device Configuration')).toBeInTheDocument()
    })

    it('groups MQTT fields into Advanced MQTT Settings section', () => {
      const schema: SchemaProperty = {
        type: 'object',
        properties: {
          qos: { type: 'integer' },
          retain: { type: 'boolean' },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByText('Advanced MQTT Settings')).toBeInTheDocument()
    })
  })

  describe('nested objects', () => {
    it('renders section header for object fields in device category', () => {
      const schema: SchemaProperty = {
        type: 'object',
        properties: {
          device: {
            type: 'object',
            properties: {
              name: { type: 'string' },
              manufacturer: { type: 'string' },
            },
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByText('Device Configuration')).toBeInTheDocument()
    })

    it('expands nested object when section and field are clicked', async () => {
      const user = userEvent.setup()
      const schema: SchemaProperty = {
        type: 'object',
        properties: {
          device: {
            type: 'object',
            properties: {
              name: { type: 'string', description: 'Device name' },
            },
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      await user.click(screen.getByText('Device Configuration'))
      await user.click(screen.getByText('Device'))

      expect(screen.getByLabelText(/Name/i)).toBeInTheDocument()
    })
  })

  describe('array fields', () => {
    it('renders add button for array fields', () => {
      const schema: SchemaProperty = {
        type: 'object',
        required: ['tags'],
        properties: {
          tags: {
            type: 'array',
            items: { type: 'string' },
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByRole('button', { name: /add/i })).toBeInTheDocument()
    })

    it('shows empty state message for empty array fields', () => {
      const schema: SchemaProperty = {
        type: 'object',
        required: ['tags'],
        properties: {
          tags: {
            type: 'array',
            items: { type: 'string' },
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByText(/No items/i)).toBeInTheDocument()
    })

    it('calls onChange when add button is clicked', async () => {
      const user = userEvent.setup()
      const onChange = vi.fn()
      const schema: SchemaProperty = {
        type: 'object',
        required: ['tags'],
        properties: {
          tags: {
            type: 'array',
            items: { type: 'string' },
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={onChange} />)

      await user.click(screen.getByRole('button', { name: /add/i }))

      expect(onChange).toHaveBeenCalled()
    })

    it('renders input for each array item', () => {
      const schema: SchemaProperty = {
        type: 'object',
        required: ['tags'],
        properties: {
          tags: {
            type: 'array',
            items: { type: 'string' },
          },
        },
      }

      render(
        <DynamicForm
          schema={schema}
          value={{ tags: ['item1', 'item2'] }}
          onChange={vi.fn()}
        />
      )

      expect(screen.getByDisplayValue('item1')).toBeInTheDocument()
      expect(screen.getByDisplayValue('item2')).toBeInTheDocument()
    })
  })

  describe('onChange callback', () => {
    it('calls onChange when text input value changes', async () => {
      const user = userEvent.setup()
      const onChange = vi.fn()

      render(
        <DynamicForm
          schema={baseSchema}
          value={{}}
          onChange={onChange}
        />
      )

      const input = screen.getByLabelText(/Command Topic/i)
      await user.type(input, 'a')

      expect(onChange).toHaveBeenCalled()
    })

    it('calls onChange when checkbox is toggled', async () => {
      const user = userEvent.setup()
      const onChange = vi.fn()
      const schema: SchemaProperty = {
        type: 'object',
        required: ['enabled'],
        properties: {
          enabled: { type: 'boolean' },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={onChange} />)

      await user.click(screen.getByRole('checkbox'))

      expect(onChange).toHaveBeenCalledWith({ enabled: true })
    })

    it('calls onChange when select value changes', async () => {
      const user = userEvent.setup()
      const onChange = vi.fn()
      const schema: SchemaProperty = {
        type: 'object',
        required: ['status'],
        properties: {
          status: {
            type: 'string',
            enum: ['active', 'inactive'],
          },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={onChange} />)

      await user.selectOptions(screen.getByRole('combobox'), 'inactive')

      expect(onChange).toHaveBeenCalledWith({ status: 'inactive' })
    })

    it('calls onChange when number input value changes', async () => {
      const user = userEvent.setup()
      const onChange = vi.fn()
      const schema: SchemaProperty = {
        type: 'object',
        required: ['count'],
        properties: {
          count: { type: 'integer', minimum: 0, maximum: 10 },
        },
      }

      render(<DynamicForm schema={schema} value={{}} onChange={onChange} />)

      const input = screen.getByLabelText(/Count/i)
      await user.clear(input)
      await user.type(input, '5')

      expect(onChange).toHaveBeenCalled()
    })
  })

  describe('value display', () => {
    it('displays existing string values', () => {
      render(
        <DynamicForm
          schema={baseSchema}
          value={{ commandTopic: 'existing/topic' }}
          onChange={vi.fn()}
        />
      )

      expect(screen.getByLabelText(/Command Topic/i)).toHaveValue('existing/topic')
    })

    it('displays existing boolean values', () => {
      const schema: SchemaProperty = {
        type: 'object',
        required: ['enabled'],
        properties: {
          enabled: { type: 'boolean' },
        },
      }

      render(
        <DynamicForm
          schema={schema}
          value={{ enabled: true }}
          onChange={vi.fn()}
        />
      )

      expect(screen.getByRole('checkbox')).toBeChecked()
    })

    it('displays existing number values', () => {
      const schema: SchemaProperty = {
        type: 'object',
        required: ['count'],
        properties: {
          count: { type: 'integer' },
        },
      }

      render(
        <DynamicForm schema={schema} value={{ count: 2 }} onChange={vi.fn()} />
      )

      expect(screen.getByLabelText(/Count/i)).toHaveValue(2)
    })

    it('displays existing select values', () => {
      const schema: SchemaProperty = {
        type: 'object',
        required: ['status'],
        properties: {
          status: {
            type: 'string',
            enum: ['active', 'inactive'],
          },
        },
      }

      render(
        <DynamicForm
          schema={schema}
          value={{ status: 'inactive' }}
          onChange={vi.fn()}
        />
      )

      expect(screen.getByRole('combobox')).toHaveValue('inactive')
    })
  })

  describe('section expand/collapse', () => {
    it('expands Entity Configuration by default', () => {
      render(<DynamicForm schema={baseSchema} value={{}} onChange={vi.fn()} />)

      expect(screen.getByLabelText(/Command Topic/i)).toBeVisible()
    })

    it('collapses section when header is clicked', async () => {
      const user = userEvent.setup()
      render(<DynamicForm schema={baseSchema} value={{}} onChange={vi.fn()} />)

      await user.click(screen.getByText('Entity Configuration'))

      expect(screen.queryByLabelText(/Command Topic/i)).not.toBeInTheDocument()
    })
  })
})
