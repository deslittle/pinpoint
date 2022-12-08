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
	input := &pb.Timezones{}
	if err := proto.Unmarshal(tzfrel.LiteData, input); err != nil {
		panic(err)
	}
	_finder, _ := pinpoint.NewFinderFromPB(input)
	finder = _finder
}

func initFull() {
	input := &pb.Timezones{}
	if err := proto.Unmarshal(tzfrel.FullData, input); err != nil {
		panic(err)
	}
	_finder, _ := pinpoint.NewFinderFromPB(input)
	fullFinder = _finder
}

func BenchmarkGetTimezoneName(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		_ = finder.GetTimezoneName(116.6386, 40.0786)
	}
}

func BenchmarkGetTimezoneNameAtEdge(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		_ = finder.GetTimezoneName(110.8571, 43.1483)
	}
}

func BenchmarkGetTimezoneName_Random_WorldCities(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		_ = finder.GetTimezoneName(p.Lng, p.Lat)
	}
}

func BenchmarkFullFinder_GetTimezoneName(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		_ = fullFinder.GetTimezoneName(116.6386, 40.0786)
	}
}

func BenchmarkFullFinder_GetTimezoneNameAtEdge(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		_ = fullFinder.GetTimezoneName(110.8571, 43.1483)
	}
}

func BenchmarkFullFinder_GetTimezoneName_Random_WorldCities(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		_ = fullFinder.GetTimezoneName(p.Lng, p.Lat)
	}
}

func ExampleFinder_GetTimezoneName() {
	input := &pb.Timezones{}

	// Lite data, about 16.7MB
	dataFile := tzfrel.LiteData

	// Full data, about 83.5MB
	// dataFile := tzfrel.FullData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)
	fmt.Println(finder.GetTimezoneName(116.6386, 40.0786))
	// Output: Asia/Shanghai
}

func ExampleFinder_GetTimezoneLoc() {
	input := &pb.Timezones{}

	// Lite data, about 16.7MB
	dataFile := tzfrel.LiteData

	// Full data, about 83.5MB
	// dataFile := tzfrel.FullData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)
	fmt.Println(finder.GetTimezoneLoc(116.6386, 40.0786))
	// Output: Asia/Shanghai <nil>
}

func ExampleFinder_GetTimezoneShapeByName() {
	input := &pb.Timezones{}

	// Lite data, about 16.7MB
	dataFile := tzfrel.LiteData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)
	pbtz, err := finder.GetTimezoneShapeByName("Asia/Shanghai")
	fmt.Printf("%v %v\n", pbtz.Name, err)
	// Output: Asia/Shanghai <nil>
}

func ExampleFinder_GetTimezoneShapeByShift() {
	input := &pb.Timezones{}

	// Lite data, about 16.7MB
	dataFile := tzfrel.LiteData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := pinpoint.NewFinderFromPB(input)
	pbtzs, _ := finder.GetTimezoneShapeByShift(28800)

	pbnames := make([]string, 0)
	for _, pbtz := range pbtzs {
		pbnames = append(pbnames, pbtz.Name)
	}
	sort.Strings(pbnames)

	fmt.Println(strings.Join(pbnames, ","))
	// Output: Asia/Brunei,Asia/Choibalsan,Asia/Hong_Kong,Asia/Irkutsk,Asia/Kuala_Lumpur,Asia/Kuching,Asia/Macau,Asia/Makassar,Asia/Manila,Asia/Shanghai,Asia/Singapore,Asia/Taipei,Asia/Ulaanbaatar,Australia/Perth,Etc/GMT-8
}