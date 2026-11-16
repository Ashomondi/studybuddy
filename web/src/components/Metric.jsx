export default function Metric({ value, label }) {
  return <div className="metric"><span className="value">{value}</span><span className="label">{label}</span></div>
}