import { describe, it, expect, vi, beforeEach } from 'vitest'
import { screen, waitFor, render } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import EntityCreatePage from './EntityCreatePage'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

const mockNavigate = vi.fn()

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  }
})

function renderWithRoute(route: string) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[route]}>
        <Routes>
          <Route path="/create/:kind" element={<EntityCreatePage />} />
          <Route path="/entities/:kind" element={<div>Entity List</div>} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  )
}

describe('EntityCreatePage', () => {
  beforeEach(() => {
    mockNavigate.mockClear()
  })

  describe('loading state', () => {
    it('shows loading state while schema loads', () => {
      renderWithRoute('/create/MQTTButton')

      expect(screen.getByText('Loading schema...')).toBeInTheDocument()
    })
  })

  describe('page structure', () => {
    it('renders page title with entity type', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByText('Create Button')).toBeInTheDocument()
      })
    })

    it('renders description for entity type', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByText('Stateless button')).toBeInTheDocument()
      })
    })

    it('renders back button', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.queryByText('Loading schema...')).not.toBeInTheDocument()
      })

      const backLink = screen.getByRole('link', { name: '' })
      expect(backLink).toHaveAttribute('href', '/entities/MQTTButton')
    })
  })

  describe('basic information form', () => {
    it('renders resource name input', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByText('Resource Name *')).toBeInTheDocument()
      })

      expect(screen.getByPlaceholderText('my-entity')).toBeInTheDocument()
    })

    it('renders namespace selector', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByText('Namespace *')).toBeInTheDocument()
      })

      const namespaceSelect = screen.getAllByRole('combobox')[0]
      expect(namespaceSelect).toBeInTheDocument()
    })

    it('validates resource name format', async () => {
      const user = userEvent.setup()
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByPlaceholderText('my-entity')).toBeInTheDocument()
      })

      const input = screen.getByPlaceholderText('my-entity')
      await user.type(input, 'Invalid Name!')

      expect(input).toHaveValue('invalid-name-')
    })

    it('shows format hint for resource name', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByText('Lowercase, alphanumeric, hyphens only')).toBeInTheDocument()
      })
    })
  })

  describe('dynamic form', () => {
    it('renders DynamicForm after schema loads', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByText('Entity Configuration')).toBeInTheDocument()
      })
    })

    it('shows schema fields', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByLabelText(/Command Topic/i)).toBeInTheDocument()
      })
    })
  })

  describe('preview panel', () => {
    it('renders Preview section', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByText('Preview')).toBeInTheDocument()
      })
    })

    it('renders Show YAML button', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Show YAML/i })).toBeInTheDocument()
      })
    })

    it('shows YAML preview when toggled', async () => {
      const user = userEvent.setup()
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Show YAML/i })).toBeInTheDocument()
      })

      await user.click(screen.getByRole('button', { name: /Show YAML/i }))

      expect(screen.getByText(/apiVersion: mqtt.home-assistant.io/)).toBeInTheDocument()
      expect(screen.getByText(/kind: MQTTButton/)).toBeInTheDocument()
    })
  })

  describe('submit button', () => {
    it('renders Create Entity button', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Create Entity/i })).toBeInTheDocument()
      })
    })

    it('button is disabled when name is empty', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Create Entity/i })).toBeInTheDocument()
      })

      expect(screen.getByRole('button', { name: /Create Entity/i })).toBeDisabled()
    })

    it('button is enabled when name is provided', async () => {
      const user = userEvent.setup()
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByPlaceholderText('my-entity')).toBeInTheDocument()
      })

      await user.type(screen.getByPlaceholderText('my-entity'), 'test-button')

      expect(screen.getByRole('button', { name: /Create Entity/i })).not.toBeDisabled()
    })
  })

  describe('form submission', () => {
    it('enables submit button when name is entered', async () => {
      const user = userEvent.setup()
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        expect(screen.getByPlaceholderText('my-entity')).toBeInTheDocument()
      })

      const submitButton = screen.getByRole('button', { name: /Create Entity/i })
      expect(submitButton).toBeDisabled()

      await user.type(screen.getByPlaceholderText('my-entity'), 'new-button')

      expect(submitButton).not.toBeDisabled()
    })
  })

  describe('namespace selection', () => {
    it('defaults to first available namespace', async () => {
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        const namespaceSelect = screen.getAllByRole('combobox')[0]
        expect(namespaceSelect).toHaveValue('default')
      })
    })

    it('allows changing namespace', async () => {
      const user = userEvent.setup()
      renderWithRoute('/create/MQTTButton')

      await waitFor(() => {
        const namespaceSelect = screen.getAllByRole('combobox')[0]
        expect(namespaceSelect).toBeInTheDocument()
      })

      const namespaceSelect = screen.getAllByRole('combobox')[0]
      await user.selectOptions(namespaceSelect, 'production')

      expect(namespaceSelect).toHaveValue('production')
    })
  })
})
