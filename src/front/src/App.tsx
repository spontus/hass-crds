import { Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/layout/Layout'
import Dashboard from './pages/Dashboard'
import EntityListPage from './pages/EntityListPage'
import EntityCreatePage from './pages/EntityCreatePage'
import EntityEditPage from './pages/EntityEditPage'

export default function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/entities" element={<EntityListPage />} />
        <Route path="/entities/:kind" element={<EntityListPage />} />
        <Route path="/create/:kind" element={<EntityCreatePage />} />
        <Route path="/edit/:kind/:namespace/:name" element={<EntityEditPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Layout>
  )
}
