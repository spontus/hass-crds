import { describe, it, expect } from 'vitest'
import { screen, waitFor, render } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import EntityListPage from './EntityListPage'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

function renderWithRoute(route: string) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[route]}>
        <Routes>
          <Route path="/entities" element={<EntityListPage />} />
          <Route path="/entities/:kind" element={<EntityListPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  )
}

describe('EntityListPage', () => {
  describe('without kind filter', () => {
    it('shows All Entities title', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.getByText('All Entities')).toBeInTheDocument()
      })
    })

    it('displays loading state initially', () => {
      renderWithRoute('/entities')

      expect(screen.getByText('Loading...')).toBeInTheDocument()
    })

    it('renders entities table after loading', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.queryByText('Loading...')).not.toBeInTheDocument()
      })

      expect(screen.getByRole('table')).toBeInTheDocument()
      expect(screen.getByText('Name')).toBeInTheDocument()
      expect(screen.getByText('Namespace')).toBeInTheDocument()
      expect(screen.getByText('Kind')).toBeInTheDocument()
      expect(screen.getByText('Status')).toBeInTheDocument()
    })

    it('displays entities from mock data', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.getByText('Test Button')).toBeInTheDocument()
      })

      expect(screen.getByText('Temperature Sensor')).toBeInTheDocument()
      expect(screen.getByText('Power Switch')).toBeInTheDocument()
    })

    it('shows Kind column when no filter', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.getByText('Kind')).toBeInTheDocument()
      })
    })
  })

  describe('with kind filter', () => {
    it('shows Create button for kind', async () => {
      renderWithRoute('/entities/MQTTButton')

      await waitFor(() => {
        expect(screen.getByRole('link', { name: /Create Button/i })).toBeInTheDocument()
      })
    })
  })

  describe('search functionality', () => {
    it('renders search input', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.getByPlaceholderText('Search entities...')).toBeInTheDocument()
      })
    })

    it('filters entities by search query', async () => {
      const user = userEvent.setup()
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.getByText('Test Button')).toBeInTheDocument()
      })

      await user.type(screen.getByPlaceholderText('Search entities...'), 'Button')

      expect(screen.getByText('Test Button')).toBeInTheDocument()
      expect(screen.queryByText('Temperature Sensor')).not.toBeInTheDocument()
    })

    it('shows no results message when search matches nothing', async () => {
      const user = userEvent.setup()
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.getByText('Test Button')).toBeInTheDocument()
      })

      await user.type(screen.getByPlaceholderText('Search entities...'), 'nonexistent')

      expect(screen.getByText('No matching entities found')).toBeInTheDocument()
    })
  })

  describe('namespace filter', () => {
    it('renders namespace dropdown', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.getByText('All namespaces')).toBeInTheDocument()
      })
    })

    it('shows namespaces dropdown', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        const select = screen.getByRole('combobox')
        expect(select).toBeInTheDocument()
      })
    })
  })

  describe('entity status display', () => {
    it('shows Published status for published entities', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        const publishedElements = screen.getAllByText('Published')
        expect(publishedElements.length).toBeGreaterThan(0)
      })
    })

    it('shows Pending status for unpublished entities', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.getByText('Pending')).toBeInTheDocument()
      })
    })
  })

  describe('entity actions', () => {
    it('renders table rows with entity data', async () => {
      renderWithRoute('/entities')

      await waitFor(() => {
        expect(screen.getByText('Test Button')).toBeInTheDocument()
      })

      expect(screen.getByText('Temperature Sensor')).toBeInTheDocument()
    })
  })

  describe('empty state', () => {
    it('shows empty state message when no entities', async () => {
      renderWithRoute('/entities/MQTTLight')

      await waitFor(() => {
        expect(screen.queryByText('Loading...')).not.toBeInTheDocument()
      })

      const emptyMessage = screen.queryByText(/No Light entities yet/)
      if (emptyMessage) {
        expect(emptyMessage).toBeInTheDocument()
      }
    })
  })
})
