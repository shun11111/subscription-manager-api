package api

// このファイルは OpenAPI からコード生成するためのエントリーポイントです。
// oapi-codegen がインストールされている環境で、以下を実行すると
// internal/api/openapi.gen.go が自動生成されます。
//
//   go generate ./...
//
// 画面を追加したいときは、基本的に:
//  1. docs/openapi.yaml に path と schema を追加
//  2. go generate ./... で型とインターフェースを再生成
//  3. 必要なら service 層のビジネスロジックだけ実装
//
//go:generate oapi-codegen -generate types,server,echo-server -package api -o openapi.gen.go ../../docs/openapi.yaml


