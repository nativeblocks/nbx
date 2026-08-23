package compiler

import (
	"os"
	"strings"
	"testing"

	"github.com/nativeblocks/nbx/internal/formatter"
	"github.com/nativeblocks/nbx/internal/lexer"
	"github.com/nativeblocks/nbx/internal/parser"
	"github.com/nativeblocks/nbx/internal/validator"
)

func TestToDsl(t *testing.T) {
	blocksJSON, err := os.ReadFile("../example/blocks.json")
	if err != nil {
		t.Fatalf("Failed to read blocks.json: %v", err)
	}

	actionsJSON, err := os.ReadFile("../example/actions.json")
	if err != nil {
		t.Fatalf("Failed to read actions.json: %v", err)
	}

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
        block(keyType = "nativeblocks/column", key = "mainColumn", visibility = visible, version = 2)
        .prop(
            horizontalAlignment = "centerHorizontally",
            width = "match",
            height = "match"
        )
        .slot("content") {
            block(keyType = "nativeblocks/spacer", key = "spacer1", visibility = visible, version = 2)
            block(keyType = "nativeblocks/column", key = "nativeblocksColumn", visibility = visible, version = 2)
            .prop(
                horizontalAlignment = "centerHorizontally",
                paddingTop = "64",
                verticalArrangement = "spaceAround"
            )
            .slot("content") {
                block(keyType = "nativeblocks/image", key = "logo", visibility = visible, version = 2)
                .prop(
                    scaleType = "inside",
                    width = "128",
                    height = "128"
                )
                .data(imageUrl = logo)

                block(keyType = "nativeblocks/text", key = "welcome", visibility = visible, version = 2)
                .prop(
                    fontSize = "24",
                    textAlign = "center",
                    width = "wrap"
                )
                .data(text = welcome)
            }
            block(keyType = "nativeblocks/row", key = "buttonsRow", visibility = visible, version = 2)
            .prop(
                horizontalArrangement = "spaceAround",
                verticalAlignment = "centerVertically",
                paddingTop = "12"
            )
            .slot("content") {
                block(keyType = "nativeblocks/button", key = "decreaseButton", visibility = visible, version = 2)
                .prop(
                    backgroundColor = "#2563EB",
                    borderColor = "#2563EB",
                    radiusTopStart = "32",
                    radiusTopEnd = "32",
                    radiusBottomStart = "32",
                    radiusBottomEnd = "32",
                    fontSize = "20"
                )
                .data(text = decreaseButton, enable = enable)
                .action(event = "onClick") {
                    trigger(keyType = "SCRIPT", name = "decrease", version = 2)
                    .prop(script = "
const count = getVariable("count")
let result = count
if (count >= 1) {
    result = count - 1
} else {
    result = count
}
updateVariable("count", result)
")
                }
                block(keyType = "nativeblocks/text", key = "countText", visibility = visible, version = 2)
                .prop(
                    fontSize = "18",
                    textAlign = "center",
                    width = "128"
                )
                .data(text = count)
                block(keyType = "nativeblocks/button", key = "increaseButton", visibility = visible, version = 2)
                .prop(
                    backgroundColor = "#2563EB",
                    borderColor = "#2563EB",
                    radiusTopStart = "32",
                    radiusTopEnd = "32",
                    radiusBottomStart = "32",
                    radiusBottomEnd = "32",
                    fontSize = "20"
                )
                .data(text = increaseButton, enable = enable)
                .action(event = "onClick") {
                    trigger(keyType = "SCRIPT", name = "increase", version = 2)
                    .prop(script = "
const count = getVariable("count")
const result = count + 1
updateVariable("count", result)
                    ")
                }
            }
            block(keyType = "nativeblocks/spacer", key = "spacer2", visibility = visible, version = 2)
        }
    }
}`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frameDSL := p.ParseNBX()

	errorCollector := p.ErrorCollector()
	if frameDSL == nil || errorCollector.HasErrors() {
		t.Fatalf("Failed to parse DSL: %v", errorCollector.FormatAll())
	}

	collector, _ := validator.ValidateWithSource(frameDSL, dsl)
	if collector != nil && collector.HasErrors() {
		t.Fatalf("Validation failed: %v", collector.FormatAll())
	}

	frameJson, nbxErrs := ToJson(*frameDSL, string(blocksJSON), string(actionsJSON), "")
	if nbxErrs != nil {
		t.Fatalf("%v", nbxErrs)
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

	l2 := lexer.NewLexer(invalidDsl)
	p2 := parser.NewParser(l2, invalidDsl)
	invalidFrameDSL := p2.ParseNBX()

	errorCollector2 := p2.ErrorCollector()
	if invalidFrameDSL == nil || errorCollector2.HasErrors() {
		t.Fatalf("Failed to parse invalid DSL: %v", errorCollector2.FormatAll())
	}

	collector2, _ := validator.ValidateWithSource(invalidFrameDSL, invalidDsl)
	if collector2 != nil && collector2.HasErrors() {
		t.Fatalf("Validation failed: %v", collector2.FormatAll())
	}

	_, nbxErrs = ToJson(*invalidFrameDSL, string(blocksJSON), string(actionsJSON), "")
	if nbxErrs != nil {
		t.Error("Expected no error but got error", nbxErrs)
	}

	customID := "custom-frame-id"
	frameJsonWithID, nbxErrs := ToJson(*frameDSL, string(blocksJSON), string(actionsJSON), customID)
	if nbxErrs != nil {
		t.Fatalf("Failed to convert to JSON with custom ID: %v", nbxErrs)
	}
	if frameJsonWithID.Id != customID {
		t.Errorf("Expected frame ID '%s', got '%s'", customID, frameJsonWithID.Id)
	}
}

func TestToString(t *testing.T) {
	content, err := os.ReadFile("../example/welcome_android.nbx")
	if err != nil {
		t.Fatal(err)
	}

	originalDSL := string(content)

	l := lexer.NewLexer(originalDSL)
	p := parser.NewParser(l, originalDSL)
	dslModel := p.ParseNBX()

	errorCollector := p.ErrorCollector()
	if dslModel == nil || errorCollector.HasErrors() {
		t.Fatal("Parse error:", errorCollector.FormatAll())
	}

	collector, _ := validator.ValidateWithSource(dslModel, originalDSL)
	if collector != nil && collector.HasErrors() {
		t.Fatal("Validation failed:", collector.FormatAll())
	}

	reconstructedDSL := formatter.FormatFrameDSL(*dslModel)

	t.Logf("=== ORIGINAL DSL ===\n%s", originalDSL)
	t.Logf("\n=== RECONSTRUCTED DSL ===\n%s", reconstructedDSL)

	if strings.TrimSpace(originalDSL) == strings.TrimSpace(reconstructedDSL) {
		t.Log("originalDsl and stringify dsl are identical")
	}

	l2 := lexer.NewLexer(reconstructedDSL)
	p2 := parser.NewParser(l2, reconstructedDSL)
	reconstructedModel := p2.ParseNBX()

	errorCollector2 := p2.ErrorCollector()
	if reconstructedModel == nil || errorCollector2.HasErrors() {
		t.Fatal("Reconstructed DSL cannot be parsed:", errorCollector2.FormatAll())
	}

	t.Log("Reconstructed DSL is valid and parseable!")
}
