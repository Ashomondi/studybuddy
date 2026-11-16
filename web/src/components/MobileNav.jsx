import Icon from './Icon'

const items = [
  { id: 'dashboard', label: 'Home', icon: 'home' },
  { id: 'goals', label: 'Goals', icon: 'target' },
  { id: 'plan', label: 'Plan', icon: 'plan' },
  { id: 'schedule', label: 'Schedule', icon: 'clock' },
  { id: 'sessions', label: 'Sessions', icon: 'clock' },
]

export default function MobileNav({ view, setView }) {
  return (
    <nav className="mobile-nav">
      {items.map((it) => (
        <button key={it.id} className={view === it.id ? 'active' : ''} onClick={() => setView(it.id)}>
          <Icon name={it.icon} /> {it.label}
        </button>
      ))}
    </nav>
  )
}
