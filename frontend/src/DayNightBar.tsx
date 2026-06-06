import type { SunResult } from './types'

interface DayNightBarProps {
  result: SunResult | null
}

const TOTAL_MINUTES = 1440

/**
 * Computes a linear gradient for a day segment so that it is brightest
 * at mid-day and dims toward the sunrise/sunset edges.
 *
 * @param segStart - Segment start in minutes (0-1440)
 * @param segEnd - Segment end in minutes (0-1440)
 * @param peak - Mid-day position in minutes
 * @param dim - Darker color at edges
 * @param bright - Brightest color at peak
 */
function computeDayGradient(segStart: number, segEnd: number, peak: number, dim: string, bright: string): string {
  const segLen = segEnd - segStart
  // Position of peak relative to segment (0 = start, 1 = end)
  const peakRel = (peak - segStart) / segLen

  if (peakRel <= 0) {
    // Peak is before segment — fade from bright to dim
    return `linear-gradient(to right, ${bright}, ${dim})`
  }
  if (peakRel >= 1) {
    // Peak is after segment — fade from dim to bright
    return `linear-gradient(to right, ${dim}, ${bright})`
  }

  // Peak is inside segment — bright at peak, dim at both edges
  const pct = Math.round(peakRel * 100)
  return `linear-gradient(to right, ${dim}, ${bright} ${pct}%, ${dim})`
}

function DayNightBar({ result }: DayNightBarProps) {
  if (!result) return null

  const { polar_type, sunrise_minutes_local, sunset_minutes_local, sunrise_local, sunset_local } = result

  if (polar_type === 'midnight_sun') {
    return (
      <div className="day-night-bar">
        <div className="day-night-bar-track">
          <div className="bar-segment day" style={{ width: '100%' }} />
        </div>
        <div className="day-night-bar-icons">
          <span className="bar-icon" style={{ left: '50%' }}>☀️</span>
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
        <div className="day-night-bar-icons">
          <span className="bar-icon" style={{ left: '50%' }}>🌙</span>
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

  // Mid-day: center of the daylight period
  // In wrapped case the naive average lands in the night gap, so shift by 12h
  const midDayMinutes = isWrapped
    ? ((sunrise_minutes_local + sunset_minutes_local) / 2 + TOTAL_MINUTES / 2) % TOTAL_MINUTES
    : (sunrise_minutes_local + sunset_minutes_local) / 2

  // Gradient colors for day segments (dim at edges, bright at mid-day)
  const dayDim = '#e08a15'
  const dayBright = '#ffd080'

  // Build segments (skip zero-width)
  const segments: Array<{ type: 'night' | 'day'; width: number; gradient?: string }> = []

  if (isWrapped) {
    // day (0-sunset), night (sunset-sunrise), day (sunrise-24)
    const w1 = sunset_minutes_local / TOTAL_MINUTES * 100
    const w2 = (sunrise_minutes_local - sunset_minutes_local) / TOTAL_MINUTES * 100
    const w3 = (TOTAL_MINUTES - sunrise_minutes_local) / TOTAL_MINUTES * 100
    if (w1 > 0) segments.push({ type: 'day', width: w1, gradient: computeDayGradient(0, sunset_minutes_local, midDayMinutes, dayDim, dayBright) })
    if (w2 > 0) segments.push({ type: 'night', width: w2 })
    // Second day segment: add TOTAL_MINUTES to peak so it falls "after" the segment
    // in linear coordinate space (handles circular time correctly)
    if (w3 > 0) segments.push({ type: 'day', width: w3, gradient: computeDayGradient(sunrise_minutes_local, TOTAL_MINUTES, midDayMinutes + TOTAL_MINUTES, dayDim, dayBright) })
  } else {
    // night (0-sunrise), day (sunrise-sunset), night (sunset-24)
    const w1 = sunrise_minutes_local / TOTAL_MINUTES * 100
    const w2 = (sunset_minutes_local - sunrise_minutes_local) / TOTAL_MINUTES * 100
    const w3 = (TOTAL_MINUTES - sunset_minutes_local) / TOTAL_MINUTES * 100
    if (w1 > 0) segments.push({ type: 'night', width: w1 })
    if (w2 > 0) segments.push({ type: 'day', width: w2, gradient: computeDayGradient(sunrise_minutes_local, sunset_minutes_local, midDayMinutes, dayDim, dayBright) })
    if (w3 > 0) segments.push({ type: 'night', width: w3 })
  }

  // Transitions: positions, labels, icons, and per-transition visibility
  const t1Minutes = isWrapped ? sunset_minutes_local : sunrise_minutes_local
  const t2Minutes = isWrapped ? sunrise_minutes_local : sunset_minutes_local
  const t1Pos = t1Minutes / TOTAL_MINUTES * 100
  const t2Pos = t2Minutes / TOTAL_MINUTES * 100
  const t1Label = isWrapped ? sunset_local : sunrise_local
  const t2Label = isWrapped ? sunrise_local : sunset_local
  const t1Icon = isWrapped ? '🌇' : '🌅'
  const t2Icon = isWrapped ? '🌅' : '🌇'

  // Omit marker/icon/label if transition is exactly at 0 or 1440 (would overlap edge labels)
  const showT1 = t1Minutes !== 0 && t1Minutes !== TOTAL_MINUTES
  const showT2 = t2Minutes !== 0 && t2Minutes !== TOTAL_MINUTES

  const ariaDescription = isWrapped
    ? `Daylight bar: ${sunset_local} to 24 then 00 to ${sunrise_local}`
    : `Daylight bar: ${sunrise_local} to ${sunset_local}`

  return (
    <div className="day-night-bar" aria-label={ariaDescription}>
      <div className="day-night-bar-track">
        {segments.map((seg, i) => (
          <div
            key={`${seg.type}-${i}`}
            className={`bar-segment ${seg.type}${i === 0 ? ' bar-segment-first' : ''}${i === segments.length - 1 ? ' bar-segment-last' : ''}`}
            style={{ width: seg.width + '%', ...(seg.gradient ? { background: seg.gradient } : {}) }}
          />
        ))}
        {showT1 && (
          <div className="bar-marker" style={{ left: t1Pos + '%' }} />
        )}
        {showT2 && (
          <div className="bar-marker" style={{ left: t2Pos + '%' }} />
        )}
      </div>
      <div className="day-night-bar-icons">
        {showT1 && (
          <span className="bar-icon" style={{ left: t1Pos + '%' }}>{t1Icon}</span>
        )}
        {showT2 && (
          <span className="bar-icon" style={{ left: t2Pos + '%' }}>{t2Icon}</span>
        )}
      </div>
      <div className="day-night-bar-labels">
        <span className="bar-label bar-label-start">00</span>
        {showT1 && (
          <span className="bar-label" style={{ left: t1Pos + '%' }}>{t1Label}</span>
        )}
        {showT2 && (
          <span className="bar-label" style={{ left: t2Pos + '%' }}>{t2Label}</span>
        )}
        <span className="bar-label bar-label-end">24</span>
      </div>
    </div>
  )
}

export default DayNightBar
