package pinpoint_test

import (
	"fmt"
	"testing"

	pinpoint "github.com/deslittle/pinpoint"
	usstates "github.com/deslittle/pinpoint-us-states"
	"github.com/deslittle/pinpoint/pb"
	"github.com/loov/hrtime/hrtesting"
	"google.golang.org/protobuf/proto"
)

var (
	finder     *pinpoint.Finder
	fullFinder *pinpoint.Finder
)

func init() {
	initLite()
	initFull()
}

func initLite() {
	input := &pb.Locations{}
	if err := proto.Unmarshal(usstates.LiteData, input); err != nil {
		panic(err)
	}
	_finder, _ := pinpoint.NewFinderFromPB(input)
	finder = _finder
}

func initFull() {
	input := &pb.Locations{}
	if err := proto.Unmarshal(usstates.FullData, input); err != nil {
		panic(err)
	}
	_finder, _ := pinpoint.NewFinderFromPB(input)
	fullFinder = _finder
}

func BenchmarkGetLocationName(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()
	for bench.Next() {
		_ = finder.GetLocationName(116.6386, 40.0786)
	}
}

func BenchmarkGetLocationNameAtEdge(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()
	for bench.Next() {
		_ = finder.GetLocationName(110.8571, 43.1483)
	}
}

func BenchmarkGetLocationName_Random_WorldCities(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()
	for bench.Next() {
		//p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		p := struct{ Lat, Lng float64 }{Lat: 40.0786, Lng: 116.6386}
		_ = finder.GetLocationName(p.Lng, p.Lat)
	}
}

func BenchmarkFullFinder_GetLocationName(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()
	for bench.Next() {
		_ = fullFinder.GetLocationName(116.6386, 40.0786)
	}
}

func BenchmarkFullFinder_GetLocationNameAtEdge(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()
	for bench.Next() {
		_ = fullFinder.GetLocationName(110.8571, 43.1483)
	}
}

func BenchmarkFullFinder_GetLocationName_Random_WorldCities(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()
	for bench.Next() {
		//p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		p := struct{ Lat, Lng float64 }{Lat: 40.0786, Lng: 116.6386}
		_ = fullFinder.GetLocationName(p.Lng, p.Lat)
	}
}

func ExampleFinder_GetLocationName() {
	input := &pb.Locations{}

	// Lite data, about 16.7MB
	//dataFile := usstates.LiteData

	// Full data, about 83.5MB
	dataFile := usstates.FullData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)

	// Need to use a point not so close to New York border for Fuzzy Search
	loc, err := finder.GetLocation(-74.03440821618342, 40.71579135708155)
	fmt.Printf("%v %v\n", loc.GetName(), err)
	// Output: 34 <nil>
}

func ExampleFinder_GetLocationShapeByName() {
	input := &pb.Locations{}

	// Lite data, about 16.7MB
	dataFile := usstates.FullData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)
	pbloc, err := finder.GetLocationShapeByName("34")
	fmt.Printf("%v %v\n", pbloc.GetName(), err)
	// Output: 34 <nil>
}
