import Icon from './Icon'

export default function Topbar({ title, apiStatus }) {
  return (
    <header className="topbar">
      <span className="topbar-title">{title}</span>
      <div className="topbar-actions">
        <span><i className={`api-dot ${apiStatus === 'Connected' ? 'online' : ''}`} /> API {apiStatus}</span>
      </div>
    </header>
  )
}