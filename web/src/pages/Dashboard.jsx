import { useEffect, useState } from 'react'
import { api } from '../api/client'
import Card from '../components/Card'
import Metric from '../components/Metric'
import EmptyState from '../components/EmptyState'
import Button from '../components/Button'
import Icon from '../components/Icon'
import Badge from '../components/Badge'
import ProgressBar from '../components/ProgressBar'

export default function Dashboard({ setView }) {
  const [goals, setGoals] = useState([])
  const [sessions, setSessions] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    let cancelled = false
    async function load() {
      try {
        const [g, s] = await Promise.all([api.goals.list(), api.sessions.list()])
        if (!cancelled) { setGoals(g.goals || []); setSessions(s.sessions || []) }
      } catch (e) {
        if (!cancelled) setError(e.message)
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    load()
    return () => { cancelled = true }
  }, [])

  const activeGoals = goals.filter((g) => g.status === 'active')
  const running = sessions.filter((s) => s.status === 'running')

  if (loading) return <div className="loading">Loading dashboard…</div>
  if (error) return <div className="toast error">{error}</div>

  return (
    <div className="grid" style={{ gap: 20 }}>
      <div className="page-header">
        <div>
          <h1>Good to see you 👋</h1>
          <p>Here's your learning space for today.</p>
        </div>
        <Button onClick={() => setView('goals')}><Icon name="plus" /> New goal</Button>
      </div>

      <div className="grid grid-4">
        <Card><Metric value={activeGoals.length} label="Active goals" /></Card>
        <Card><Metric value={sessions.length} label="Study sessions" /></Card>
        <Card><Metric value={running.length} label="Running now" /></Card>
        <Card><Metric value="—" label="Focus minutes (sample)" /></Card>
      </div>

      <div className="grid grid-2">
        <Card>
          <h3 style={{ marginTop: 0, fontFamily: 'var(--font-heading)' }}>Active goals</h3>
          {activeGoals.length === 0 ? (
            <EmptyState title="No active goals" description="Create your first learning goal to get started." action={<Button onClick={() => setView('goals')}>Create goal</Button>} />
          ) : (
            activeGoals.slice(0, 4).map((g) => (
              <div key={g.id} className="schedule-item">
                <span className="schedule-time"><Badge status={g.status} /></span>
                <div style={{ flex: 1 }}>
                  <div style={{ fontWeight: 600 }}>{g.title}</div>
                  <div style={{ fontSize: 13, color: 'var(--color-muted)' }}>{g.subject || 'No subject'}</div>
                </div>
              </div>
            ))
          )}
        </Card>

        <Card>
          <h3 style={{ marginTop: 0, fontFamily: 'var(--font-heading)' }}>Today's schedule (sample)</h3>
          <p style={{ fontSize: 12, color: 'var(--color-muted)', marginTop: 0 }}>Sample data — schedule endpoints not yet implemented.</p>
          {['09:00 Review notes', '11:00 Practice problems', '14:00 Reading'].map((t) => (
            <div key={t} className="schedule-item"><span className="schedule-time">{t.slice(0,5)}</span><span>{t.slice(6)}</span></div>
          ))}
        </Card>
      </div>

      <Card>
        <h3 style={{ marginTop: 0, fontFamily: 'var(--font-heading)' }}>Weekly progress (placeholder)</h3>
        <ProgressBar value={42} />
        <p style={{ fontSize: 12, color: 'var(--color-muted)', marginTop: 8 }}>Sample progress indicator.</p>
      </Card>
    </div>
  )
}
