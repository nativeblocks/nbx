package test

import (
	"github.com/nativeblocks/nbx"
	"github.com/nativeblocks/nbx/internal/compiler"
	"testing"
)

func TestToDsl(t *testing.T) {
	schema := `{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"required": ["name", "route", "type", "blocks"],
		"properties": {
			"name": {"type": "string"},
			"route": {"type": "string"},
			"blocks": {
				"type": "array",
				"minItems": 1,
				"items": {
					"type": "object",
					"required": ["keyType", "key"],
					"properties": {
						"keyType": {"type": "string"},
						"key": {"type": "string"},
						"slot": {"type": "string"}
					}
				}
			},
			"variables": {
				"type": "array",
				"items": {
					"type": "object",
					"required": ["key", "type", "value"],
					"properties": {
						"key": {"type": "string"},
						"type": {"type": "string"},
						"value": {"type": ["string", "number", "boolean"]}
					}
				}
			}
		}
	}`

	dsl := `frame(
    name = "welcome",
    route = "/welcome"
) {
    var visible: BOOLEAN = true
    var count: INT = 0
	var welcome: STRING = "Welcome"

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "nativeblocks/column", key = "mainColumn", visibility = visible, version = 1)
        .slot("content") {
            block(keyType = "nativeblocks/text", key = "welcome", visibility = visible, version = 1)
            .assignData(text = welcome)
        }
    }
}`

	frameDSL, err := nbx.Parse(dsl)
	if err != nil {
		t.Fatalf("Failed to parse DSL: %v", err)
	}

	frameJson, err := compiler.ToJson(frameDSL, schema, "")
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	if frameJson.Name != "welcome" {
		t.Errorf("Expected name 'welcome', got '%s'", frameJson.Name)
	}
	if frameJson.Route != "/welcome" {
		t.Errorf("Expected route '/welcome', got '%s'", frameJson.Route)
	}
	if len(frameJson.Variables) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(frameJson.Variables))
	}
	if len(frameJson.Blocks) < 3 {
		t.Errorf("Expected at least 3 blocks, got %d", len(frameJson.Blocks))
	}

	invalidDsl := `frame(
    name = "welcome",
    route = "/welcome"
) {
    block(keyType = "ROOT", key = "root")
}`

	invalidFrameDSL, err := nbx.Parse(invalidDsl)
	if err != nil {
		t.Fatalf("Failed to parse invalid DSL: %v", err)
	}

	_, err = compiler.ToJson(invalidFrameDSL, schema, "")
	if err != nil {
		t.Error("Expected no error but error", err)
	}

	customID := "custom-frame-id"
	frameJsonWithID, err := compiler.ToJson(frameDSL, schema, customID)
	if err != nil {
		t.Fatalf("Failed to convert to JSON with custom ID: %v", err)
	}
	if frameJsonWithID.Id != customID {
		t.Errorf("Expected frame ID '%s', got '%s'", customID, frameJsonWithID.Id)
	}
}
