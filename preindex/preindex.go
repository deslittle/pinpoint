// Package preindex
//
// # Background
//
// The Ray Casting algorithm's time complexity is $O(n^2)$ which is expensive for
// high throughput API that severing geo based data like weather forecasts.
// And most of these requests are came from big cities around the world.
//
// If we can reduce these location's query execution times, our API could got
// performance improvements.
//
// # How to
//
// Preindex's logic is very simple, generate map tiles around a multi polygon,
// and exclude 1/2 edge layer, then merge to upper tiles. Then dumps all the tiles's
// X/Y/Z and location to Protocol Buffer based data.
//
// A sample image of output tiles show on maps:
// https://user-images.githubusercontent.com/13536789/200174943-7d40661e-bda5-4b79-a867-ec637e245a49.png
package preindex

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/deslittle/pinpoint/convert"
	"github.com/deslittle/pinpoint/pb"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/maptile/tilecover"
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/lotsa"
	"golang.org/x/exp/maps"
)

// Drop most outside layer of tile, since tile may cover area not included in location.
func DropEdgeTiles(tiles []maptile.Tile) []maptile.Tile {
	ret := []maptile.Tile{}
	tilehash := map[maptile.Tile]bool{}

	// setup tilehash
	for _, tile := range tiles {
		tilehash[tile] = true
	}

	// filter all neighbor in tiles
	for _, tile := range tiles {
		neighbors := []maptile.Tile{
			maptile.New(tile.X-1, tile.Y-1, tile.Z),
			maptile.New(tile.X, tile.Y-1, tile.Z),
			maptile.New(tile.X+1, tile.Y-1, tile.Z),

			maptile.New(tile.X-1, tile.Y, tile.Z),
			// maptile.New(tile.X, tile.Y, tile.Z),
			maptile.New(tile.X+1, tile.Y, tile.Z),

			maptile.New(tile.X-1, tile.Y+1, tile.Z),
			maptile.New(tile.X, tile.Y+1, tile.Z),
			maptile.New(tile.X+1, tile.Y+1, tile.Z),
		}

		var allNeighorIn bool = func() bool {
			for _, neighborTile := range neighbors {
				if _, ok := tilehash[neighborTile]; !ok {
					return false
				}
			}
			return true
		}()
		if !allNeighorIn {
			continue
		}
		ret = append(ret, tile)
	}

	return ret
}

func EnsureInside(geopolys []*geometry.Poly, tiles []maptile.Tile) []maptile.Tile {
	insideLocTiles := []maptile.Tile{}
	for _, tile := range tiles {
		minLng := tile.Bound().Min.Lon()
		minLat := tile.Bound().Min.Lat()
		maxLng := tile.Bound().Max.Lon()
		maxLat := tile.Bound().Max.Lat()

		geometryPoints := []geometry.Point{
			{X: minLng, Y: minLat},
			{X: maxLng, Y: minLat},
			{X: maxLng, Y: maxLat},
			{X: minLng, Y: maxLat},
			{X: minLng, Y: minLat},
		}
		tilePoly := geometry.NewPoly(geometryPoints, nil, nil)

		for _, geopoly := range geopolys {
			if !geopoly.ContainsPoly(tilePoly) {
				continue
			}
			for _, point := range geometryPoints {
				if !geopoly.ContainsPoint(point) {
					continue
				}
			}
		}
		insideLocTiles = append(insideLocTiles, tile)
	}
	return insideLocTiles
}

// PreIndexLocation will gen tiles at idxZoom level and merge up to aggZoom.
//
// The `idxZoom` level tiles will be removed before final return.
func PreIndexLocation(input *pb.Location, idxZoom, aggZoom, maxZoomLevelToKeep maptile.Zoom, dropEdgeLayger int) ([]*pb.PreindexLocation, error) {
	// Generate all tiles event not included in location shape
	tiles := []maptile.Tile{}
	for _, poly := range input.Polygons {
		orbPoly := orb.Polygon{}

		ring := orb.Ring{}
		for _, point := range poly.Points {
			ring = append(ring, orb.Point{float64(point.Lng), float64(point.Lat)})
		}
		// bypass too little
		if len(ring) < 10 {
			continue
		}
		// add first point
		ring = append(ring, ring[0])
		orbPoly = append(orbPoly, ring)

		// add polygon holes
		for _, hole := range poly.Holes {
			holering := orb.Ring{}
			for _, point := range hole.Points {
				holering = append(holering, orb.Point{float64(point.Lng), float64(point.Lat)})
			}
			if len(holering) < 3 {
				continue
			}
			holering = append(holering, holering[0])
			orbPoly = append(orbPoly, holering)
		}

		// gen polygon tiles
		polytiles, err := tilecover.Geometry(orbPoly, idxZoom)
		if err != nil {
			panic(err)
		}
		tiles = append(tiles, maps.Keys(polytiles)...)
	}
	// unable to agg
	if len(tiles) < 9 {
		return nil, fmt.Errorf("too little")
	}

	// Iter all tile's polygon if inside original polygon
	geopolys := convert.FromLocationPBToGeometryPoly(input)
	insideLocTiles := EnsureInside(geopolys, tiles)

	// Drop edge tiles
	for i := 0; i < dropEdgeLayger; i++ {
		insideLocTiles = DropEdgeTiles(insideLocTiles)
	}

	// Gen tileset
	newtileset := maptile.Set{}
	for _, tile := range insideLocTiles {
		newtileset[tile] = true
	}

	// Merge all filterd tiles
	mergedtiles := maptile.Set{}
	for _, tile := range EnsureInside(geopolys, maps.Keys(tilecover.MergeUp(newtileset, aggZoom))) {
		mergedtiles[tile] = true
	}

	// // Dumps JSON for debug
	// b, _ := json.Marshal(mergedtiles.ToFeatureCollection())
	// _ = os.WriteFile("preindexex.geojson", b, 0644)

	// Dumps as pb
	ret := []*pb.PreindexLocation{}
	for _, v := range maps.Keys(mergedtiles) {
		if int(v.Z) > int(maxZoomLevelToKeep) {
			continue
		}
		ret = append(ret, &pb.PreindexLocation{
			Name: input.Name,
			X:    int32(v.X),
			Y:    int32(v.Y),
			Z:    int32(v.Z),
		})
	}
	return ret, nil
}

func PreIndexLocations(input *pb.Locations, idxZoom, aggZoom, maxZoomLevelToKeep maptile.Zoom, dropEdgeLayger int) *pb.PreindexLocations {
	ret := &pb.PreindexLocations{
		IdxZoom: int32(idxZoom),
		AggZoom: int32(aggZoom),
		Keys:    make([]*pb.PreindexLocation, 0),
	}

	lock := &sync.Mutex{}
	lotsa.Ops(len(input.Locations), runtime.NumCPU()*2, func(i, thread int) {
		tz := input.Locations[i]
		preindexes, err := PreIndexLocation(tz, idxZoom, aggZoom, maxZoomLevelToKeep, dropEdgeLayger)
		if err != nil {
			return
		}
		lock.Lock()
		ret.Keys = append(ret.Keys, preindexes...)
		lock.Unlock()
	})
	return ret
}

func PreIndexLocationsToGeoJSON(input *pb.PreindexLocations) []byte {
	tileset := maptile.Set{}
	for _, key := range input.Keys {
		tileset[maptile.New(uint32(key.X), uint32(key.Y), maptile.Zoom(key.Z))] = true
	}
	b, _ := json.Marshal(tileset.ToFeatureCollection())
	return b
}
