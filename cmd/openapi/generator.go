package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"

	"github.com/b1tray3r/go-openapi3/internal/openapi"
)

func main() {
	var output string

	flag.StringVar(&output, "path", "", "Path to use for generating OpenAPI 3 files")
	flag.Parse()

	if output == "" {
		log.Fatalln("path is required")
	}

	openapi3 := openapi.NewDefinition()

	// openapi3.json
	data, err := json.Marshal(&openapi3)
	if err != nil {
		log.Fatalf("Couldn't marshal json: %s", err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		if err := os.MkdirAll(output, 0750); err != nil {
			log.Fatalf("Couldn't create directory: %s", err)
		}
	}

	target := path.Join(output, "openapi3.json")
	if _, err := os.Stat(output); os.IsNotExist(err) {
		file, err := os.Create(target)
		if err != nil {
			log.Fatalf("Couldn't create file: %s", err)
		}
		defer file.Close()
	}

	if err := os.WriteFile(target, data, 0640); err != nil {
		log.Fatalf("Couldn't write json: %s", err)
	}
}
