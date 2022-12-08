package reduce

import (
	"fmt"

	"github.com/deslittle/pinpoint/pb"
	"github.com/twpayne/go-polyline"
)

func CompressedPointsToPolylineBytes(points []*pb.Point) []byte {
	expect := [][]float64{}
	for _, point := range points {
		expect = append(expect, []float64{float64(point.Lng), float64(point.Lat)})
	}
	return polyline.EncodeCoords(expect)
}

func DecompressedPolylineBytesToPoints(input []byte) []*pb.Point {
	expect := []*pb.Point{}
	coords, _, _ := polyline.DecodeCoords(input)
	for _, coord := range coords {
		expect = append(expect, &pb.Point{
			Lng: float32(coord[0]), Lat: float32(coord[1]),
		})
	}
	return expect
}

func CompressWithPolyline(input *pb.Locations) *pb.CompressedLocations {
	output := &pb.CompressedLocations{
		Method: pb.CompressMethod_Polyline,
	}
	for _, location := range input.Locations {
		reducedLocation := &pb.CompressedLocation{
			Name: location.Name,
		}
		for _, polygon := range location.Polygons {
			newPoly := &pb.CompressedPolygon{
				Points: CompressedPointsToPolylineBytes(polygon.Points),
				Holes:  make([]*pb.CompressedPolygon, 0),
			}
			for _, hole := range polygon.Holes {
				newPoly.Holes = append(newPoly.Holes, &pb.CompressedPolygon{
					Points: CompressedPointsToPolylineBytes(hole.Points),
				})
			}
			reducedLocation.Data = append(reducedLocation.Data, newPoly)
		}
		output.Locations = append(output.Locations, reducedLocation)
	}
	return output
}

func Compress(input *pb.Locations, method pb.CompressMethod) (*pb.CompressedLocations, error) {
	switch method {
	case pb.CompressMethod_Polyline:
		return CompressWithPolyline(input), nil
	default:
		return nil, fmt.Errorf("pinpoint/reduce: unknown method %v", method)
	}
}

func DecompressWithPolyline(input *pb.CompressedLocations) *pb.Locations {
	output := &pb.Locations{}
	for _, location := range input.Locations {
		reducedLocation := &pb.Location{
			Name: location.Name,
		}
		for _, polygon := range location.Data {
			newPoly := &pb.Polygon{
				Points: DecompressedPolylineBytesToPoints(polygon.Points),
				Holes:  make([]*pb.Polygon, 0),
			}
			for _, hole := range polygon.Holes {
				newPoly.Holes = append(newPoly.Holes, &pb.Polygon{
					Points: DecompressedPolylineBytesToPoints(hole.Points),
				})
			}
			reducedLocation.Polygons = append(reducedLocation.Polygons, newPoly)
		}
		output.Locations = append(output.Locations, reducedLocation)
	}
	return output
}

func Decompress(input *pb.CompressedLocations) (*pb.Locations, error) {
	switch input.Method {
	case pb.CompressMethod_Polyline:
		return DecompressWithPolyline(input), nil
	default:
		return nil, fmt.Errorf("pinpoint/reduce: unknown method %v", input.Method)
	}
}
