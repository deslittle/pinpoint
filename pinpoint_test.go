package pinpoint_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"sort"

	gocitiesjson "github.com/deslittle/go-cities.json"
	pinpoint "github.com/deslittle/pinpoint"
	"github.com/deslittle/pinpoint/pb"
	tzfrel "github.com/deslittle/tzf-rel"
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
	if err := proto.Unmarshal(tzfrel.LiteData, input); err != nil {
		panic(err)
	}
	_finder, _ := pinpoint.NewFinderFromPB(input)
	finder = _finder
}

func initFull() {
	input := &pb.Locations{}
	if err := proto.Unmarshal(tzfrel.FullData, input); err != nil {
		panic(err)
	}
	_finder, _ := pinpoint.NewFinderFromPB(input)
	fullFinder = _finder
}

func BenchmarkGetLocationName(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		_ = finder.GetLocationName(116.6386, 40.0786)
	}
}

func BenchmarkGetLocationNameAtEdge(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		_ = finder.GetLocationName(110.8571, 43.1483)
	}
}

func BenchmarkGetLocationName_Random_WorldCities(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		_ = finder.GetLocationName(p.Lng, p.Lat)
	}
}

func BenchmarkFullFinder_GetLocationName(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		_ = fullFinder.GetLocationName(116.6386, 40.0786)
	}
}

func BenchmarkFullFinder_GetLocationNameAtEdge(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		_ = fullFinder.GetLocationName(110.8571, 43.1483)
	}
}

func BenchmarkFullFinder_GetLocationName_Random_WorldCities(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		_ = fullFinder.GetLocationName(p.Lng, p.Lat)
	}
}

func ExampleFinder_GetLocationName() {
	input := &pb.Locations{}

	// Lite data, about 16.7MB
	dataFile := tzfrel.LiteData

	// Full data, about 83.5MB
	// dataFile := tzfrel.FullData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)
	fmt.Println(finder.GetLocationName(116.6386, 40.0786))
	// Output: Asia/Shanghai
}

func ExampleFinder_GetLocationLoc() {
	input := &pb.Locations{}

	// Lite data, about 16.7MB
	dataFile := tzfrel.LiteData

	// Full data, about 83.5MB
	// dataFile := tzfrel.FullData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)
	fmt.Println(finder.GetLocationLoc(116.6386, 40.0786))
	// Output: Asia/Shanghai <nil>
}

func ExampleFinder_GetLocationShapeByName() {
	input := &pb.Locations{}

	// Lite data, about 16.7MB
	dataFile := tzfrel.LiteData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)
	pbtz, err := finder.GetLocationShapeByName("Asia/Shanghai")
	fmt.Printf("%v %v\n", pbtz.Name, err)
	// Output: Asia/Shanghai <nil>
}

func ExampleFinder_GetLocationShapeByShift() {
	input := &pb.Locations{}

	// Lite data, about 16.7MB
	dataFile := tzfrel.LiteData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)
	pbtzs, _ := finder.GetLocationShapeByShift(28800)

	pbnames := make([]string, 0)
	for _, pbtz := range pbtzs {
		pbnames = append(pbnames, pbtz.Name)
	}
	sort.Strings(pbnames)

	fmt.Println(strings.Join(pbnames, ","))
	// Output: Asia/Brunei,Asia/Choibalsan,Asia/Hong_Kong,Asia/Irkutsk,Asia/Kuala_Lumpur,Asia/Kuching,Asia/Macau,Asia/Makassar,Asia/Manila,Asia/Shanghai,Asia/Singapore,Asia/Taipei,Asia/Ulaanbaatar,Australia/Perth,Etc/GMT-8
}
