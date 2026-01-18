package main

// geofenced checks if the current GPS fix is within any of the defined geofences
func geofenced() bool {
	if currentLatitude == 0 && currentLongitude == 0 {
		return false
	}

	// check if inside any geofence
	if pointInPolygon(currentLatitude, currentLongitude, latvia) {
		return true
	}
	if pointInPolygon(currentLatitude, currentLongitude, norKor) {
		return true
	}
	if pointInPolygon(currentLatitude, currentLongitude, uk) {
		return true
	}
	if pointInPolygon(currentLatitude, currentLongitude, yemen) {
		return true
	}

	return false
}

var latvia = []float32{
	26.6418000, 55.6838000,
	28.2238770, 56.2280850,
	27.7349854, 57.4272096,
	26.4221191, 57.6395208,
	25.2465820, 58.0662560,
	21.6485596, 57.7598683,
	20.8850098, 56.9809114,
	20.7421875, 56.0904271,
	22.1044922, 56.4230166,
	25.2136230, 56.1975370,
	26.5869141, 55.6713893,
}

var norKor = []float32{
	123.7829590, 39.9434365,
	124.5410156, 37.4574181,
	129.0124512, 38.7112325,
	127.8259277, 39.4955634,
	131.1547852, 41.9921602,
	129.7595215, 43.2291951,
	123.7829590, 39.9518589,
}

var uk = []float32{
	-0.5932617, 61.0263703,
	-7.9541016, 58.2228110,
	-7.9650879, 54.1752967,
	-5.0427246, 53.8913913,
	-6.4819336, 49.7670741,
	1.4941406, 50.9791824,
	1.8457031, 52.6430634,
	-2.0983887, 56.3287209,
	-0.5932617, 61.0263703,
}

var yemen = []float32{
	52.1081543, 19.3111434,
	42.0556641, 17.3401517,
	43.4838867, 12.1360052,
	53.8330078, 15.6653542,
	52.1081543, 19.3111434,
}

// pointInPolygon determines if a point (latitude, longitude) is inside a polygon
// polygon is a flat array of coordinates: [lon0, lat0, lon1, lat1, ...]
func pointInPolygon(latitude, longitude float32, polygon []float32) bool {
	numCorners := len(polygon) / 2
	if numCorners < 3 {
		return false
	}

	j := numCorners*2 - 2
	oddNodes := false

	for i := 0; i < numCorners*2; i += 2 {
		if ((polygon[i+1] < latitude && polygon[j+1] >= latitude) ||
			(polygon[j+1] < latitude && polygon[i+1] >= latitude)) &&
			(polygon[i] <= longitude || polygon[j] <= longitude) {

			oddNodes = oddNodes != (polygon[i]+(latitude-polygon[i+1])/(polygon[j+1]-polygon[i+1])*(polygon[j]-polygon[i]) < longitude)
		}
		j = i
	}

	return oddNodes
}
