export default function ProgressBar({ value = 0 }) {
  const pct = Math.min(100, Math.max(0, value))
  return <div className="progress"><span style={{ width: `${pct}%` }} /></div>
}