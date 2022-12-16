package pinpoint_test

import (
	"fmt"
	"testing"

	pinpoint "github.com/deslittle/pinpoint"
	usstates "github.com/deslittle/pinpoint-us-states"
	"github.com/deslittle/pinpoint/pb"
	"google.golang.org/protobuf/proto"
)

var (
	fuzzyFinder *pinpoint.FuzzyFinder
)

func init() {
	input := &pb.PreindexLocations{}
	if err := proto.Unmarshal(usstates.PreindexData, input); err != nil {
		panic(err)
	}
	_fuzzyFinder, err := pinpoint.NewFuzzyFinderFromPB(input)
	if err != nil {
		panic(err)
	}
	fuzzyFinder = _fuzzyFinder
}

// func TestFuzzySupports(t *testing.T) {
// 	failCount := 0
// 	for _, city := range gocitiesjson.Cities {
// 		name := fuzzyFinder.GetLocationName(city.Lng, city.Lat)
// 		if name == "" {
// 			failCount += 1
// 		}
// 	}
// 	// more than 10%
// 	if failCount/len(gocitiesjson.Cities)*100 > 10 {
// 		t.Errorf("has too many covered cities %v", failCount)
// 	}
// }

func ExampleFuzzyFinder_GetLocationName() {
	input := &pb.PreindexLocations{}
	if err := proto.Unmarshal(usstates.PreindexData, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFuzzyFinderFromPB(input)

	loc := finder.GetLocationName(-74.666645, 40.736032)
	fmt.Println(loc)
	// Output: 34

}

func BenchmarkFuzzyFinder_GetLocationName_Random_WorldCities(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		//p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		p := struct{ Lat, Lng float64 }{Lat: 40.0786, Lng: 116.6386}
		_ = fuzzyFinder.GetLocationName(p.Lng, p.Lat)
	}
}

func FuzzFuzzyFinder_GetLocationName(f *testing.F) {
	f.Add(116.3883, 39.9289)
	f.Fuzz(func(t *testing.T, a float64, b float64) {
		ret, err := fuzzyFinder.GetLocationNames(a, b)
		if err == nil && len(ret) == 0 {
			t.Errorf("bad return %v, %v", ret, err)
		}
	})
}
