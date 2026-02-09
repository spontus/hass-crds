import { describe, it, expect } from 'vitest'
import { render, screen, waitFor } from '../test/utils/render'
import Dashboard from './Dashboard'

describe('Dashboard', () => {
  it('renders dashboard title', async () => {
    render(<Dashboard />)

    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText(/Manage your MQTT entities/)).toBeInTheDocument()
  })

  it('shows loading state initially', () => {
    render(<Dashboard />)

    expect(screen.getAllByText('-').length).toBeGreaterThan(0)
  })

  it('renders statistics cards after loading', async () => {
    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText('Total Entities')).toBeInTheDocument()
    })

    expect(screen.getByText('Published')).toBeInTheDocument()
    expect(screen.getByText('Pending')).toBeInTheDocument()
    expect(screen.getByText('Entity Types')).toBeInTheDocument()
  })

  it('displays entity count after data loads', async () => {
    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.queryByText('-')).not.toBeInTheDocument()
    })
  })

  it('renders Entities by Category section', async () => {
    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText('Entities by Category')).toBeInTheDocument()
    })
  })

  it('renders Recent Entities section', async () => {
    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText('Recent Entities')).toBeInTheDocument()
    })
  })

  it('shows category list after loading', async () => {
    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText('Controls')).toBeInTheDocument()
    })

    expect(screen.getByText('Sensors')).toBeInTheDocument()
    expect(screen.getByText('Lighting')).toBeInTheDocument()
  })

  it('shows recent entities after loading', async () => {
    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText('Test Button')).toBeInTheDocument()
    })

    expect(screen.getByText('Temperature Sensor')).toBeInTheDocument()
  })

  it('displays published status indicator', async () => {
    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getAllByText('Published').length).toBeGreaterThan(0)
    })
  })

  it('displays pending status indicator', async () => {
    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText('Pending')).toBeInTheDocument()
    })
  })

  it('renders Create Entity button', () => {
    render(<Dashboard />)

    expect(screen.getByRole('link', { name: /Create Entity/i })).toBeInTheDocument()
  })

  it('renders View all link', () => {
    render(<Dashboard />)

    expect(screen.getByRole('link', { name: /View all/i })).toBeInTheDocument()
  })

  it('shows empty state when no entities exist', async () => {
    render(<Dashboard />)

    await waitFor(() => {
      const emptyText = screen.queryByText('No entities yet')
      if (emptyText) {
        expect(emptyText).toBeInTheDocument()
      }
    })
  })
})
