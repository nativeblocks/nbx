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
        .assignProperty(horizontalAlignment = (valueMobile = "centerHorizontally", valueTablet = "centerHorizontally", valueDesktop = "centerHorizontally"))
        .assignProperty(width = (valueMobile = "match", valueTablet = "match", valueDesktop = "match"))
        .assignProperty(height = (valueMobile = "match", valueTablet = "match", valueDesktop = "match"))
        .slot("content") {
            block(keyType = "nativeblocks/column", key = "nativeblocksColumn", VISIBILITY = visible, version = 1)
            .assignProperty(horizontalAlignment = (valueMobile = "centerHorizontally", valueTablet = "centerHorizontally", valueDesktop = "centerHorizontally"))
            .assignProperty(paddingTop = (valueMobile = "64", valueTablet = "64", valueDesktop = "64"))
            .assignProperty(weight = (valueMobile = "0.4f", valueTablet = "0.4f", valueDesktop = "0.4f"))
            .assignProperty(verticalArrangement = (valueMobile = "spaceAround", valueTablet = "spaceAround", valueDesktop = "spaceAround"))
            .slot("content") {
                block(keyType = "nativeblocks/image", key = "logo", visibility = visible, version = 1)
                .assignProperty(scaleType = (valueMobile = "inside", valueTablet = "inside", valueDesktop = "inside"))
                .assignProperty(width = (valueMobile = "128", valueTablet = "128", valueDesktop = "128"))
                .assignProperty(height = (valueMobile = "128", valueTablet = "128", valueDesktop = "128"))
                .assignData(imageUrl = logo)

                block(keyType = "nativeblocks/text", key = "welcome", visibility = visible, version = 1)
                .assignProperty(fontSize = (valueMobile = "24", valueTablet = "24", valueDesktop = "24"))
                .assignProperty(textAlign = (valueMobile = "center", valueTablet = "center", valueDesktop = "center"))
                .assignProperty(width = (valueMobile = "wrap", valueTablet = "wrap", valueDesktop = "wrap"))
                .assignData(text = welcome)
            }
            block(keyType = "nativeblocks/row", key = "buttonsRow", visibility = visible, version = 1)
            .assignProperty(horizontalArrangement = (valueMobile = "spaceAround", valueTablet = "spaceAround", valueDesktop = "spaceAround"))
            .assignProperty(verticalAlignment = (valueMobile = "centerVertically", valueTablet = "centerVertically", valueDesktop = "centerVertically"))
            .assignProperty(paddingTop = (valueMobile = "12", valueTablet = "12", valueDesktop = "12"))
            .assignProperty(weight = (valueMobile = "0.6f", valueTablet = "0.6f", valueDesktop = "0.6f"))
            .slot("content") {
                block(keyType = "nativeblocks/button", key = "decreaseButton", visibility = visible, version = 1)
                .assignProperty(backgroundColor = (valueMobile = "#2563EB", valueTablet = "#2563EB", valueDesktop = "#2563EB"))
                .assignProperty(borderColor = (valueMobile = "#2563EB", valueTablet = "#2563EB", valueDesktop = "#2563EB"))
                .assignProperty(radiusTopStart = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .assignProperty(radiusTopEnd = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .assignProperty(radiusBottomStart = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .assignProperty(radiusBottomEnd = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .assignProperty(fontSize = (valueMobile = "20", valueTablet = "20", valueDesktop = "20"))
                .assignData(text = decreaseButton)
                .assignData(enable = enable)
                .action(event = "onClick") {
                    trigger(keyType = "nativeblocks/change_variable", name = "decrease", version = 1)
                    .assignProperty(variableValue = "#SCRIPT
                    const count = {var:count}
                    let result = count
                    if (count >= 1) {
                        result = count - 1
                    } else {
                        result = count
                    }
                    result
                    #ENDSCRIPT")
                    .assignData(variableKey = count)
                }
                block(keyType = "nativeblocks/text", key = "countText", visibility = visible, version = 1)
                .assignProperty(fontSize = (valueMobile = "18", valueTablet = "18", valueDesktop = "18"))
                .assignProperty(textAlign = (valueMobile = "center", valueTablet = "center", valueDesktop = "center"))
                .assignProperty(width = (valueMobile = "128", valueTablet = "128", valueDesktop = "128"))
                .assignData(text = count)
                block(keyType = "nativeblocks/button", key = "increaseButton", visibility = visible, version = 1)
                .assignProperty(backgroundColor = (valueMobile = "#2563EB", valueTablet = "#2563EB", valueDesktop = "#2563EB"))
                .assignProperty(borderColor = (valueMobile = "#2563EB", valueTablet = "#2563EB", valueDesktop = "#2563EB"))
                .assignProperty(radiusTopStart = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .assignProperty(radiusTopEnd = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .assignProperty(radiusBottomStart = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .assignProperty(radiusBottomEnd = (valueMobile = "32", valueTablet = "32", valueDesktop = "32"))
                .assignProperty(fontSize = (valueMobile = "20", valueTablet = "20", valueDesktop = "20"))
                .assignData(text = increaseButton)
                .assignData(enable = enable)
                .action(event = "onClick") {
                    trigger(keyType = "nativeblocks/change_variable", name = "increase", version = 1)
                    .assignProperty(variableValue = "#SCRIPT
                    const count = {var:count}
                    const result = count + 1
                    result
                    #ENDSCRIPT")
                    .assignData(variableKey = count)
                }
            }
        }
    }
}
`

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
