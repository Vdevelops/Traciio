"use client";

import { useEffect } from "react";
import { useMap } from "react-leaflet";
import L from "leaflet";
import "leaflet-draw";

interface EditControlProps {
  position?: L.ControlPosition;
  onCreated: (e: any) => void;
  onEdited?: (e: any) => void;
  onDeleted?: (e: any) => void;
  draw?: any; // Use any for draw options to avoid complex typing
  edit?: {
    featureGroup: L.FeatureGroup;
    remove?: boolean;
  };
}

export function EditControl({
  position = "topright",
  onCreated,
  onEdited,
  onDeleted,
  draw = {
    polygon: {
      shapeOptions: {
        color: "#3B82F6",
        fillColor: "#3B82F6",
        fillOpacity: 0.3,
        weight: 3,
        opacity: 0.9,
      },
      allowIntersection: false,
      drawError: {
        color: "#e74c3c",
        message: "<strong>Error:</strong> Shape edges cannot cross!",
      },
    },
    rectangle: {
      shapeOptions: {
        color: "#3B82F6",
        fillColor: "#3B82F6",
        fillOpacity: 0.3,
        weight: 3,
        opacity: 0.9,
      },
    },
    polyline: false,
    circle: false,
    marker: false,
    circlemarker: false,
  },
  edit,
}: EditControlProps) {
  const map = useMap();

  useEffect(() => {
    if (!map || !L.Control.Draw) return;

    const drawControl = new L.Control.Draw({
      position,
      draw,
      edit,
    });

    map.addControl(drawControl);

    // Event handlers
    map.on(L.Draw.Event.CREATED, onCreated);
    if (onEdited) {
      map.on(L.Draw.Event.EDITED, onEdited);
    }
    if (onDeleted) {
      map.on(L.Draw.Event.DELETED, onDeleted);
    }

    return () => {
      map.removeControl(drawControl);
      map.off(L.Draw.Event.CREATED, onCreated);
      if (onEdited) {
        map.off(L.Draw.Event.EDITED, onEdited);
      }
      if (onDeleted) {
        map.off(L.Draw.Event.DELETED, onDeleted);
      }
    };
  }, [map, onCreated, onEdited, onDeleted, position, draw, edit]);

  return null;
}