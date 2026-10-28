import React from 'react'
import { createRoot } from 'react-dom/client'
import './style.css'

const API_BASE = (import.meta.env.VITE_API_URL || 'http://localhost:8080').replace(/\/$/, '')

function App() {
  const [apiStatus, setApiStatus] = React.useState('Checking…')

  React.useEffect(() => {
    fetch(`${API_BASE}/api/health`)
      .then((response) => {
        if (!response.ok) throw new Error('API unavailable')
        return response.json()
      })
      .then((health) => setApiStatus(health.status === 'ok' ? 'Connected' : 'Unavailable'))
      .catch(() => setApiStatus('Unavailable'))
  }, [])

  return (
    <main className="page-shell">
      <header className="topbar">
        <a className="brand" href="#home" aria-label="StudyBuddy home">
          <span className="brand-mark">✳</span>
          <span>studybuddy<span className="brand-ai">.ai</span></span>
        </a>
        <div className="profile"><span className="status-dot" /> Learning space <span className="avatar">S</span></div>
      </header>

      <section className="welcome" id="home">
        <div className="eyebrow">YOUR LEARNING SPACE</div>
        <h1>Make room for<br /><span>what’s next.</span></h1>
        <p>A little progress, consistently, adds up to something big.</p>
        <button className="primary-button" type="button" disabled>
          <span aria-hidden="true">＋</span> Create a learning goal
        </button>
      </section>

      <section className="empty-card" aria-label="Your goals">
        <div className="empty-art" aria-hidden="true">
          <div className="sun" /><div className="book book-back" /><div className="book book-front" />
          <div className="sparkle sparkle-one">✦</div><div className="sparkle sparkle-two">✳</div>
        </div>
        <div className="empty-copy">
          <div className="eyebrow">A FRESH START</div>
          <h2>Your next chapter<br />starts with one goal.</h2>
          <p>Choose something you’ve been curious about. We’ll help you turn it into a plan you can stick with.</p>
        </div>
        <div className="empty-action">
          <button className="secondary-button" type="button" disabled>Set your first goal <span>↗</span></button>
          <small>Goal creation is coming in the next step</small>
        </div>
      </section>

      <footer><span>Small steps. Real momentum.</span><span className="api-indicator"><i className={apiStatus === 'Connected' ? 'online' : ''} /> API {apiStatus}</span></footer>
    </main>
  )
}

createRoot(document.getElementById('root')).render(<App />)
