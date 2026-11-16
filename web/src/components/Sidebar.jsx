import Icon from './Icon'

const items = [
  { id: 'dashboard', label: 'Dashboard', icon: 'home' },
  { id: 'goals', label: 'Learning Goals', icon: 'target' },
  { id: 'plan', label: 'Learning Plan', icon: 'plan' },
  { id: 'schedule', label: 'Schedule', icon: 'clock' },
  { id: 'sessions', label: 'Study Sessions', icon: 'clock' },
]

export default function Sidebar({ view, setView }) {
  return (
    <aside className="sidebar">
      <div className="sidebar-brand">
        <span className="brand-mark">✳</span>
        <span>studybuddy<span className="brand-ai">.ai</span></span>
      </div>
      <nav className="nav-list">
        {items.map((it) => (
          <button key={it.id} className={`nav-item ${view === it.id ? 'active' : ''}`} onClick={() => setView(it.id)}>
            <Icon name={it.icon} /> {it.label}
          </button>
        ))}
      </nav>
      <div className="sidebar-footer">Small steps. Real momentum.</div>
    </aside>
  )
}
