import { describe, it, expect, beforeEach, vi } from 'vitest'
import apiClient from '@/service/AxiosClient'

describe('AxiosClient.ts', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.restoreAllMocks()
  })

  it('creates an axios client with correct base configuration', () => {
    expect(apiClient.defaults.timeout).toBe(10000)
    expect(apiClient.defaults.headers.common['Accept']).toBe('application/json, text/plain, */*')
    expect(apiClient.defaults.headers['Content-Type']).toBe('application/json')
    expect(apiClient.defaults.baseURL).toBeDefined()
  })

  it('adds Authorization header in request interceptor when token exists in localStorage', async () => {
    localStorage.setItem('kd5f_token', 'mock-jwt-token-123')

    const config = { headers: {} as Record<string, string> }
    const handlers = (apiClient.interceptors.request as any).handlers
    const requestHandler = handlers[handlers.length - 1].fulfilled

    const updatedConfig = await requestHandler(config)
    expect(updatedConfig.headers.Authorization).toBe('Bearer mock-jwt-token-123')
  })

  it('does not add Authorization header in request interceptor when token is missing', async () => {
    const config = { headers: {} as Record<string, string> }
    const handlers = (apiClient.interceptors.request as any).handlers
    const requestHandler = handlers[handlers.length - 1].fulfilled

    const updatedConfig = await requestHandler(config)
    expect(updatedConfig.headers.Authorization).toBeUndefined()
  })
})