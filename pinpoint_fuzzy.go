package pinpoint

import (
	"github.com/deslittle/pinpoint/pb"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
)

// FuzzyFinder use a tile map to store location name. Data are made by
// [github.com/deslittle/pinpoint/cmd/preindexlocpb] which powerd by
// [github.com/deslittle/pinpoint/preindex.PreIndexLocations].
type FuzzyFinder struct {
	idxZoom int
	aggZoom int
	m       map[maptile.Tile][]string // locations may have common area
}

func NewFuzzyFinderFromPB(input *pb.PreindexLocations) (*FuzzyFinder, error) {
	f := &FuzzyFinder{
		m:       make(map[maptile.Tile][]string),
		idxZoom: int(input.IdxZoom),
		aggZoom: int(input.AggZoom),
	}
	for _, item := range input.Keys {
		tile := maptile.New(uint32(item.X), uint32(item.Y), maptile.Zoom(item.Z))
		if _, ok := f.m[tile]; !ok {
			f.m[tile] = make([]string, 0)
		}
		f.m[tile] = append(f.m[tile], item.Name)
	}
	return f, nil
}

func (f *FuzzyFinder) GetLocationName(lng float64, lat float64) string {
	names, err := f.GetLocationNames(lng, lat)
	if err != nil {
		return ""
	}
	return names[0]
}

func (f *FuzzyFinder) GetLocationNames(lng float64, lat float64) ([]string, error) {
	p := orb.Point{lng, lat}
	for z := f.aggZoom; z <= f.idxZoom; z++ {
		key := maptile.At(p, maptile.Zoom(z))
		v, ok := f.m[key]
		if ok {
			return v, nil
		}
	}
	return nil, ErrNoLocationFound
}
