// CLI tool to combine multiple geojson files into one.
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []Features `json:"features"`
}

type Features struct {
	Type       string                 `json:"type"`
	Properties Properties             `json:"properties"`
	Geometry   map[string]interface{} `json:"geometry"`
}

type Properties struct {
	Name string `json:"name"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: combinegeojson <geojsonin file> <geojsonout file> <whitelist string 1> <whitelist string 1> ...")
		return
	}

	// Create a map of the whitelist strings
	whitelist := make(map[string]struct{})
	for _, allowedName := range os.Args[3:] {
		whitelist[allowedName] = struct{}{}
	}

	// Parse the json into a FeatureCollection
	rawFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	collection := &FeatureCollection{}
	err = json.Unmarshal(rawFile, collection)
	if err != nil {
		panic(err)
	}

	filteredCollection := &FeatureCollection{
		Type: "FeatureCollection",
	}
	// Loop through the features and only add the ones that are in the whitelist
	for _, feature := range collection.Features {
		if _, ok := whitelist[feature.Properties.Name]; ok {
			filteredCollection.Features = append(filteredCollection.Features, feature)
		}
	}

	// Convert back into json
	outFile, _ := json.Marshal(filteredCollection)
	outputPath := os.Args[2]
	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	_, _ = f.Write(outFile)
	fmt.Println(outputPath)

}
