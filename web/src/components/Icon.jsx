const paths = {
  home: 'M3 10.5 12 3l9 7.5V20a1 1 0 0 1-1 1h-5v-6H9v6H4a1 1 0 0 1-1-1v-9.5z',
  target: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20zm0 4a6 6 0 1 1 0 12 6 6 0 0 1 0-12zm0 4a2 2 0 1 0 0 4 2 2 0 0 0 0-4z',
  plan: 'M4 4h16v16H4zM8 8h8M8 12h8M8 16h5',
  clock: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20zm0 5v5l3.5 2',
  plus: 'M12 5v14M5 12h14',
  trash: 'M3 6h18M8 6V4h8v2M6 6l1 14h10l1-14',
  edit: 'M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04a1 1 0 0 0 0-1.41l-2.34-2.34a1 1 0 0 0-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z',
  check: 'M20 6 9 17l-5-5',
  x: 'M18 6 6 18M6 6l12 12',
  chevron: 'M9 6l6 6-6 6',
  calendar: 'M3 6h18v15H3zM3 10h18M8 3v4M16 3v4',
  book: 'M4 5a2 2 0 0 1 2-2h6v18H6a2 2 0 0 1-2-2V5zM20 5a2 2 0 0 0-2-2h-6v18h6a2 2 0 0 0 2-2V5z',
}

export default function Icon({ name, size = 18, className }) {
  const d = paths[name]
  if (!d) return null
  return (
    <svg className={className} width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <path d={d} />
    </svg>
  )
}