package convert

import (
	"github.com/deslittle/pinpoint/pb"
	"github.com/tidwall/geojson/geometry"
)

func FromLocationPBToGeometryPoly(location *pb.Location) []*geometry.Poly {
	ret := []*geometry.Poly{}
	for _, polygon := range location.Polygons {

		newPoints := make([]geometry.Point, 0)
		for _, point := range polygon.Points {
			newPoints = append(newPoints, geometry.Point{
				X: float64(point.Lng),
				Y: float64(point.Lat),
			})
		}

		holes := [][]geometry.Point{}
		for _, holePoly := range polygon.Holes {
			newHolePoints := make([]geometry.Point, 0)
			for _, point := range holePoly.Points {
				newHolePoints = append(newHolePoints, geometry.Point{
					X: float64(point.Lng),
					Y: float64(point.Lat),
				})
			}
			holes = append(holes, newHolePoints)
		}

		newPoly := geometry.NewPoly(newPoints, holes, nil)
		ret = append(ret, newPoly)
	}
	return ret
}
