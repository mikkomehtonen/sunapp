import { useState } from 'react'
import './App.css'

interface SunResult {
  sunrise_utc: string
  sunset_utc: string
  sunrise_local: string
  sunset_local: string
  day_length: string
  timezone: string
}

function App() {
  const [lat, setLat] = useState('51.5074')
  const [lon, setLon] = useState('-0.1278')
  const [date, setDate] = useState(() => {
    const today = new Date()
    return today.toISOString().split('T')[0]
  })
  const [tz, setTz] = useState(() => Intl.DateTimeFormat().resolvedOptions().timeZone)
  const [result, setResult] = useState<SunResult | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setResult(null)

    const latNum = parseFloat(lat)
    const lonNum = parseFloat(lon)

    if (isNaN(latNum) || isNaN(lonNum)) {
      setError('Latitude and longitude must be valid numbers')
      return
    }

    if (latNum < -90 || latNum > 90) {
      setError('Latitude must be between -90 and 90')
      return
    }

    if (lonNum < -180 || lonNum > 180) {
      setError('Longitude must be between -180 and 180')
      return
    }

    setLoading(true)

    try {
      const params = new URLSearchParams({
        lat: lat.toString(),
        lon: lon.toString(),
        date: date,
        tz,
      })

      const response = await fetch(`/api/sun?${params}`)

      if (!response.ok) {
        const text = await response.text()
        throw new Error(text || 'Failed to fetch sun data')
      }

      const data = await response.json()
      setResult(data)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'An unknown error occurred')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="app">
      <header>
        <h1>Sunrise &amp; Sunset Calculator</h1>
        <p>Enter a location and date to see sunrise, sunset, and day length.</p>
      </header>

      <main>
        <form onSubmit={handleSubmit} className="sun-form">
          <div className="form-grid">
            <div className="form-group">
              <label htmlFor="lat">Latitude</label>
              <input
                id="lat"
                type="number"
                step="any"
                value={lat}
                onChange={(e) => setLat(e.target.value)}
                placeholder="e.g. 51.5074"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="lon">Longitude</label>
              <input
                id="lon"
                type="number"
                step="any"
                value={lon}
                onChange={(e) => setLon(e.target.value)}
                placeholder="e.g. -0.1278"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="date">Date</label>
              <input
                id="date"
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="tz">Timezone</label>
              <input
                id="tz"
                type="text"
                value={tz}
                onChange={(e) => setTz(e.target.value)}
                placeholder="e.g. Europe/Helsinki"
              />
            </div>
          </div>

          <button type="submit" disabled={loading} className="submit-btn">
            {loading ? 'Calculating...' : 'Calculate'}
          </button>
        </form>

        {error && <div className="error-message">{error}</div>}

        {result && (
          <div className="results">
            <h2>Results ({result.timezone})</h2>
            <div className="result-cards">
              <div className="card">
                <h3>Sunrise (UTC)</h3>
                <p className="time">{result.sunrise_utc}</p>
              </div>
              <div className="card">
                <h3>Sunset (UTC)</h3>
                <p className="time">{result.sunset_utc}</p>
              </div>
              <div className="card">
                <h3>Sunrise ({result.timezone})</h3>
                <p className="time">{result.sunrise_local}</p>
              </div>
              <div className="card">
                <h3>Sunset ({result.timezone})</h3>
                <p className="time">{result.sunset_local}</p>
              </div>
              <div className="card">
                <h3>Day Length</h3>
                <p className="time">{result.day_length}</p>
              </div>
            </div>
          </div>
        )}
      </main>
    </div>
  )
}

export default App
