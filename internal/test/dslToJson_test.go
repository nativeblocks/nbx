package test

import (
	"github.com/nativeblocks/nbx"
	"github.com/nativeblocks/nbx/internal/compiler"
	"testing"
)

func TestToDsl(t *testing.T) {
	schema := `
{
  "$schema" : "http://json-schema.org/draft-07/schema#",
  "schema-version" : "projectId_0197d45c-90cf-7c11-a55e-edd2cef8db42",
  "type" : "object",
  "required" : [ "name", "route", "type", "variables", "blocks" ],
  "properties" : {
    "blocks" : {
      "items" : {
        "$ref" : "#/definitions/block"
      },
      "maxItems" : 1,
      "type" : "array"
    },
    "name" : {
      "type" : "string"
    },
    "route" : {
      "type" : "string"
    },
    "type" : {
      "enum" : [ "FRAME", "BOTTOM_SHEET", "DIALOG" ],
      "type" : "string"
    },
    "variables" : {
      "items" : {
        "properties" : {
          "key" : {
            "type" : "string"
          },
          "type" : {
            "enum" : [ "STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN" ],
            "type" : "string"
          },
          "value" : {
            "type" : "string"
          }
        },
        "required" : [ "key", "value", "type" ],
        "type" : "object"
      },
      "type" : "array",
      "uniqueItems" : true
    }
  },
  "definitions" : {
    "block" : {
      "properties" : {
        "actions" : {
          "items" : {
            "properties" : {
              "event" : {
                "enum" : [ "onAppear", "onDisappear", "onClick" ],
                "type" : "string"
              },
              "triggers" : {
                "items" : {
                  "$ref" : "#/definitions/trigger"
                },
                "type" : "array"
              }
            },
            "required" : [ "event", "triggers" ],
            "type" : "object"
          },
          "type" : "array"
        },
        "blocks" : {
          "items" : {
            "$ref" : "#/definitions/block"
          },
          "type" : "array"
        },
        "data" : {
          "items" : {
            "properties" : {
              "key" : {
                "enum" : [ ],
                "type" : "string"
              },
              "type" : {
                "enum" : [ "STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN" ],
                "type" : "string"
              },
              "value" : {
                "type" : "string"
              }
            },
            "required" : [ "key", "value", "type" ],
            "type" : "object"
          },
          "type" : "array"
        },
        "integrationVersion" : {
          "type" : "integer"
        },
        "key" : {
          "type" : "string"
        },
        "keyType" : {
          "enum" : [ "ROOT", "nativeblocks/column", "nativeblocks/row", "nativeblocks/text", "nativeblocks/button", "nativeblocks/image" ],
          "type" : "string"
        },
        "properties" : {
          "items" : {
            "properties" : {
              "key" : {
                "type" : "string"
              },
              "type" : {
                "enum" : [ "STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN" ],
                "type" : "string"
              },
              "valueDesktop" : {
                "type" : "string"
              },
              "valueMobile" : {
                "type" : "string"
              },
              "valueTablet" : {
                "type" : "string"
              }
            },
            "required" : [ "key", "valueMobile", "valueTablet", "valueDesktop", "type" ],
            "type" : "object"
          },
          "type" : "array"
        },
        "slot" : {
          "type" : "string"
        },
        "slots" : {
          "items" : {
            "properties" : {
              "slot" : {
                "enum" : [ "content" ],
                "type" : "string"
              }
            },
            "type" : "object"
          },
          "type" : "array"
        },
        "visibilityKey" : {
          "type" : "string"
        }
      },
      "required" : [ "keyType", "key", "visibilityKey", "slot", "slots", "integrationVersion", "data", "properties", "actions", "blocks" ],
      "type" : "object"
    },
    "trigger" : {
      "properties" : {
        "data" : {
          "items" : {
            "properties" : {
              "key" : {
                "enum" : [ ],
                "type" : "string"
              },
              "type" : {
                "enum" : [ "STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN" ],
                "type" : "string"
              },
              "value" : {
                "type" : "string"
              }
            },
            "required" : [ "key", "value", "type" ],
            "type" : "object"
          },
          "type" : "array"
        },
        "integrationVersion" : {
          "type" : "integer"
        },
        "keyType" : {
          "enum" : [ "nativeblocks/change_variable" ],
          "type" : "string"
        },
        "name" : {
          "type" : "string"
        },
        "properties" : {
          "items" : {
            "properties" : {
              "key" : {
                "enum" : [ ],
                "type" : "string"
              },
              "type" : {
                "enum" : [ "STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN" ],
                "type" : "string"
              },
              "value" : {
                "type" : "string"
              }
            },
            "required" : [ "key", "value", "type" ],
            "type" : "object"
          },
          "type" : "array"
        },
        "then" : {
          "enum" : [ "NEXT", "END", "SUCCESS", "FAILURE" ],
          "type" : "string"
        },
        "triggers" : {
          "items" : {
            "$ref" : "#/definitions/trigger"
          },
          "type" : "array"
        }
      },
      "required" : [ "keyType", "then", "name", "integrationVersion", "properties", "data", "triggers" ],
      "type" : "object"
    }
  }
}
`

	dsl := `
frame(
    name = "welcome",
    route = "/welcome"
) {
    var visible: BOOLEAN = true
    var enable: BOOLEAN = true
    var count: INT = 0
    var increaseButton: STRING = "+"
    var decreaseButton: STRING = "-"
    var welcome: STRING = "Welcome to Nativeblocks"
    var logo: STRING = "https://nativeblocks.io/nativeblocks_logo.png"

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "nativeblocks/column", key = "mainColumn", visibility = visible, version = 1)
        .prop(horizontalAlignment = (valueMobile = "centerHorizontally", valueTablet = "centerHorizontally", valueDesktop = "centerHorizontally"))
        .prop(width = (valueMobile = "match", valueTablet = "match", valueDesktop = "match"))
        .prop(height = (valueMobile = "match", valueTablet = "match", valueDesktop = "match"))
        .slot("content") {
            block(keyType = "nativeblocks/column", key = "nativeblocksColumn", VISIBILITY = visible, version = 1)
            .prop(horizontalAlignment = (valueMobile = "centerHorizontally", valueTablet = "centerHorizontally", valueDesktop = "centerHorizontally"))
            .prop(paddingTop = (valueMobile = "64", valueTablet = "64", valueDesktop = "64"))
            .prop(weight = (valueMobile = "0.4f", valueTablet = "0.4f", valueDesktop = "0.4f"))
            .prop(verticalArrangement = (valueMobile = "spaceAround", valueTablet = "spaceAround", valueDesktop = "spaceAround"))
            .slot("content") {
                block(keyType = "nativeblocks/image", key = "logo", visibility = visible, version = 1)
                .prop(scaleType = (valueMobile = "inside", valueTablet = "inside", valueDesktop = "inside"))
                .prop(width = (valueMobile = "128", valueTablet = "128", valueDesktop = "128"))
                .prop(height = (valueMobile = "128", valueTablet = "128", valueDesktop = "128"))
                .data(imageUrl = logo)

                block(keyType = "nativeblocks/text", key = "welcome", visibility = visible, version = 1)
                .prop(fontSize = (valueMobile = "24", valueTablet = "24", valueDesktop = "24"))
                .prop(textAlign = (valueMobile = "center", valueTablet = "center", valueDesktop = "center"))
                .prop(width = (valueMobile = "wrap", valueTablet = "wrap", valueDesktop = "wrap"))
                .data(text = welcome)
            }
            block(keyType = "nativeblocks/row", key = "buttonsRow", visibility = visible, version = 1)
            .prop(horizontalArrangement = (valueMobile = "spaceAround", valueTablet = "spaceAround", valueDesktop = "spaceAround"))
            .prop(verticalAlignment = (valueMobile = "centerVertically", valueTablet = "centerVertically", valueDesktop = "centerVertically"))
            .prop(paddingTop = (valueMobile = "12", valueTablet = "12", valueDesktop = "12"))
            .prop(weight = (valueMobile = "0.6f", valueTablet = "0.6f", valueDesktop = "0.6f"))
            .slot("content") {
                block(keyType = "nativeblocks/button", key = "decreaseButton", visibility = visible, version = 1)
                .prop(backgroundColor = (valueMobile = "#2563EB", valueTablet = "#2563EB", valueDesktop = "#2563EB"))
                .prop(borderColor = (valueMobile = "#2563EB", valueTablet = "#2563EB", valueDesktop = "#2563EB"))
                .prop(radiusTopStart = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .prop(radiusTopEnd = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .prop(radiusBottomStart = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .prop(radiusBottomEnd = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .prop(fontSize = (valueMobile = "20", valueTablet = "20", valueDesktop = "20"))
                .data(text = decreaseButton)
                .data(enable = enable)
                .action(event = "onClick") {
                    trigger(keyType = "nativeblocks/change_variable", name = "decrease", version = 1)
                    .prop(variableValue = "#SCRIPT
                    const count = {var:count}
                    let result = count
                    if (count >= 1) {
                        result = count - 1
                    } else {
                        result = count
                    }
                    result
                    #ENDSCRIPT")
                    .data(variableKey = count)
                }
                block(keyType = "nativeblocks/text", key = "countText", visibility = visible, version = 1)
                .prop(fontSize = (valueMobile = "18", valueTablet = "18", valueDesktop = "18"))
                .prop(textAlign = (valueMobile = "center", valueTablet = "center", valueDesktop = "center"))
                .prop(width = (valueMobile = "128", valueTablet = "128", valueDesktop = "128"))
                .data(text = count)
                block(keyType = "nativeblocks/button", key = "increaseButton", visibility = visible, version = 1)
                .prop(backgroundColor = (valueMobile = "#2563EB", valueTablet = "#2563EB", valueDesktop = "#2563EB"))
                .prop(borderColor = (valueMobile = "#2563EB", valueTablet = "#2563EB", valueDesktop = "#2563EB"))
                .prop(radiusTopStart = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .prop(radiusTopEnd = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .prop(radiusBottomStart = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .prop(radiusBottomEnd = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .prop(fontSize = (valueMobile = "20", valueTablet = "20", valueDesktop = "20"))
                .data(text = increaseButton)
                .data(enable = enable)
                .action(event = "onClick") {
                    trigger(keyType = "nativeblocks/change_variable", name = "increase", version = 1)
                    .prop(variableValue = "#SCRIPT
                    const count = {var:count}
                    const result = count + 1
                    result
                    #ENDSCRIPT")
                    .data(variableKey = count)
                }
            }
        }
    }
}`

	frameDSL, err := nbx.Parse(dsl)
	if err != nil {
		t.Fatalf("Failed to parse DSL: %v", err)
	}

	frameJson, err := compiler.ToJson(frameDSL, schema, "")
	if err != nil {
		t.Fatalf("%v", err)
	}

	if frameJson.Name != "welcome" {
		t.Errorf("Expected name 'welcome', got '%s'", frameJson.Name)
	}
	if frameJson.Route != "/welcome" {
		t.Errorf("Expected route '/welcome', got '%s'", frameJson.Route)
	}
	if len(frameJson.Variables) != 7 {
		t.Errorf("Expected 7 variables, got %d", len(frameJson.Variables))
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
