export default function Button({ children, variant = 'primary', size, className = '', ...props }) {
  const cls = `btn btn-${variant}${size === 'sm' ? ' btn-sm' : ''} ${className}`.trim()
  return <button className={cls} {...props}>{children}</button>
}