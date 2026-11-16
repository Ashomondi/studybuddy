const API_BASE = (import.meta.env.VITE_API_URL || 'http://localhost:8080').replace(/\/$/, '')

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...(options.headers || {}) },
    ...options,
  })
  const text = await response.text()
  let data = null
  if (text) {
    try { data = JSON.parse(text) } catch { data = { raw: text } }
  }
  if (!response.ok) {
    const message = data?.error?.message || response.statusText || 'Request failed'
    const err = new Error(message)
    err.status = response.status
    err.code = data?.error?.code
    err.data = data
    throw err
  }
  return data
}

export const api = {
  health: () => request('/api/health'),
  goals: {
    list: (status) => request(`/api/goals${status ? `?status=${encodeURIComponent(status)}` : ''}`),
    get: (id) => request(`/api/goals/${id}`),
    create: (payload) => request('/api/goals', { method: 'POST', body: JSON.stringify(payload) }),
    update: (id, payload) => request(`/api/goals/${id}`, { method: 'PUT', body: JSON.stringify(payload) }),
    remove: (id) => request(`/api/goals/${id}`, { method: 'DELETE' }),
  },
  plans: {
    list: (goalId) => request(`/api/plans${goalId ? `?goal_id=${goalId}` : ''}`),
    get: (id) => request(`/api/plans/${id}`),
    create: (payload) => request('/api/plans', { method: 'POST', body: JSON.stringify(payload) }),
    update: (id, payload) => request(`/api/plans/${id}`, { method: 'PUT', body: JSON.stringify(payload) }),
    remove: (id) => request(`/api/plans/${id}`, { method: 'DELETE' }),
    activities: (id) => request(`/api/plans/${id}/topics`),
  },
  activities: {
    create: (payload) => request('/api/activities', { method: 'POST', body: JSON.stringify(payload) }),
    update: (id, payload) => request(`/api/activities/${id}`, { method: 'PUT', body: JSON.stringify(payload) }),
    remove: (id) => request(`/api/activities/${id}`, { method: 'DELETE' }),
  },
  schedule: (from, to) => request(`/api/schedule?from=${from}&to=${to}`),
  sessions: {
    list: (params = {}) => {
      const q = new URLSearchParams()
      if (params.goalId) q.set('goal_id', params.goalId)
      if (params.status) q.set('status', params.status)
      const s = q.toString()
      return request(`/api/sessions${s ? `?${s}` : ''}`)
    },
    get: (id) => request(`/api/sessions/${id}`),
    create: (payload) => request('/api/sessions', { method: 'POST', body: JSON.stringify(payload) }),
  },
}
