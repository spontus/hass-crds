import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { screen, waitFor, render } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import EntityEditPage from './EntityEditPage'
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
          <Route path="/edit/:kind/:namespace/:name" element={<EntityEditPage />} />
          <Route path="/entities/:kind" element={<div>Entity List</div>} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  )
}

describe('EntityEditPage', () => {
  beforeEach(() => {
    mockNavigate.mockClear()
    vi.spyOn(window, 'confirm').mockReturnValue(true)
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('loading state', () => {
    it('shows loading state while data loads', () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      expect(screen.getByText('Loading...')).toBeInTheDocument()
    })
  })

  describe('page structure', () => {
    it('renders page title with entity type', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByText('Edit Button')).toBeInTheDocument()
      })
    })

    it('renders namespace/name subtitle', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByText('default/test-button')).toBeInTheDocument()
      })
    })

    it('renders back button', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.queryByText('Loading...')).not.toBeInTheDocument()
      })

      const backLink = screen.getByRole('link', { name: '' })
      expect(backLink).toHaveAttribute('href', '/entities/MQTTButton')
    })

    it('renders delete button', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Delete/i })).toBeInTheDocument()
      })
    })
  })

  describe('resource info section', () => {
    it('displays readonly name field', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        const nameInput = screen.getByDisplayValue('test-button')
        expect(nameInput).toBeInTheDocument()
        expect(nameInput).toBeDisabled()
      })
    })

    it('displays readonly namespace field', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        const nsInput = screen.getByDisplayValue('default')
        expect(nsInput).toBeInTheDocument()
        expect(nsInput).toBeDisabled()
      })
    })
  })

  describe('dynamic form', () => {
    it('renders DynamicForm after data loads', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByText('Entity Configuration')).toBeInTheDocument()
      })
    })

    it('populates form with existing entity data', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByDisplayValue('Test Button')).toBeInTheDocument()
      })

      expect(screen.getByDisplayValue('homeassistant/button/test/command')).toBeInTheDocument()
    })
  })

  describe('status panel', () => {
    it('renders Status section', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByText('Status')).toBeInTheDocument()
      })
    })

    it('shows Published status', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByText('Published')).toBeInTheDocument()
      })

      const publishedValue = screen.getByText('Yes')
      expect(publishedValue).toBeInTheDocument()
    })

    it('shows Created date', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByText('Created')).toBeInTheDocument()
      })
    })

    it('shows Discovery Topic', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByText('Discovery Topic')).toBeInTheDocument()
      })

      expect(
        screen.getByText('homeassistant/button/default-test-button/config')
      ).toBeInTheDocument()
    })
  })

  describe('submit button', () => {
    it('renders Save Changes button', async () => {
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Save Changes/i })).toBeInTheDocument()
      })
    })
  })

  describe('form submission', () => {
    it('submits form and navigates on success', async () => {
      const user = userEvent.setup()
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Save Changes/i })).toBeInTheDocument()
      })

      await user.click(screen.getByRole('button', { name: /Save Changes/i }))

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/entities/MQTTButton')
      })
    })
  })

  describe('delete functionality', () => {
    it('shows confirmation dialog on delete', async () => {
      const confirmSpy = vi.spyOn(window, 'confirm')
      const user = userEvent.setup()
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Delete/i })).toBeInTheDocument()
      })

      await user.click(screen.getByRole('button', { name: /Delete/i }))

      expect(confirmSpy).toHaveBeenCalledWith('Delete test-button?')
    })

    it('deletes entity and navigates on confirm', async () => {
      const user = userEvent.setup()
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Delete/i })).toBeInTheDocument()
      })

      await user.click(screen.getByRole('button', { name: /Delete/i }))

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/entities/MQTTButton')
      })
    })

    it('does not delete when confirm is cancelled', async () => {
      vi.spyOn(window, 'confirm').mockReturnValue(false)
      const user = userEvent.setup()
      renderWithRoute('/edit/MQTTButton/default/test-button')

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Delete/i })).toBeInTheDocument()
      })

      await user.click(screen.getByRole('button', { name: /Delete/i }))

      expect(mockNavigate).not.toHaveBeenCalled()
    })
  })

  describe('not found state', () => {
    it('shows not found message for nonexistent entity', async () => {
      renderWithRoute('/edit/MQTTButton/default/nonexistent')

      await waitFor(() => {
        expect(screen.getByText('Entity not found')).toBeInTheDocument()
      })
    })

    it('shows back to list link on not found', async () => {
      renderWithRoute('/edit/MQTTButton/default/nonexistent')

      await waitFor(() => {
        expect(screen.getByRole('link', { name: /Back to list/i })).toBeInTheDocument()
      })
    })
  })
})
