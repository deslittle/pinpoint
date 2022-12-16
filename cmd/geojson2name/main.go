// CLI tool to convert GeoJSON based Location boundary to pinpoints's Probuf format.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/deslittle/pinpoint/convert"
)

type BoundaryProperties struct {
	Type     string     `json:"type"` // FeatureCollection
	Features []Features `json:"features"`
}

type Properties map[string]interface{}

type Features struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Geometry struct {
	Coordinates interface{} `json:"coordinates"`
	Type        string      `json:"type"`
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: geojson2name <geojson file> <propertyKey string>")
		return
	}

	jsonFilePath := os.Args[1]
	featureId := os.Args[2]

	rawFile, err := os.ReadFile(jsonFilePath)
	if err != nil {
		panic(err)
	}

	// Parse the json into a BoundaryProperties
	var boundaryProperties BoundaryProperties
	err = json.Unmarshal(rawFile, &boundaryProperties)
	if err != nil {
		panic(err)
	}

	// Parse the json into a BoundaryFile
	outFileJson := &convert.BoundaryFile{}
	if err := json.Unmarshal(rawFile, outFileJson); err != nil {
		panic(err)
	}

	for i, feature := range boundaryProperties.Features {
		if feature.Properties[featureId] != nil {

			outFileJson.Features[i].Properties.Name = feature.Properties[featureId].(string)
		}
	}

	// Convert back into json
	outFile, _ := json.Marshal(outFileJson)
	outputPath := strings.Replace(jsonFilePath, ".json", "-named.json", 1)

	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	_, _ = f.Write(outFile)
	fmt.Println(outputPath)
}
