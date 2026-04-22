"use client";

import { useEffect, useRef, useState } from "react";
import L from "leaflet";
import "leaflet-draw";
import { MapDebugPanel } from "./MapDebugPanel";

// Import CSS explicitly in case globals.css isn't loaded yet
import "leaflet/dist/leaflet.css";
import "leaflet-draw/dist/leaflet.draw.css";

// Fix default marker icons
if (typeof window !== "undefined") {
  delete (L.Icon.Default.prototype as any)._getIconUrl;
  L.Icon.Default.mergeOptions({
    iconRetinaUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png",
    iconUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png",
    shadowUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png",
  });
}

interface TerritoryMapDrawerProps {
  initialPolygon?: number[][][]; // GeoJSON polygon coordinates
  color?: string;
  onPolygonChange: (polygon: number[][][]) => void;
  center?: [number, number];
  zoom?: number;
}

export function TerritoryMapDrawer({
  initialPolygon,
  color = "#3B82F6",
  onPolygonChange,
  center = [-6.2, 106.816666], // Jakarta default
  zoom = 12,
}: TerritoryMapDrawerProps) {
  const mapRef = useRef<L.Map | null>(null);
  const drawnItemsRef = useRef<L.FeatureGroup | null>(null);
  const mapIdRef = useRef<string>(`territory-map-${Math.random().toString(36).substr(2, 9)}`);
  const [isClient, setIsClient] = useState(false);
  const initialPolygonLoadedRef = useRef<string>("");

  useEffect(() => {
    setIsClient(true);
  }, []);

  useEffect(() => {
    if (!isClient) return;

    // Wait for DOM element to be ready
    const mapElement = document.getElementById(mapIdRef.current);
    if (!mapElement) {
      // Retry after a short delay
      const timeout = setTimeout(() => {
        const retryElement = document.getElementById(mapIdRef.current);
        if (retryElement && !mapRef.current) {
          // Force re-render
          setIsClient(false);
          setTimeout(() => setIsClient(true), 50);
        }
      }, 50);
      return () => clearTimeout(timeout);
    }

    // Initialize map only once
    if (!mapRef.current) {
      // Small delay to ensure DOM is fully ready and Leaflet is loaded
      const initTimeout = setTimeout(() => {
        if (mapRef.current) return; // Already initialized
        
        const mapElement = document.getElementById(mapIdRef.current);
        if (!mapElement) {
          return;
        }

        const map = L.map(mapIdRef.current, {
          preferCanvas: false,
        }).setView(center, zoom);

        // Add tile layer
        L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
          attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
          maxZoom: 19,
        }).addTo(map);

        // Initialize FeatureGroup for drawn items
        const drawnItems = new L.FeatureGroup();
        // Add to map first - this is critical for visibility
        map.addLayer(drawnItems);
        drawnItemsRef.current = drawnItems;
        
        // Verify FeatureGroup is on map and force it to be visible
        
        // Ensure the FeatureGroup pane exists and is properly styled
        try {
          const pane = map.getPane('overlayPane');
          if (pane) {
            pane.style.zIndex = '400'; // Ensure drawn items are visible
          }
        } catch (err) {
        }

        // Check if Leaflet Draw is available
        if (!L.Control.Draw) {
          // Show user-friendly error
          const errorDiv = document.createElement('div');
          errorDiv.innerHTML = `
            <div style="position: absolute; top: 10px; right: 10px; background: #f8d7da; color: #721c24; padding: 8px; border-radius: 4px; font-size: 12px; z-index: 1000;">
              ⚠️ Drawing tools not available. Please refresh the page.
            </div>
          `;
          mapElement.appendChild(errorDiv);
          return;
        }

        // Add draw control with proper styling
        const drawControl = new L.Control.Draw({
          position: "topright",
          draw: {
            polygon: {
              shapeOptions: {
                color: color,
                fillColor: color,
                fillOpacity: 0.3,
                weight: 2,
                opacity: 0.8,
              },
              allowIntersection: false,
              drawError: {
                color: "#e74c3c",
                message: "<strong>Error:</strong> Shape edges cannot cross!",
              },
            },
            polyline: false,
            circle: false,
            rectangle: {
              shapeOptions: {
                color: color,
                fillColor: color,
                fillOpacity: 0.3,
                weight: 2,
                opacity: 0.8,
              },
            },
            marker: false,
            circlemarker: false,
          },
          edit: {
            featureGroup: drawnItems,
            remove: true,
          },
        });
        map.addControl(drawControl);
        
        // IMPORTANT: Override the default draw behavior to ensure layers go to our FeatureGroup
        // This is a critical fix for layer visibility
        const originalOnAdd = (L.Draw.Polygon.prototype as any)._onTouch || (L.Draw.Polygon.prototype as any)._onClick;
        if (originalOnAdd) {
        }

        // Handle polygon creation - use string event names for compatibility
        map.on("draw:created", (e: any) => {
          const layer = e.layer;
          
          
          // Clear existing polygons first (only one polygon allowed)
          drawnItems.clearLayers();
          
          // Apply proper styling to the layer BEFORE adding
          layer.setStyle({
            color: color,
            fillColor: color,
            fillOpacity: 0.3,
            weight: 3,
            opacity: 0.9,
          });
          
          // IMPORTANT: Add layer to FeatureGroup which should make it visible
          // Since FeatureGroup is already on the map, adding to it should display the layer
          drawnItems.addLayer(layer);
          
          // Force layer to be visible by also adding directly to map as backup
          if (!map.hasLayer(layer)) {
            map.addLayer(layer);
          }
          
          
          // Force immediate redraw
          map.invalidateSize(false);
          
          // Extract coordinates for parent component
          const latlngs = layer.getLatLngs()[0];
          if (!latlngs || latlngs.length < 3) {
            return;
          }

          // Convert coordinates to GeoJSON format
          const coordinates = latlngs.map((latlng: L.LatLng) => [
            latlng.lng,
            latlng.lat,
          ]);
          
          // Close the polygon by adding first point at the end
          if (coordinates.length > 0 && 
              (coordinates[0][0] !== coordinates[coordinates.length - 1][0] || 
               coordinates[0][1] !== coordinates[coordinates.length - 1][1])) {
            coordinates.push(coordinates[0]);
          }
          
          // Validate polygon format for GeoJSON
          if (coordinates.length < 4) {
            return;
          }
          
          // Use timeout to ensure visibility after coordinates are processed
          setTimeout(() => {
            if (mapRef.current && drawnItemsRef.current) {
              try {
                
                // Force layer visibility one more time
                if (drawnItemsRef.current.hasLayer(layer)) {
                  layer.bringToFront();
                  
                  // Fit bounds to show the polygon
                  const bounds = layer.getBounds();
                  if (bounds && bounds.isValid && bounds.isValid()) {
                    mapRef.current.fitBounds(bounds, { padding: [20, 20] });
                  }
                }
                
                // Force complete map redraw
                mapRef.current.invalidateSize(false);
              } catch (err) {
              }
            }
          }, 100);
          
          // Update parent component
          onPolygonChange([coordinates]);
        });

        // Handle polygon edit
        map.on("draw:edited", (e: any) => {
          const layers = e.layers;
          layers.eachLayer((layer: any) => {
            // Apply styling after edit
            layer.setStyle({
              color: color,
              fillColor: color,
              fillOpacity: 0.3,
              weight: 3,
              opacity: 0.9,
            });
            
            const latlngs = layer.getLatLngs()[0];
            if (!latlngs || latlngs.length < 3) {
              return;
            }
            const coordinates = latlngs.map((latlng: L.LatLng) => [
              latlng.lng,
              latlng.lat,
            ]);
            if (coordinates.length > 0 && 
                (coordinates[0][0] !== coordinates[coordinates.length - 1][0] || 
                 coordinates[0][1] !== coordinates[coordinates.length - 1][1])) {
              coordinates.push(coordinates[0]);
            }
            onPolygonChange([coordinates]);
          });
        });

        // Handle polygon delete
        map.on("draw:deleted", () => {
          onPolygonChange([]);
        });

        mapRef.current = map;
      }, 100); // Increased delay to ensure everything is ready

      return () => {
        clearTimeout(initTimeout);
        // Better cleanup
        if (mapRef.current) {
          try {
            // Remove all event listeners
            mapRef.current.off();
            // Clear drawn items
            if (drawnItemsRef.current) {
              drawnItemsRef.current.clearLayers();
            }
            // Remove map
            mapRef.current.remove();
          } catch (err) {
          } finally {
            mapRef.current = null;
            drawnItemsRef.current = null;
          }
        }
      };
    }

    // Cleanup
    return () => {
      if (mapRef.current) {
        try {
          mapRef.current.off();
          if (drawnItemsRef.current) {
            drawnItemsRef.current.clearLayers();
          }
          mapRef.current.remove();
        } catch (err) {
        } finally {
          mapRef.current = null;
          drawnItemsRef.current = null;
        }
      }
    };
  }, [isClient, color, center, zoom]);

  // Separate effect to handle initialPolygon changes
  useEffect(() => {
    if (!isClient || !mapRef.current || !drawnItemsRef.current) return;


    // Create a unique key for the polygon to detect changes
    const polygonKey = JSON.stringify(initialPolygon);
    
    // Only reload if polygon actually changed
    if (polygonKey === initialPolygonLoadedRef.current) return;

    // Clear existing layers
    drawnItemsRef.current.clearLayers();
    initialPolygonLoadedRef.current = polygonKey;

    // Load initial polygon if provided
    if (initialPolygon && initialPolygon.length > 0 && initialPolygon[0] && initialPolygon[0].length > 0) {
      try {
        
        // Convert coordinates: GeoJSON format is [lng, lat], Leaflet needs [lat, lng]
        const latlngs = initialPolygon[0]
          .map((coord, index) => {
            if (Array.isArray(coord) && coord.length >= 2) {
              const lng = coord[0];
              const lat = coord[1];
              
              // Validate coordinates
              if (typeof lng === 'number' && typeof lat === 'number' && 
                  lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180) {
                return L.latLng(lat, lng); // Leaflet uses [lat, lng]
              } else {
              }
            }
            return null;
          })
          .filter((latlng): latlng is L.LatLng => latlng !== null);

        if (latlngs.length >= 3) {
          
          const polygon = L.polygon(latlngs, {
            color: color,
            fillColor: color,
            fillOpacity: 0.3,
            weight: 3,
            opacity: 0.9,
          });
          
          // Add to FeatureGroup first
          drawnItemsRef.current.addLayer(polygon);
          
          // Also add directly to map to ensure visibility
          if (!mapRef.current.hasLayer(polygon)) {
            mapRef.current.addLayer(polygon);
          }
          
          
          // Use defensive approach for visibility
          const ensureInitialVisibility = () => {
            if (!mapRef.current || !drawnItemsRef.current) {
              return;
            }
            
            try {
              // Check if map container still exists
              const container = mapRef.current.getContainer();
              if (!container) {
                return;
              }
              
              
              // Bring polygon to front and fit bounds
              if (drawnItemsRef.current.hasLayer(polygon) || mapRef.current.hasLayer(polygon)) {
                polygon.bringToFront();
                
                // Force map to redraw and fit bounds
                mapRef.current.invalidateSize(false);
                
                const bounds = polygon.getBounds();
                if (bounds && bounds.isValid()) {
                  mapRef.current.fitBounds(bounds, { padding: [20, 20] });
                }
              }
            } catch (err) {
            }
          };

          // Use shorter timeout
          setTimeout(ensureInitialVisibility, 100);
        } else {
        }
      } catch (error) {
      }
    } else {
    }
  }, [isClient, initialPolygon, color]);

  if (!isClient) {
    return (
      <div className="h-[400px] w-full bg-muted rounded-md flex items-center justify-center">
        <p className="text-muted-foreground">Loading map...</p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <div id={mapIdRef.current} className="h-[400px] w-full rounded-md border" />
      <div className="text-xs text-muted-foreground space-y-1">
        <p>
          Use the drawing tools on the right to draw a polygon for this territory. 
          Click points to create a shape, or use the rectangle tool for quick setup.
        </p>
      </div>
      <MapDebugPanel 
        mapRef={mapRef} 
        drawnItemsRef={drawnItemsRef} 
        polygon={initialPolygon}
      />
    </div>
  );
}
