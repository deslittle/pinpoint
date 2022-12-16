// CLI tool to combine multiple geojson files into one.
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FeatureCollection struct {
	Type     string                   `json:"type"`
	Features []map[string]interface{} `json:"features"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: combinegeojson <geojsonout file> <geojson file 1> <geojson file 2> ...")
		return
	}

	combinedCollection := &FeatureCollection{
		Type: "FeatureCollection",
	}

	for _, filePath := range os.Args[2:] {
		rawFile, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		// Parse the json into a BoundaryFile
		collection := &FeatureCollection{}
		err = json.Unmarshal(rawFile, collection)
		if err != nil {
			panic(err)
		}
		// Append the features to the combined collection
		combinedCollection.Features = append(combinedCollection.Features, collection.Features...)
	}

	// Convert back into json
	outFile, _ := json.Marshal(combinedCollection)
	outputPath := os.Args[1]
	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	_, _ = f.Write(outFile)
	fmt.Println(outputPath)

}
