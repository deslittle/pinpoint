package convert

import "github.com/deslittle/pinpoint/pb"

func FromPbPolygonToGeoMultipolygon(pbpoly []*pb.Polygon) MultiPolygonCoordinates {
	res := MultiPolygonCoordinates{}
	for _, poly := range pbpoly {
		newGeoPoly := make(PolygonCoordinates, 0)

		mainpoly := [][2]float64{}
		for _, point := range poly.Points {
			mainpoly = append(mainpoly, [2]float64{float64(point.Lng), float64(point.Lat)})
		}
		newGeoPoly = append(newGeoPoly, mainpoly)

		for _, holepoly := range poly.Holes {
			holepolyCoords := [][2]float64{}
			for _, point := range holepoly.Points {
				holepolyCoords = append(holepolyCoords, [2]float64{float64(point.Lng), float64(point.Lat)})
			}
			newGeoPoly = append(newGeoPoly, holepolyCoords)
		}
		res = append(res, newGeoPoly)
	}
	return res
}

func RevertItem(input *pb.Location) *FeatureItem {
	return &FeatureItem{
		Type: FeatureType,
		Properties: PropertiesDefine{
			Name: input.Name,
		},
		Geometry: GeometryDefine{
			Type:        MultiPolygonType,
			Coordinates: FromPbPolygonToGeoMultipolygon(input.Polygons),
		},
	}
}

// Revert could convert pb define data to GeoJSON format.
func Revert(input *pb.Locations) *BoundaryFile {
	output := &BoundaryFile{}
	for _, location := range input.Locations {
		item := RevertItem(location)
		output.Features = append(output.Features, item)
	}
	output.Type = "FeatureCollection"
	return output
}
