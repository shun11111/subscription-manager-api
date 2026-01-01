package main

// generate.go はコード生成用のスクリプト
// 将来的に oapi-codegen を使用して OpenAPI からコードを生成する際に使用

import (
	"fmt"
)

func main() {
	// oapi-codegen を使用してコードを生成
	// 例: oapi-codegen -generate types,server -package api docs/openapi.yaml > internal/api/generated.go

	fmt.Println("Code generation script")
	fmt.Println("Usage:")
	fmt.Println("  go run generate.go")
	fmt.Println("")
	fmt.Println("To generate code from OpenAPI spec:")
	fmt.Println("  oapi-codegen -generate types,server -package api docs/openapi.yaml > internal/api/generated.go")

	// 実際の生成コマンドを実行する場合は以下をコメントアウト
	/*
	cmd := exec.Command("oapi-codegen",
		"-generate", "types,server",
		"-package", "api",
		"docs/openapi.yaml",
	)
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(string(output))
	*/
}

