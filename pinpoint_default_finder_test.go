package pinpoint_test

import (
	"fmt"
	"math/rand"
	"testing"

	gocitiesjson "github.com/deslittle/go-cities.json"
	pinpoint "github.com/deslittle/pinpoint"
)

var (
	defaultFinder *pinpoint.DefaultFinder
)

func init() {
	finder, err := pinpoint.NewDefaultFinder()
	if err != nil {
		panic(err)
	}
	defaultFinder = finder
}

func ExampleDefaultFinder_GetLocationName() {
	finder, err := pinpoint.NewDefaultFinder()
	if err != nil {
		panic(err)
	}
	fmt.Println(finder.GetLocationName(116.6386, 40.0786))
	// Output: Asia/Shanghai
}

func BenchmarkDefaultFinder_GetLocationName_Random_WorldCities(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		_ = defaultFinder.GetLocationName(p.Lng, p.Lat)
	}
}
