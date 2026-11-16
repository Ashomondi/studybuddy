export default function Field({ label, error, children, hint }) {
  return (
    <div className="field">
      {label && <label>{label}</label>}
      {children}
      {error && <span className="error">{error}</span>}
      {hint && !error && <span style={{ fontSize: 12, color: 'var(--color-muted)' }}>{hint}</span>}
    </div>
  )
}