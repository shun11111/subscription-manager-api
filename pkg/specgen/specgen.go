package specgen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// SpecDoc は仕様書の構造を表す
type SpecDoc struct {
	Title       string
	Description string
	Endpoints   []Endpoint
	Models      []Model
}

// Endpoint はAPIエンドポイントの定義
type Endpoint struct {
	Method      string   // GET, POST, PUT, DELETE
	Path        string   // /subscriptions
	Summary     string   // 説明
	Auth        bool     // 認証が必要か
	Request     string   // リクエスト型名（空文字列の場合はなし）
	Response    string   // レスポンス型名
	Params      []Param  // パスパラメータ
	ErrorCodes  []string // エラーコード（例: ["400", "401"]）
}

// Param はパスパラメータ
type Param struct {
	Name     string // id
	Type     string // uuid
	Location string // path
}

// Model はデータモデルの定義
type Model struct {
	Name   string
	Fields []Field
}

// Field はモデルのフィールド
type Field struct {
	Name     string
	Type     string
	Required bool
	Options  map[string]string // 追加のオプション（enum値など）
}

// ParseSpecFile はMarkdown仕様書をパースする
func ParseSpecFile(filePath string) (*SpecDoc, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	doc := &SpecDoc{
		Endpoints: []Endpoint{},
		Models:    []Model{},
	}

	// タイトルと概要を抽出
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			doc.Title = strings.TrimPrefix(line, "# ")
		}
		if strings.HasPrefix(line, "## 概要") && i+1 < len(lines) {
			doc.Description = strings.TrimSpace(lines[i+1])
		}
	}

	// エンドポイントを抽出
	doc.Endpoints = parseEndpoints(content)

	// データモデルを抽出
	doc.Models = parseModels(content)

	return doc, nil
}

// parseEndpoints はエンドポイント定義を抽出
func parseEndpoints(content string) []Endpoint {
	var endpoints []Endpoint

	// ### GET /path のパターンを探す
	re := regexp.MustCompile(`###\s+(GET|POST|PUT|DELETE|PATCH)\s+(/[^\s]+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		method := match[1]
		path := match[2]

		// エンドポイントの詳細を抽出
		endpoint := Endpoint{
			Method:     method,
			Path:       path,
			ErrorCodes: []string{},
		}

		// 説明、認証、リクエスト、レスポンスを抽出
		sectionStart := strings.Index(content, match[0])
		sectionEnd := strings.Index(content[sectionStart:], "\n##")
		if sectionEnd == -1 {
			sectionEnd = len(content) - sectionStart
		}
		section := content[sectionStart : sectionStart+sectionEnd]

		// 説明を抽出
		if descMatch := regexp.MustCompile(`\*\*説明\*\*:\s*(.+?)(?:\n|$)`).FindStringSubmatch(section); len(descMatch) > 1 {
			endpoint.Summary = strings.TrimSpace(descMatch[1])
		}

		// 認証を抽出
		if authMatch := regexp.MustCompile(`\*\*認証\*\*:\s*(必要|不要)`).FindStringSubmatch(section); len(authMatch) > 1 {
			endpoint.Auth = authMatch[1] == "必要"
		}

		// リクエストを抽出
		if reqMatch := regexp.MustCompile(`\*\*リクエスト\*\*:\s*([^\n]+)`).FindStringSubmatch(section); len(reqMatch) > 1 {
			req := strings.TrimSpace(reqMatch[1])
			if req != "なし" {
				endpoint.Request = req
			}
		}

		// レスポンスを抽出
		if resMatch := regexp.MustCompile(`\*\*レスポンス\*\*:\s*([^\n]+)`).FindStringSubmatch(section); len(resMatch) > 1 {
			endpoint.Response = strings.TrimSpace(resMatch[1])
		}

		// パラメータを抽出
		if paramSection := regexp.MustCompile(`\*\*パラメータ\*\*:\s*\n((?:\s+-[^\n]+\n?)+)`).FindStringSubmatch(section); len(paramSection) > 1 {
			paramLines := strings.Split(paramSection[1], "\n")
			for _, line := range paramLines {
				if strings.HasPrefix(line, "-") {
					// - `id`: uuid (path) の形式
					if paramMatch := regexp.MustCompile(`\`([^\`]+)\`:\s*(\w+)\s*\((\w+)\)`).FindStringSubmatch(line); len(paramMatch) > 3 {
						endpoint.Params = append(endpoint.Params, Param{
							Name:     paramMatch[1],
							Type:     paramMatch[2],
							Location: paramMatch[3],
						})
					}
				}
			}
		}

		// エラーコードを抽出
		if errMatch := regexp.MustCompile(`\*\*エラー\*\*:\s*([^\n]+)`).FindStringSubmatch(section); len(errMatch) > 1 {
			errStr := errMatch[1]
			// "400 (バリデーションエラー), 401 (認証エラー)" のような形式
			errCodes := regexp.MustCompile(`(\d{3})`).FindAllString(errStr, -1)
			endpoint.ErrorCodes = errCodes
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}

// parseModels はデータモデル定義を抽出
func parseModels(content string) []Model {
	var models []Model

	// ### ModelName のパターンを探す（データモデルセクション内）
	dataModelSection := regexp.MustCompile(`## データモデル\n(.*)`).FindStringSubmatch(content)
	if len(dataModelSection) < 2 {
		return models
	}

	modelSection := dataModelSection[1]

	// ### ModelName のパターン
	re := regexp.MustCompile(`###\s+(\w+)`)
	matches := re.FindAllStringSubmatch(modelSection, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		modelName := match[1]

		model := Model{Name: modelName, Fields: []Field{}}

		// モデルのセクションを抽出
		sectionStart := strings.Index(modelSection, match[0])
		sectionEnd := strings.Index(modelSection[sectionStart:], "\n###")
		if sectionEnd == -1 {
			sectionEnd = len(modelSection) - sectionStart
		}
		section := modelSection[sectionStart : sectionStart+sectionEnd]

		// フィールドを抽出
		fieldRe := regexp.MustCompile(`-\s*\`([^\`]+)\`:\s*(\w+)\s*(\(必須\)|\(オプション\))?`)
		fieldMatches := fieldRe.FindAllStringSubmatch(section, -1)

		for _, fieldMatch := range fieldMatches {
			if len(fieldMatch) < 3 {
				continue
			}
			field := Field{
				Name:     fieldMatch[1],
				Type:     fieldMatch[2],
				Required: fieldMatch[3] == "(必須)",
			}

			// enum値などを抽出（例: "monthly" | "yearly"）
			if enumMatch := regexp.MustCompile(`enum\s*\(([^)]+)\)`).FindStringSubmatch(section); len(enumMatch) > 1 {
				field.Options = map[string]string{
					"enum": enumMatch[1],
				}
			}

			model.Fields = append(model.Fields, field)
		}

		models = append(models, model)
	}

	return models
}

// GenerateOpenAPI は仕様書からOpenAPIを生成
func GenerateOpenAPI(specsDir string, outputPath string) error {
	// 全ての仕様書を読み込む
	var allEndpoints []Endpoint
	var allModels []Model

	err := filepath.WalkDir(specsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		if strings.HasSuffix(path, "template.md") || strings.HasSuffix(path, "README.md") {
			return nil
		}

		spec, err := ParseSpecFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		allEndpoints = append(allEndpoints, spec.Endpoints...)
		allModels = append(allModels, spec.Models...)

		return nil
	})

	if err != nil {
		return err
	}

	// OpenAPI形式に変換
	openAPI := buildOpenAPI(allEndpoints, allModels)

	// YAMLとして出力
	data, err := yaml.Marshal(openAPI)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}

// OpenAPISpec はOpenAPI仕様の構造
type OpenAPISpec struct {
	OpenAPI    string                 `yaml:"openapi"`
	Info       map[string]interface{} `yaml:"info"`
	Servers    []map[string]interface{} `yaml:"servers"`
	Paths      map[string]interface{} `yaml:"paths"`
	Components map[string]interface{} `yaml:"components"`
}

// buildOpenAPI はエンドポイントとモデルからOpenAPIを構築
func buildOpenAPI(endpoints []Endpoint, models []Model) *OpenAPISpec {
	// 簡易実装：既存のopenapi.yamlを読み込んで、pathsとcomponentsを更新する
	// 完全な実装は複雑なので、ここでは構造だけ定義
	return &OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: map[string]interface{}{
			"title":   "Subscription Manager API",
			"version": "1.0.0",
		},
		Servers: []map[string]interface{}{
			{
				"url":         "http://localhost:8080",
				"description": "ローカル開発サーバー",
			},
		},
		Paths:      map[string]interface{}{},
		Components: map[string]interface{}{},
	}
}

