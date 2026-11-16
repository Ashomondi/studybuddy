import { useEffect, useState } from 'react'
import { api } from '../api/client'
import Card from '../components/Card'
import Button from '../components/Button'
import Field from '../components/Field'
import Modal from '../components/Modal'
import EmptyState from '../components/EmptyState'
import Loading from '../components/Loading'
import Badge from '../components/Badge'
import Icon from '../components/Icon'
import Toast from '../components/Toast'

export default function Sessions() {
  const [sessions, setSessions] = useState([])
  const [goals, setGoals] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [toast, setToast] = useState(null)

  const [modalOpen, setModalOpen] = useState(false)
  const [goalId, setGoalId] = useState('')
  const [startedAt, setStartedAt] = useState('')
  const [saving, setSaving] = useState(false)
  const [formError, setFormError] = useState(null)

  async function load() {
    setLoading(true); setError(null)
    try {
      const [s, g] = await Promise.all([api.sessions.list(), api.goals.list()])
      setSessions(s.sessions || [])
      setGoals(g.goals || [])
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [])

  async function submit(e) {
    e.preventDefault()
    if (!goalId) { setFormError('Select a goal'); return }
    setSaving(true); setFormError(null)
    try {
      const payload = { goal_id: Number(goalId) }
      if (startedAt.trim()) payload.started_at = new Date(startedAt).toISOString()
      await api.sessions.create(payload)
      setToast('Session started')
      setModalOpen(false)
      setGoalId(''); setStartedAt('')
      await load()
    } catch (e) {
      setFormError(e.message)
    } finally {
      setSaving(false)
    }
  }

  const goalTitle = (id) => goals.find((g) => g.id === id)?.title || `Goal ${id}`

  if (loading) return <Loading />
  if (error) return <div className="toast error">{error}</div>

  return (
    <div>
      <div className="page-header">
        <div>
          <h1>Study sessions</h1>
          <p>Track your study time with reliable start/stop events.</p>
        </div>
        <Button onClick={() => setModalOpen(true)} disabled={goals.length === 0}><Icon name="plus" /> Start session</Button>
      </div>

      {sessions.length === 0 ? (
        <Card>
          <EmptyState title="No sessions yet" description={goals.length === 0 ? 'Create a learning goal first.' : 'Start a session to begin tracking.'} />
        </Card>
      ) : (
        <Card>
          {sessions.map((s) => (
            <div key={s.id} className="schedule-item">
              <span className="schedule-time">{new Date(s.started_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
              <div style={{ flex: 1 }}>
                <div style={{ fontWeight: 600 }}>{goalTitle(s.goal_id)}</div>
                <div style={{ fontSize: 13, color: 'var(--color-muted)' }}>{new Date(s.started_at).toLocaleDateString()}</div>
              </div>
              <Badge status={s.status} />
            </div>
          ))}
        </Card>
      )}

      <Modal open={modalOpen} title="Start session" onClose={() => setModalOpen(false)}>
        <form onSubmit={submit} className="grid" style={{ gap: 12 }}>
          <Field label="Goal" error={formError}>
            <select value={goalId} onChange={(e) => setGoalId(e.target.value)}>
              <option value="">Select a goal…</option>
              {goals.map((g) => <option key={g.id} value={g.id}>{g.title}</option>)}
            </select>
          </Field>
          <Field label="Start time (optional)" hint="Defaults to now">
            <input type="datetime-local" value={startedAt} onChange={(e) => setStartedAt(e.target.value)} />
          </Field>
          <div className="modal-actions">
            <Button variant="secondary" type="button" onClick={() => setModalOpen(false)}>Cancel</Button>
            <Button type="submit" disabled={saving}>{saving ? 'Starting…' : 'Start'}</Button>
          </div>
        </form>
      </Modal>

      <Toast message={toast} onClose={() => setToast(null)} />
    </div>
  )
}