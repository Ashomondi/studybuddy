import { useEffect, useState } from 'react'
import { api } from './api/client'
import Sidebar from './components/Sidebar'
import Topbar from './components/Topbar'
import MobileNav from './components/MobileNav'
import Dashboard from './pages/Dashboard'
import Goals from './pages/Goals'
import PlanDetail from './pages/PlanDetail'
import Sessions from './pages/Sessions'
import Schedule from './pages/Schedule'

const titles = {
  dashboard: 'Dashboard',
  goals: 'Learning goals',
  plan: 'Learning plan',
  sessions: 'Study sessions',
  schedule: 'Study schedule',
}

export default function App() {
  const [view, setView] = useState('dashboard')
  const [apiStatus, setApiStatus] = useState('Checking…')

  useEffect(() => {
    api.health()
      .then((h) => setApiStatus(h.status === 'ok' ? 'Connected' : 'Unavailable'))
      .catch(() => setApiStatus('Unavailable'))
  }, [])

  let page = null
  if (view === 'dashboard') page = <Dashboard setView={setView} />
  else if (view === 'goals') page = <Goals />
  else if (view === 'plan') page = <PlanDetail />
  else if (view === 'sessions') page = <Sessions />
  else if (view === 'schedule') page = <Schedule />

  return (
    <div className="app-shell">
      <Sidebar view={view} setView={setView} />
      <div className="main">
        <Topbar title={titles[view]} apiStatus={apiStatus} />
        <main className="content">{page}</main>
      </div>
      <MobileNav view={view} setView={setView} />
    </div>
  )
}
