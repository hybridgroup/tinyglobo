package main

import "testing"

func TestPointInPolygon_Latvia(t *testing.T) {
	// A point inside Latvia polygon
	lat, lon := float32(56.95), float32(24.1)
	if !pointInPolygon(lat, lon, latvia) {
		t.Errorf("Expected point (%f, %f) to be inside Latvia polygon", lat, lon)
	}

	// A point outside Latvia polygon
	lat, lon = float32(60.0), float32(25.0)
	if pointInPolygon(lat, lon, latvia) {
		t.Errorf("Expected point (%f, %f) to be outside Latvia polygon", lat, lon)
	}
}

func TestPointInPolygon_NorKor(t *testing.T) {
	// A point inside NorKor polygon
	lat, lon := float32(40.0), float32(126.0)
	if !pointInPolygon(lat, lon, norKor) {
		t.Errorf("Expected point (%f, %f) to be inside NorKor polygon", lat, lon)
	}

	// A point outside NorKor polygon
	lat, lon = float32(36.0), float32(130.0)
	if pointInPolygon(lat, lon, norKor) {
		t.Errorf("Expected point (%f, %f) to be outside NorKor polygon", lat, lon)
	}
}

func TestPointInPolygon_UK(t *testing.T) {
	// A point inside UK polygon
	lat, lon := float32(53.5), float32(-1.5)
	if !pointInPolygon(lat, lon, uk) {
		t.Errorf("Expected point (%f, %f) to be inside UK polygon", lat, lon)
	}

	// A point outside UK polygon
	lat, lon = float32(48.0), float32(-4.0)
	if pointInPolygon(lat, lon, uk) {
		t.Errorf("Expected point (%f, %f) to be outside UK polygon", lat, lon)
	}
}

func TestPointInPolygon_Yemen(t *testing.T) {
	// A point inside Yemen polygon
	lat, lon := float32(15.5), float32(48.0)
	if !pointInPolygon(lat, lon, yemen) {
		t.Errorf("Expected point (%f, %f) to be inside Yemen polygon", lat, lon)
	}

	// A point outside Yemen polygon
	lat, lon = float32(25.0), float32(50.0)
	if pointInPolygon(lat, lon, yemen) {
		t.Errorf("Expected point (%f, %f) to be outside Yemen polygon", lat, lon)
	}
}

func TestPointInPolygon_Ukraine(t *testing.T) {
	// A point inside Ukraine polygon
	lat, lon := float32(48.0), float32(30.0)
	if !pointInPolygon(lat, lon, ukraine) {
		t.Errorf("Expected point (%f, %f) to be inside Ukraine polygon", lat, lon)
	}

	// A point outside Ukraine polygon
	lat, lon = float32(60.0), float32(30.0)
	if pointInPolygon(lat, lon, ukraine) {
		t.Errorf("Expected point (%f, %f) to be outside Ukraine polygon", lat, lon)
	}
}

func TestGeofenced(t *testing.T) {
	// Save and restore globals
	oldLat, oldLon := currentLatitude, currentLongitude
	defer func() {
		currentLatitude = oldLat
		currentLongitude = oldLon
	}()

	// Inside Latvia
	currentLatitude = 56.95
	currentLongitude = 24.1
	if !geofenced() {
		t.Error("Expected geofenced() to return true for point inside Latvia")
	}

	// Outside all geofences
	currentLatitude = 0
	currentLongitude = 0
	if geofenced() {
		t.Error("Expected geofenced() to return false for (0,0)")
	}

	currentLatitude = 10.0
	currentLongitude = 10.0
	if geofenced() {
		t.Error("Expected geofenced() to return false for point outside all geofences")
	}
}
