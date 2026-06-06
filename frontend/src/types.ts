export interface SunResult {
  sunrise_local: string
  sunset_local: string
  day_length: string
  timezone: string
  sunrise_minutes_local: number
  sunset_minutes_local: number
  polar_type: '' | 'midnight_sun' | 'polar_night'
}
