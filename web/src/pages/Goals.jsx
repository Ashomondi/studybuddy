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

const emptyForm = { title: '', subject: '', description: '', status: 'active', target_date: '', current_level: '', weekly_study_minutes: '', preferred_study_days: [] }

function validate(form) {
  const errors = {}
  if (!form.title.trim()) errors.title = 'Title is required'
  if (form.target_date && !/^\d{4}-\d{2}-\d{2}$/.test(form.target_date)) errors.target_date = 'Use YYYY-MM-DD'
  if (!['active', 'paused', 'completed', 'archived'].includes(form.status)) errors.status = 'Invalid status'
  return errors
}

export default function Goals() {
  const [goals, setGoals] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [toast, setToast] = useState(null)

  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState(null)
  const [form, setForm] = useState(emptyForm)
  const [errors, setErrors] = useState({})
  const [formError, setFormError] = useState('')
  const [saving, setSaving] = useState(false)

  const [deleteTarget, setDeleteTarget] = useState(null)

  async function load() {
    setLoading(true); setError(null)
    try {
      const data = await api.goals.list()
      setGoals(data.goals || [])
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [])

  function openCreate() {
    setEditing(null); setForm(emptyForm); setErrors({}); setFormError(''); setModalOpen(true)
  }
  function openEdit(g) {
    setEditing(g)
    setForm({
      title: g.title || '',
      subject: g.subject || '',
      description: g.description || '',
      status: g.status || 'active',
      target_date: g.target_date || '',
      current_level: g.current_level || '', weekly_study_minutes: g.weekly_study_minutes || '', preferred_study_days: g.preferred_study_days || [],
    })
    setErrors({})
    setFormError('')
    setModalOpen(true)
  }

  async function submit(e) {
    e.preventDefault()
    const errs = validate(form)
    setErrors(errs)
    setFormError('')
    if (Object.keys(errs).length) return
    setSaving(true)
    try {
      const payload = {
        title: form.title.trim(),
        subject: form.subject.trim() || null,
        description: form.description.trim() || null,
        status: form.status,
        target_date: form.target_date.trim() || null,
        current_level: form.current_level.trim() || null,
        weekly_study_minutes: form.weekly_study_minutes ? Number(form.weekly_study_minutes) : null,
        preferred_study_days: form.preferred_study_days,
      }
      let legacyAPI = false
      if (editing) {
        try {
          await api.goals.update(editing.id, payload)
        } catch (e) {
          if (e.status !== 400 || e.message !== 'invalid JSON body') throw e
          await api.goals.update(editing.id, {
            title: payload.title, subject: payload.subject, description: payload.description,
            status: payload.status, target_date: payload.target_date,
          })
          legacyAPI = true
        }
        setToast('Goal updated')
      } else {
        try {
          await api.goals.create(payload)
        } catch (e) {
          if (e.status !== 400 || e.message !== 'invalid JSON body') throw e
          await api.goals.create({
            title: payload.title, subject: payload.subject, description: payload.description,
            status: payload.status, target_date: payload.target_date,
          })
          legacyAPI = true
        }
        setToast(legacyAPI ? 'Goal saved. Restart the API to save study preferences.' : 'Goal created')
      }
      if (editing && legacyAPI) setToast('Goal updated. Restart the API to save study preferences.')
      setModalOpen(false)
      await load()
    } catch (e) {
      // Keep API/JSON errors separate from title validation so the message
      // does not misleadingly appear under the title field.
      setFormError(e.message)
    } finally {
      setSaving(false)
    }
  }

  async function confirmDelete() {
    if (!deleteTarget) return
    try {
      await api.goals.remove(deleteTarget.id)
      setToast('Goal deleted')
      setDeleteTarget(null)
      await load()
    } catch (e) {
      setToast(e.message)
    }
  }

  if (loading) return <Loading />
  if (error) return <div className="toast error">{error}</div>

  return (
    <div>
      <div className="page-header">
        <div>
          <h1>Learning goals</h1>
          <p>Define what you want to learn and track progress.</p>
        </div>
        <Button onClick={openCreate}><Icon name="plus" /> New goal</Button>
      </div>

      {goals.length === 0 ? (
        <Card>
          <EmptyState title="No goals yet" description="Create your first learning goal to get started." action={<Button onClick={openCreate}>Create goal</Button>} />
        </Card>
      ) : (
        <div className="grid grid-2">
          {goals.map((g) => (
            <Card key={g.id}>
              <div className="goal-card">
                <div className="row">
                  <h3>{g.title}</h3>
                  <Badge status={g.status} />
                </div>
                {g.subject && <div className="meta"><span>{g.subject}</span></div>}
                {g.description && <p style={{ margin: 0, fontSize: 14, color: 'var(--color-muted)' }}>{g.description}</p>}
                <div className="meta">
                  {g.target_date && <span>Target: {g.target_date}</span>}
                </div>
                <div className="row">
                  <Button variant="secondary" size="sm" onClick={() => openEdit(g)}><Icon name="edit" /> Edit</Button>
                  <Button variant="danger" size="sm" onClick={() => setDeleteTarget(g)}><Icon name="trash" /> Delete</Button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      <Modal open={modalOpen} title={editing ? 'Edit goal' : 'Create goal'} onClose={() => setModalOpen(false)}>
        <form onSubmit={submit} className="grid" style={{ gap: 12 }}>
          {formError && <div className="form-error" role="alert">Could not save goal: {formError}</div>}
          <Field label="Title" error={errors.title}>
            <input value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} placeholder="e.g. Learn Go" />
          </Field>
          <Field label="Subject">
            <input value={form.subject} onChange={(e) => setForm({ ...form, subject: e.target.value })} placeholder="e.g. Programming" />
          </Field>
          <Field label="Description">
            <textarea value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
          </Field>
          <Field label="Status" error={errors.status}>
            <select value={form.status} onChange={(e) => setForm({ ...form, status: e.target.value })}>
              <option value="active">Active</option>
              <option value="paused">Paused</option>
              <option value="completed">Completed</option>
              <option value="archived">Archived</option>
            </select>
          </Field>
          <Field label="Target date" error={errors.target_date} hint="YYYY-MM-DD">
            <input value={form.target_date} onChange={(e) => setForm({ ...form, target_date: e.target.value })} placeholder="2026-12-01" />
          </Field>
          <Field label="Current level"><select value={form.current_level} onChange={(e) => setForm({ ...form, current_level: e.target.value })}><option value="">Choose level</option><option>Beginner</option><option>Intermediate</option><option>Advanced</option></select></Field>
          <Field label="Available study time (minutes per week)"><input type="number" min="1" max="10080" value={form.weekly_study_minutes} onChange={(e) => setForm({ ...form, weekly_study_minutes: e.target.value })} /></Field>
          <Field label="Preferred study days"><div className="day-options">{['Sun','Mon','Tue','Wed','Thu','Fri','Sat'].map((day, index) => <label key={day}><input type="checkbox" checked={form.preferred_study_days.includes(index)} onChange={(e) => setForm({ ...form, preferred_study_days: e.target.checked ? [...form.preferred_study_days, index].sort() : form.preferred_study_days.filter((d) => d !== index) })} />{day}</label>)}</div></Field>
          <div className="modal-actions">
            <Button variant="secondary" type="button" onClick={() => setModalOpen(false)}>Cancel</Button>
            <Button type="submit" disabled={saving}>{saving ? 'Saving…' : (editing ? 'Save changes' : 'Create goal')}</Button>
          </div>
        </form>
      </Modal>

      <Modal open={!!deleteTarget} title="Delete goal?" onClose={() => setDeleteTarget(null)}>
        <p style={{ color: 'var(--color-muted)', marginTop: 0 }}>This will permanently delete "{deleteTarget?.title}" and its sessions.</p>
        <div className="modal-actions">
          <Button variant="secondary" onClick={() => setDeleteTarget(null)}>Cancel</Button>
          <Button variant="danger" onClick={confirmDelete}>Delete</Button>
        </div>
      </Modal>

      <Toast message={toast} onClose={() => setToast(null)} />
    </div>
  )
}
