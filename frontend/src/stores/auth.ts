import { defineStore } from 'pinia'
import client from '../api/client'

export interface CurrentUser {
  id: number
  username: string
  display_name: string
  permissions: string[]
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as CurrentUser | null,
    loading: false
  }),
  getters: {
    isAuthenticated: () => Boolean(localStorage.getItem('access_token')),
    permissions: (state) => state.user?.permissions ?? []
  },
  actions: {
    async login(username: string, password: string) {
      this.loading = true
      try {
        const { data } = await client.post('/auth/login', { username, password })
        localStorage.setItem('access_token', data.data.access_token)
        localStorage.setItem('refresh_token', data.data.refresh_token)
        await this.fetchMe()
      } finally {
        this.loading = false
      }
    },
    async register(payload: Record<string, unknown>) {
      await client.post('/auth/register', payload)
    },
    async fetchMe() {
      const { data } = await client.get('/me')
      this.user = data.data
      return this.user
    },
    async logout() {
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        await client.post('/auth/logout', { refresh_token: refreshToken }).catch(() => undefined)
      }
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      this.user = null
    }
  }
})

