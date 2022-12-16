package pinpoint_test

import (
	"fmt"
	"testing"

	pinpoint "github.com/deslittle/pinpoint"
)

var (
	CombinedFinder *pinpoint.ExampleCombinedFinder
)

func init() {
	finder, err := pinpoint.NewExampleCombinedFinder()
	if err != nil {
		panic(err)
	}
	CombinedFinder = finder
}

func ExampleCombinedFinder_GetLocationName() {
	finder, err := pinpoint.NewExampleCombinedFinder()
	if err != nil {
		panic(err)
	}
	fmt.Println(finder.GetLocationName(-74.03440821618342, 40.71579135708155))
	// Output: 34
}

func BenchmarkExampleCombinedFinder_GetLocationName_Random_WorldCities(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		//p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		p := struct{ Lat, Lng float64 }{Lat: 40.0786, Lng: 116.6386}
		_ = CombinedFinder.GetLocationName(p.Lng, p.Lat)
	}
}
