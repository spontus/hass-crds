import { describe, it, expect } from 'vitest'
import { api, ApiError } from './client'

describe('api client', () => {
  describe('get', () => {
    it('returns entity types successfully', async () => {
      const response = await api.get<{ entityTypes: unknown[] }>('/entity-types')

      expect(response.entityTypes).toBeDefined()
      expect(Array.isArray(response.entityTypes)).toBe(true)
    })

    it('returns namespaces successfully', async () => {
      const response = await api.get<{ namespaces: unknown[] }>('/namespaces')

      expect(response.namespaces).toBeDefined()
      expect(Array.isArray(response.namespaces)).toBe(true)
    })

    it('returns entities successfully', async () => {
      const response = await api.get<{ items: unknown[] }>('/entities')

      expect(response.items).toBeDefined()
      expect(Array.isArray(response.items)).toBe(true)
    })
  })

  describe('post', () => {
    it('creates entity successfully', async () => {
      const response = await api.post<{ kind: string }>('/entities/MQTTButton/default', {
        metadata: { name: 'new-button' },
        spec: { name: 'New Button', commandTopic: 'test/topic' },
      })

      expect(response.kind).toBe('MQTTButton')
    })
  })

  describe('put', () => {
    it('updates entity successfully', async () => {
      const response = await api.put<{ spec: { name: string } }>(
        '/entities/MQTTButton/default/test-button',
        {
          metadata: { name: 'test-button' },
          spec: { name: 'Updated Button' },
        }
      )

      expect(response.spec.name).toBe('Updated Button')
    })
  })

  describe('delete', () => {
    it('deletes entity successfully', async () => {
      const response = await api.delete('/entities/MQTTButton/default/test-button')

      expect(response).toBeUndefined()
    })
  })

  describe('error handling', () => {
    it('throws ApiError for 404 responses', async () => {
      await expect(
        api.get('/entities/MQTTButton/default/nonexistent')
      ).rejects.toBeInstanceOf(ApiError)
    })

    it('includes status in ApiError', async () => {
      try {
        await api.get('/entities/MQTTButton/default/nonexistent')
        expect.fail('should have thrown')
      } catch (error) {
        expect(error).toBeInstanceOf(ApiError)
        expect((error as ApiError).status).toBe(404)
      }
    })
  })
})
