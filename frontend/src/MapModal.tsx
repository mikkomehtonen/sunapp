import { useState } from 'react'
import { MapContainer, TileLayer, Marker, useMapEvents } from 'react-leaflet'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'

delete (L.Icon.Default.prototype as unknown as Record<string, unknown>)._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
})

interface MapModalProps {
  onClose: () => void
  onSelect: (lat: number, lon: number) => void
  initialLat?: number
  initialLon?: number
}

function MapClickHandler({ onSelect }: { onSelect: (lat: number, lon: number) => void }) {
  useMapEvents({
    click(e) {
      onSelect(e.latlng.lat, e.latlng.lng)
    },
  })
  return null
}

function getInitialPosition(initialLat?: number, initialLon?: number): [number, number] {
  if (initialLat !== undefined && initialLon !== undefined && !isNaN(initialLat) && !isNaN(initialLon)) {
    return [initialLat, initialLon]
  }
  return [61.5, 23.75]
}

function MapModal({ onClose, onSelect, initialLat, initialLon }: MapModalProps) {
  const [position, setPosition] = useState<[number, number]>(() => getInitialPosition(initialLat, initialLon))

  const handleSelect = (lat: number, lon: number) => {
    setPosition([lat, lon])
    onSelect(lat, lon)
  }

  return (
    <div className="map-modal-overlay" onClick={onClose}>
      <div className="map-modal" onClick={(e) => e.stopPropagation()}>
        <button className="map-modal-close" onClick={onClose} aria-label="Close">
          ×
        </button>
        <h3 className="map-modal-title">Click on the map to select a location</h3>
        <div className="map-container">
          <MapContainer
            center={position}
            zoom={4}
            style={{ height: '100%', width: '100%' }}
          >
            <TileLayer
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            <MapClickHandler onSelect={handleSelect} />
            <Marker position={position} />
          </MapContainer>
        </div>
        <div className="map-coords">
          {position[0].toFixed(4)}, {position[1].toFixed(4)}
        </div>
      </div>
    </div>
  )
}

export default MapModal
