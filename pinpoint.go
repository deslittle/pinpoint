// package pinpoint is a package convert (lng,lat) to location.

package pinpoint

import (
	"errors"
	"fmt"
	"sort"

	"github.com/deslittle/pinpoint/convert"
	"github.com/deslittle/pinpoint/pb"
	"github.com/deslittle/pinpoint/reduce"
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/rtree"
)

var ErrNoLocationFound = errors.New("pinpoint: no location found")

type Option struct {
	DropPBLoc bool
}

type OptionFunc = func(opt *Option)

// SetDropPBLoc will make Finder not save [github.com/deslittle/pinpoint/pb.Location] in memory
func SetDropPBLoc(opt *Option) {
	opt.DropPBLoc = true
}

type locitem struct {
	pbloc *pb.Location
	name  string
	polys []*geometry.Poly
	min   [2]float64
	max   [2]float64
}

func newNotFoundErr(lng float64, lat float64) error {
	return fmt.Errorf("pinpoint: not found for %v,%v", lng, lat)
}

func (i *locitem) ContainsPoint(p geometry.Point) bool {
	for _, poly := range i.polys {
		if poly.ContainsPoint(p) {
			return true
		}
	}
	return false
}

func (i *locitem) GetMinMax() ([2]float64, [2]float64) {
	retmin := [2]float64{
		i.polys[0].Rect().Min.X,
		i.polys[0].Rect().Min.Y,
	}
	retmax := [2]float64{
		i.polys[0].Rect().Max.X,
		i.polys[0].Rect().Max.Y,
	}

	for _, poly := range i.polys {
		minx := poly.Rect().Min.X
		miny := poly.Rect().Min.Y
		if minx < retmin[0] {
			retmin[0] = minx
		}
		if miny < retmin[1] {
			retmin[1] = miny
		}

		maxx := poly.Rect().Max.X
		maxy := poly.Rect().Max.Y
		if maxx > retmax[0] {
			retmax[0] = maxx

		}
		if maxy > retmax[1] {
			retmax[1] = maxy
		}
	}
	return retmin, retmax
}

// Finder is based on point-in-polygon search algo.
//
// Memeory will use about 100MB if lite data and 1G if full data.
// Performance is very stable and very accuate.
type Finder struct {
	items   []*locitem
	names   []string
	reduced bool
	tr      *rtree.RTreeG[*locitem]
	opt     *Option
}

func NewFinderFromRawJSON(input *convert.BoundaryFile, opts ...OptionFunc) (*Finder, error) {
	locations, err := convert.Do(input)
	if err != nil {
		return nil, err
	}
	return NewFinderFromPB(locations, opts...)
}

func NewFinderFromPB(input *pb.Locations, opts ...OptionFunc) (*Finder, error) {

	items := make([]*locitem, 0)
	names := make([]string, 0)

	opt := &Option{}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	tr := &rtree.RTreeG[*locitem]{}
	for _, location := range input.Locations {
		names = append(names, location.Name)

		newItem := &locitem{
			name: location.Name,
		}
		if !opt.DropPBLoc {
			newItem.pbloc = location
		}
		for _, polygon := range location.Polygons {

			newPoints := make([]geometry.Point, 0)
			for _, point := range polygon.Points {
				newPoints = append(newPoints, geometry.Point{
					X: float64(point.Lng),
					Y: float64(point.Lat),
				})
			}

			holes := [][]geometry.Point{}
			for _, holePoly := range polygon.Holes {
				newHolePoints := make([]geometry.Point, 0)
				for _, point := range holePoly.Points {
					newHolePoints = append(newHolePoints, geometry.Point{
						X: float64(point.Lng),
						Y: float64(point.Lat),
					})
				}
				holes = append(holes, newHolePoints)
			}

			newPoly := geometry.NewPoly(newPoints, holes, nil)
			newItem.polys = append(newItem.polys, newPoly)
		}
		minp, maxp := newItem.GetMinMax()

		newItem.min = minp
		newItem.max = maxp

		items = append(items, newItem)
		tr.Insert(minp, maxp, newItem)
	}
	finder := &Finder{}
	finder.items = items
	finder.names = names
	finder.reduced = input.Reduced
	finder.tr = tr
	finder.opt = opt
	return finder, nil
}

func NewFinderFromCompressed(input *pb.CompressedLocations, opts ...OptionFunc) (*Finder, error) {
	locs, err := reduce.Decompress(input)
	if err != nil {
		return nil, err
	}
	return NewFinderFromPB(locs, opts...)
}

func getRTreeRangeShifed(lng float64, lat float64) float64 {
	if 73 < lng && lng < 140 && 8 < lat && lat < 54 {
		return 70.0
	}
	return 30.0
}

func (f *Finder) getItemInRanges(lng float64, lat float64) []*locitem {
	candicates := []*locitem{}

	// TODO(tzf): fix this range
	shifted := getRTreeRangeShifed(lng, lat)
	f.tr.Search([2]float64{lng - shifted, lat - shifted}, [2]float64{lng + shifted, lat + shifted}, func(min, max [2]float64, data *locitem) bool {
		candicates = append(candicates, data)
		return true
	})
	if len(candicates) == 0 {
		candicates = f.items
	}

	return candicates
}

func (f *Finder) getItem(lng float64, lat float64) ([]*locitem, error) {
	p := geometry.Point{
		X: float64(lng),
		Y: float64(lat),
	}
	ret := []*locitem{}
	candicates := f.getItemInRanges(lng, lat)
	if len(candicates) == 0 {
		return nil, ErrNoLocationFound
	}
	for _, item := range candicates {
		if item.ContainsPoint(p) {
			ret = append(ret, item)
		}
	}
	if len(ret) == 0 {
		return nil, newNotFoundErr(lng, lat)
	}
	return ret, nil
}

// GetLocationName will use alphabet order and return first matched result.
func (f *Finder) GetLocationName(lng float64, lat float64) string {
	p := geometry.Point{
		X: float64(lng),
		Y: float64(lat),
	}
	for _, item := range f.items {
		if item.ContainsPoint(p) {
			return item.name
		}
	}
	return ""
}

func (f *Finder) GetLocationNames(lng float64, lat float64) ([]string, error) {
	item, err := f.getItem(lng, lat)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for i := 0; i < len(item); i++ {
		ret = append(ret, item[i].name)
	}
	sort.Strings(ret)
	return ret, nil
}

func (f *Finder) GetLocation(lng float64, lat float64) (*pb.Location, error) {
	if f.opt.DropPBLoc {
		return nil, errors.New("pinpoint: not suppor when reduce mem")
	}
	item, err := f.getItem(lng, lat)
	if err != nil {
		return nil, err
	}
	return item[0].pbloc, nil
}

func (f *Finder) GetLocationShapeByName(name string) (*pb.Location, error) {
	for _, item := range f.items {
		if item.name == name {
			return item.pbloc, nil
		}
	}
	return nil, fmt.Errorf("location=%v not found", name)
}

func (f *Finder) LocationNames() []string {
	return f.names
}
