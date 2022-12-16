package pinpoint

import (
	"runtime"

	usstates "github.com/deslittle/pinpoint-us-states"
	"github.com/deslittle/pinpoint/pb"
	"google.golang.org/protobuf/proto"
)

// combinedFinder is an example implimentation using the `pinpoint-us-states`
// repo which combines both [FuzzyFinder] and [Finder].
//
// It's designed for performance first and allow some not so correct return at some area.
type ExampleCombinedFinder struct {
	fuzzyFinder *FuzzyFinder
	finder      *Finder
}

func NewExampleCombinedFinder() (*ExampleCombinedFinder, error) {
	fuzzyFinder, err := func() (*FuzzyFinder, error) {
		input := &pb.PreindexLocations{}
		if err := proto.Unmarshal(usstates.PreindexData, input); err != nil {
			panic(err)
		}
		return NewFuzzyFinderFromPB(input)
	}()
	if err != nil {
		return nil, err
	}

	finder, err := func() (*Finder, error) {
		input := &pb.CompressedLocations{}
		if err := proto.Unmarshal(usstates.LiteCompressData, input); err != nil {
			panic(err)
		}
		return NewFinderFromCompressed(input, SetDropPBLoc)
	}()
	if err != nil {
		return nil, err
	}

	f := &ExampleCombinedFinder{}
	f.fuzzyFinder = fuzzyFinder
	f.finder = finder

	// Force free mem by probuf, about 80MB
	runtime.GC()

	return f, nil
}

func (f *ExampleCombinedFinder) GetLocationName(lng float64, lat float64) string {
	fuzzyRes := f.fuzzyFinder.GetLocationName(lng, lat)
	if fuzzyRes != "" {
		return fuzzyRes
	}
	name := f.finder.GetLocationName(lng, lat)
	if name != "" {
		return name
	}
	for _, dx := range []float64{-0.02, 0, 0.02} {
		for _, dy := range []float64{-0.02, 0, 0.02} {
			dlng := dx + lng
			dlat := dy + lat
			fuzzyRes := f.fuzzyFinder.GetLocationName(dlng, dlat)
			if fuzzyRes != "" {
				return fuzzyRes
			}
			name := f.finder.GetLocationName(dlng, dlat)
			if name != "" {
				return name
			}
		}
	}
	return ""
}

func (f *ExampleCombinedFinder) GetLocationNames(lng float64, lat float64) ([]string, error) {
	fuzzyRes, err := f.fuzzyFinder.GetLocationNames(lng, lat)
	if err == nil {
		return fuzzyRes, nil
	}
	return f.finder.GetLocationNames(lng, lat)
}

func (f *ExampleCombinedFinder) LocationNames() []string {
	return f.finder.LocationNames()
}
