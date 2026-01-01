package main

import (
	"flag"
	"log"
	"path/filepath"

	"subscription-manager-api/pkg/specgen"
)

func main() {
	specsDir := flag.String("specs", "docs/specs", "仕様書ディレクトリ")
	output := flag.String("output", "docs/openapi.yaml", "出力ファイル")
	flag.Parse()

	absSpecsDir, err := filepath.Abs(*specsDir)
	if err != nil {
		log.Fatalf("Failed to resolve specs directory: %v", err)
	}

	absOutput, err := filepath.Abs(*output)
	if err != nil {
		log.Fatalf("Failed to resolve output path: %v", err)
	}

	log.Printf("Generating OpenAPI from specs in %s...", absSpecsDir)
	if err := specgen.GenerateOpenAPI(absSpecsDir, absOutput); err != nil {
		log.Fatalf("Failed to generate OpenAPI: %v", err)
	}

	log.Printf("OpenAPI generated: %s", absOutput)
}

