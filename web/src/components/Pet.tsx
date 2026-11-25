import type { Mood } from '../health'

// Face features sit on the pink body in both themes, so they keep the light-theme ink.
const ink = '#2b2d42'

function Eyes({ mood, still }: { mood: Mood; still?: boolean }) {
  if (mood === 'sleepy') {
    return (
      <g fill="none" stroke={ink} strokeWidth="4" strokeLinecap="round">
        <path d="M50 73q8 6 16 0" />
        <path d="M84 73q8 6 16 0" />
      </g>
    )
  }
  if (mood === 'dizzy') {
    return (
      <g fill="none" stroke={ink} strokeWidth="3.5" strokeLinecap="round">
        <path d="M58 72a3 3 0 1 1 3 3 6 6 0 1 1 -6 -6 9 9 0 0 1 9 9" />
        <path d="M92 72a3 3 0 1 0 -3 3 6 6 0 1 0 6 -6 9 9 0 0 0 -9 9" />
      </g>
    )
  }
  return (
    <g className={still ? undefined : 'pet-eyes'}>
      <ellipse cx="58" cy="72" rx="7" ry="9" fill={ink} />
      <ellipse cx="92" cy="72" rx="7" ry="9" fill={ink} />
      <circle cx="60.5" cy="68.5" r="2.6" fill="#fff" />
      <circle cx="94.5" cy="68.5" r="2.6" fill="#fff" />
    </g>
  )
}

function Mouth({ mood }: { mood: Mood }) {
  switch (mood) {
    case 'happy':
      return <path d="M63 92q12 16 24 0z" fill={ink} stroke={ink} strokeWidth="3" strokeLinejoin="round" />
    case 'worried':
      return <path d="M63 98q4-4 8 0t8 0 8 0" fill="none" stroke={ink} strokeWidth="3.5" strokeLinecap="round" />
    case 'sad':
      return <path d="M64 100q11-10 22 0" fill="none" stroke={ink} strokeWidth="4" strokeLinecap="round" />
    case 'sleepy':
      return <ellipse cx="75" cy="97" rx="4" ry="3" fill={ink} />
    case 'dizzy':
      return <ellipse cx="75" cy="98" rx="6" ry="7" fill={ink} />
  }
}

function Extra({ mood }: { mood: Mood }) {
  if (mood === 'sad') return <path d="M100 82q5 11 0 16q-5-5 0-16z" fill="var(--tear)" />
  if (mood === 'worried') return <path d="M108 44q6 9 1 13q-6-3-1-13z" fill="var(--tear)" />
  if (mood === 'sleepy')
    return (
      <g className="pet-zzz" fill="none" stroke="var(--sky)" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round">
        <path d="M112 34h9l-9 10h9" />
      </g>
    )
  return null
}

type PetProps = {
  mood: Mood
  size?: number
  /** Accessible description. Omit for a decorative pet. */
  label?: string
  /** No idle motion: for the brand mark, so only one pet is alive on screen. */
  still?: boolean
  /** CSS size override, e.g. a smaller pet on phones. */
  className?: string
}

/** The original status monster. Its face is the system's state. */
export function Pet({ mood, size = 150, label, still, className = '' }: PetProps) {
  const a11y = label ? { role: 'img', 'aria-label': label } : { 'aria-hidden': true }
  return (
    <svg width={size} height={size} viewBox="0 0 150 150" className={`shrink-0 ${className}`} {...a11y}>
      <ellipse cx="75" cy="139" rx="46" ry="6" fill="rgb(var(--shade) / 0.12)" />
      <g className={still ? undefined : 'pet-body'}>
        <path d="M44 40 36 14l24 18zm62 0 8-26-24 18z" fill="var(--pink)" />
        <path d="M30 120C14 80 28 30 75 28s61 52 45 92c-10 14-80 14-90 0z" fill="var(--pink)" />
        <path d="M31 118c4 9 84 9 88 0-3 8-9 12-14 13-20 4-46 4-60 0-6-2-11-6-14-13z" fill="#f47396" opacity=".55" />
        <ellipse cx="75" cy="123.5" rx="16" ry="5.5" fill="#ffb8ca" />
        <ellipse cx="52" cy="46" rx="14" ry="7" fill="#fff" opacity=".4" transform="rotate(-28 52 46)" />
        <ellipse cx="42" cy="88" rx="8" ry="5" fill="#ff6b93" opacity=".55" />
        <ellipse cx="108" cy="88" rx="8" ry="5" fill="#ff6b93" opacity=".55" />
        <Eyes mood={mood} still={still} />
        <Mouth mood={mood} />
        <Extra mood={mood} />
      </g>
    </svg>
  )
}
