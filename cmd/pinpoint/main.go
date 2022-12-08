// tzf-cli tool for local query.
package main

import (
	"flag"
	"fmt"

	pinpoint "github.com/deslittle/pinpoint"
	"github.com/deslittle/pinpoint/pb"
	tzfrel "github.com/deslittle/tzf-rel"
	"google.golang.org/protobuf/proto"
)

var finder *pinpoint.Finder

func init() {
	input := &pb.CompressedLocations{}
	dataFile := tzfrel.LiteCompressData
	err := proto.Unmarshal(dataFile, input)
	if err != nil {
		panic(err)
	}
	finder, err = pinpoint.NewFinderFromCompressed(input)
	if err != nil {
		panic(err)
	}
}

func main() {
	var lng float64
	var lat float64
	flag.Float64Var(&lng, "lng", 116.3883, "longitude")
	flag.Float64Var(&lat, "lat", 39.9289, "lontitude")
	flag.Parse()

	fmt.Println(finder.GetTimezoneName(lng, lat))
}
