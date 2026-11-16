export default function Badge({ status }) {
  const map = { active: 'active', in_progress: 'active', pending: 'paused', completed: 'completed', paused: 'paused', archived: 'archived', running: 'active' }
  const cls = map[status] || 'archived'
  return <span className={`badge badge-${cls}`}>{status}</span>
}
