import type { SunResult } from './types'

interface DayNightBarProps {
  result: SunResult | null
}

const TOTAL_MINUTES = 1440

function DayNightBar({ result }: DayNightBarProps) {
  if (!result) return null

  const { polar_type, sunrise_minutes_local, sunset_minutes_local, sunrise_local, sunset_local } = result

  if (polar_type === 'midnight_sun') {
    return (
      <div className="day-night-bar">
        <div className="day-night-bar-track">
          <div className="bar-segment day" style={{ width: '100%' }} />
        </div>
        <div className="day-night-bar-label-center">24h daylight</div>
      </div>
    )
  }

  if (polar_type === 'polar_night') {
    return (
      <div className="day-night-bar">
        <div className="day-night-bar-track">
          <div className="bar-segment night" style={{ width: '100%' }} />
        </div>
        <div className="day-night-bar-label-center">24h darkness</div>
      </div>
    )
  }

  // Normal or wrapped case
  if (sunrise_minutes_local < 0 || sunset_minutes_local < 0 || sunrise_minutes_local > TOTAL_MINUTES || sunset_minutes_local > TOTAL_MINUTES) {
    return null
  }

  const isWrapped = sunrise_minutes_local > sunset_minutes_local

  const seg1: { type: 'night' | 'day'; width: number } = isWrapped
    ? { type: 'day', width: sunset_minutes_local / TOTAL_MINUTES * 100 }
    : { type: 'night', width: sunrise_minutes_local / TOTAL_MINUTES * 100 }

  const seg2: { type: 'night' | 'day'; width: number } = isWrapped
    ? { type: 'night', width: (sunrise_minutes_local - sunset_minutes_local) / TOTAL_MINUTES * 100 }
    : { type: 'day', width: (sunset_minutes_local - sunrise_minutes_local) / TOTAL_MINUTES * 100 }

  const seg3: { type: 'night' | 'day'; width: number } = isWrapped
    ? { type: 'day', width: (TOTAL_MINUTES - sunrise_minutes_local) / TOTAL_MINUTES * 100 }
    : { type: 'night', width: (TOTAL_MINUTES - sunset_minutes_local) / TOTAL_MINUTES * 100 }

  const label2 = isWrapped ? sunset_local : sunrise_local
  const label3 = isWrapped ? sunrise_local : sunset_local
  const label2Pos = seg1.width
  const label3Pos = seg1.width + seg2.width

  const ariaDescription = isWrapped
    ? `Daylight bar: ${sunset_local} to 24 then 00 to ${sunrise_local}`
    : `Daylight bar: ${sunrise_local} to ${sunset_local}`

  return (
    <div className="day-night-bar" aria-label={ariaDescription}>
      <div className="day-night-bar-track">
        <div className={`bar-segment ${seg1.type}`} style={{ width: seg1.width + '%' }} />
        <div className={`bar-segment ${seg2.type}`} style={{ width: seg2.width + '%' }} />
        <div className={`bar-segment ${seg3.type}`} style={{ width: seg3.width + '%' }} />
      </div>
      <div className="day-night-bar-labels">
        <span className="bar-label bar-label-start">00</span>
        <span className="bar-label" style={{ left: label2Pos + '%' }}>{label2}</span>
        <span className="bar-label" style={{ left: label3Pos + '%' }}>{label3}</span>
        <span className="bar-label bar-label-end">24</span>
      </div>
    </div>
  )
}

export default DayNightBar
