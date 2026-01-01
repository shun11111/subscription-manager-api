.PHONY: gen-openapi gen-code help

# OpenAPI仕様書を生成（Markdown仕様書から）
gen-openapi:
	@echo "Generating OpenAPI from Markdown specs..."
	@go run cmd/specgen/main.go -specs docs/specs -output docs/openapi.yaml
	@echo "✅ OpenAPI generated: docs/openapi.yaml"

# Goコードを生成（OpenAPIから）
gen-code:
	@echo "Generating Go code from OpenAPI..."
	@export PATH=$$PATH:$$HOME/go/bin:/usr/local/go/bin && go generate ./...
	@echo "✅ Go code generated"

# 全てのコード生成を実行
gen-all: gen-openapi gen-code

help:
	@echo "Available targets:"
	@echo "  gen-openapi  - Generate OpenAPI YAML from Markdown specs"
	@echo "  gen-code     - Generate Go code from OpenAPI"
	@echo "  gen-all      - Run all code generation"

