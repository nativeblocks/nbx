package test

import (
	"strings"
	"testing"

	"github.com/nativeblocks/nbx/internal/lexer"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/nativeblocks/nbx/internal/parser"
	"github.com/nativeblocks/nbx/internal/validator"
)

func _countBlocksRecursive(blocks []model.BlockDSLModel) int {
	count := len(blocks)
	for _, block := range blocks {
		count += _countBlocksRecursive(block.Blocks)
	}
	return count
}

func _countActionsRecursive(blocks []model.BlockDSLModel) int {
	count := 0
	for _, block := range blocks {
		count += len(block.Actions)
		count += _countActionsRecursive(block.Blocks)
	}
	return count
}

func TestAllTypes(t *testing.T) {
	dsl := `
frame(name = "types_test", route = "/types") {
	var boolVar: BOOLEAN = true
	var intVar: INT = 42
	var longVar: LONG = 9223372036854775807
	var floatVar: FLOAT = 3.14
	var doubleVar: DOUBLE = 3.141592653589793
	var stringVar: STRING = "Hello World"
}
`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors: %s", errorCollector.FormatAll())
	}

	validatorCollector, _ := validator.ValidateWithSource(frame, dsl)
	if validatorCollector.HasErrors() {
		t.Fatalf("Validation errors: %s", validatorCollector.FormatAll())
	}

	if len(frame.Variables) != 6 {
		t.Errorf("Expected 6 variables, got %d", len(frame.Variables))
	}

	expectedTypes := map[string]string{
		"boolVar":   "BOOLEAN",
		"intVar":    "INT",
		"longVar":   "LONG",
		"floatVar":  "FLOAT",
		"doubleVar": "DOUBLE",
		"stringVar": "STRING",
	}

	for _, v := range frame.Variables {
		if expectedType, ok := expectedTypes[v.Key]; ok {
			if v.Type != expectedType {
				t.Errorf("Variable %s: expected type %s, got %s", v.Key, expectedType, v.Type)
			}
		}
	}
}

func TestBlockProperties(t *testing.T) {
	dsl := `
frame(name = "props_test", route = "/props") {
	var visible: BOOLEAN = true

	block(keyType = "ROOT", key = "root", visibility = visible)
	.prop(
		width = "match",
		height = "wrap"
	)
}
`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors: %s", errorCollector.FormatAll())
	}

	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(frame.Blocks))
	}

	block := frame.Blocks[0]
	if len(block.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(block.Properties))
	}
}

func TestNestedBlocks(t *testing.T) {
	dsl := `
frame(name = "nested_test", route = "/nested") {
	var visible: BOOLEAN = true

	block(keyType = "ROOT", key = "root", visibility = visible)
	.slot("content") {
		block(keyType = "nativeblocks/column", key = "column1", visibility = visible, version = 1)
		.slot("content") {
			block(keyType = "nativeblocks/text", key = "text1", visibility = visible, version = 1)
		}
	}
}
`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors: %s", errorCollector.FormatAll())
	}

	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 root block, got %d", len(frame.Blocks))
	}

	rootBlock := frame.Blocks[0]
	if len(rootBlock.Blocks) != 1 {
		t.Errorf("Expected 1 child block in root, got %d", len(rootBlock.Blocks))
	}

	column1 := rootBlock.Blocks[0]
	if column1.Key != "column1" {
		t.Errorf("Expected column1, got %s", column1.Key)
	}
	if len(column1.Blocks) != 1 {
		t.Errorf("Expected 1 block in column1, got %d", len(column1.Blocks))
	}
}

func TestActions(t *testing.T) {
	dsl := `
frame(name = "actions_test", route = "/actions") {
	var count: INT = 0
	var buttonText: STRING = "Click"
	var enabled: BOOLEAN = true

	block(keyType = "ROOT", key = "root")
	.slot("content") {
		block(keyType = "nativeblocks/button", key = "incrementBtn", version = 1)
		.data(text = buttonText, enable = enabled)
		.action(event = "onClick") {
			trigger(keyType = "nativeblocks/change_variable", name = "increment", version = 1)
			.data(variableKey = count)
			.prop(variableValue = "1")
		}
	}
}
`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors: %s", errorCollector.FormatAll())
	}

	rootBlock := frame.Blocks[0]
	button := rootBlock.Blocks[0]

	if len(button.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(button.Actions))
	}

	action := button.Actions[0]
	if action.Event != "onClick" {
		t.Errorf("Expected onClick event, got %s", action.Event)
	}

	if len(action.Triggers) != 1 {
		t.Errorf("Expected 1 trigger, got %d", len(action.Triggers))
	}

	trigger := action.Triggers[0]
	if trigger.KeyType != "nativeblocks/change_variable" {
		t.Errorf("Expected nativeblocks/change_variable keyType, got %s", trigger.KeyType)
	}
}

func TestTriggerChaining(t *testing.T) {
	dsl := `
frame(name = "chaining_test", route = "/chaining") {
	var count: INT = 0
	var buttonText: STRING = "Click"
	var enabled: BOOLEAN = true

	block(keyType = "ROOT", key = "root")
	.slot("content") {
		block(keyType = "nativeblocks/button", key = "btn", version = 1)
		.data(text = buttonText, enable = enabled)
		.action(event = "onClick") {
			trigger(keyType = "nativeblocks/change_variable", name = "first", version = 1)
			.data(variableKey = count)
			.prop(variableValue = "1")
			.then("SUCCESS") {
				trigger(keyType = "nativeblocks/change_variable", name = "second", version = 1)
				.data(variableKey = count)
				.prop(variableValue = "2")
			}
		}
	}
}
`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors: %s", errorCollector.FormatAll())
	}

	button := frame.Blocks[0].Blocks[0]
	action := button.Actions[0]
	trigger := action.Triggers[0]

	if trigger.Name != "first" {
		t.Errorf("Expected first trigger, got %s", trigger.Name)
	}

	if len(trigger.Triggers) != 1 {
		t.Errorf("Expected 1 nested trigger, got %d", len(trigger.Triggers))
	}

	nestedTrigger := trigger.Triggers[0]
	if nestedTrigger.Name != "second" {
		t.Errorf("Expected second trigger, got %s", nestedTrigger.Name)
	}
}

func TestScripts(t *testing.T) {
	dsl := `
frame(name = "script_test", route = "/script") {
	var count: INT = 0
	var buttonText: STRING = "Click"
	var enabled: BOOLEAN = true

	block(keyType = "ROOT", key = "root")
	.slot("content") {
		block(keyType = "nativeblocks/button", key = "btn", version = 1)
		.data(text = buttonText, enable = enabled)
		.action(event = "onClick") {
			trigger(keyType = "nativeblocks/change_variable", name = "calculate", version = 1)
			.prop(variableValue = "#SCRIPT
			const count = {var:count}
			const result = count + 1
			result
			#ENDSCRIPT")
			.data(variableKey = count)
		}
	}
}
`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors: %s", errorCollector.FormatAll())
	}

	button := frame.Blocks[0].Blocks[0]
	trigger := button.Actions[0].Triggers[0]

	if len(trigger.Properties) != 1 {
		t.Errorf("Expected 1 property with script, got %d", len(trigger.Properties))
	}

	prop := trigger.Properties[0]
	if !strings.Contains(prop.Value, "const count = {var:count}") {
		t.Error("Script content not preserved correctly")
	}
}

func TestComplexExample(t *testing.T) {
	dsl := `
frame(name = "welcome", route = "/welcome") {
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
		.prop(
			horizontalAlignment = "centerHorizontally",
			width = "match",
			height = "match"
		)
		.slot("content") {
			block(keyType = "nativeblocks/image", key = "logo", visibility = visible, version = 1)
			.prop(
				scaleType = "inside",
				width = "128",
				height = "128"
			)
			.data(imageUrl = logo)

			block(keyType = "nativeblocks/text", key = "welcome", visibility = visible, version = 1)
			.prop(
				fontSize = "24",
				textAlign = "center",
				width = "wrap"
			)
			.data(text = welcome)

			block(keyType = "nativeblocks/row", key = "buttonsRow", visibility = visible, version = 1)
			.prop(
				horizontalArrangement = "spaceAround",
				verticalAlignment = "centerVertically",
				paddingTop = "12"
			)
			.slot("content") {
				block(keyType = "nativeblocks/button", key = "decreaseButton", visibility = visible, version = 1)
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
				.prop(
					fontSize = "18",
					textAlign = "center",
					width = "128"
				)
				.data(text = count)

				block(keyType = "nativeblocks/button", key = "increaseButton", visibility = visible, version = 1)
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
}
`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors: %s", errorCollector.FormatAll())
	}

	validatorCollector, _ := validator.ValidateWithSource(frame, dsl)
	if validatorCollector.HasErrors() {
		t.Fatalf("Validation errors: %s", validatorCollector.FormatAll())
	}

	if frame.Name != "welcome" {
		t.Errorf("Expected frame name 'welcome', got '%s'", frame.Name)
	}
	if frame.Route != "/welcome" {
		t.Errorf("Expected route '/welcome', got '%s'", frame.Route)
	}

	if len(frame.Variables) != 7 {
		t.Errorf("Expected 7 variables, got %d", len(frame.Variables))
	}

	if validatorCollector.HasWarnings() {
		t.Logf("Warnings (should be none): %s", validatorCollector.FormatAll())
	}

	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 root block, got %d", len(frame.Blocks))
	}

	rootBlock := frame.Blocks[0]
	if rootBlock.Key != "root" {
		t.Errorf("Expected root block key 'root', got '%s'", rootBlock.Key)
	}

	totalBlocks := _countBlocksRecursive(frame.Blocks)
	if totalBlocks < 7 {
		t.Errorf("Expected at least 7 blocks in total, got %d", totalBlocks)
	}

	actionCount := _countActionsRecursive(frame.Blocks)
	if actionCount != 2 {
		t.Errorf("Expected 2 actions (increase/decrease), got %d", actionCount)
	}
}

func TestDataBindings(t *testing.T) {
	dsl := `
frame(name = "bindings_test", route = "/bindings") {
	var text: STRING = "Hello"
	var imageUrl: STRING = "https://example.com/image.png"

	block(keyType = "ROOT", key = "root")
	.slot("content") {
		block(keyType = "nativeblocks/image", key = "img", version = 1)
		.data(imageUrl = imageUrl)
	}
}
`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors: %s", errorCollector.FormatAll())
	}

	imageBlock := frame.Blocks[0].Blocks[0]
	if len(imageBlock.Data) != 1 {
		t.Errorf("Expected 1 data binding, got %d", len(imageBlock.Data))
	}

	if imageBlock.Data[0].Key != "imageUrl" {
		t.Errorf("Expected data key 'imageUrl', got '%s'", imageBlock.Data[0].Key)
	}

	if imageBlock.Data[0].Value != "imageUrl" {
		t.Errorf("Expected data value to be variable reference 'imageUrl', got '%s'", imageBlock.Data[0].Value)
	}
}

func TestTypeSafety(t *testing.T) {
	validDSL := `
frame(name = "type_safe", route = "/safe") {
	var myBool: BOOLEAN = true
	var myInt: INT = 42
	var myLong: LONG = 9223372036854775807
	var myFloat: FLOAT = 3.14
	var myDouble: DOUBLE = 3.141592653589793
	var myString: STRING = "Hello"
}
`

	l := lexer.NewLexer(validDSL)
	p := parser.NewParser(l, validDSL)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse valid DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors on valid DSL: %s", errorCollector.FormatAll())
	}

	validatorCollector, _ := validator.ValidateWithSource(frame, validDSL)
	if validatorCollector.HasErrors() {
		t.Fatalf("Validation errors on valid DSL: %s", validatorCollector.FormatAll())
	}

	if len(frame.Variables) != 6 {
		t.Errorf("Expected 6 valid type variables, got %d", len(frame.Variables))
	}
}

func TestBlockVersion(t *testing.T) {
	dsl := `
frame(name = "version_test", route = "/version") {
	var visible: BOOLEAN = true

	block(keyType = "ROOT", key = "root")
	.slot("content") {
		block(keyType = "nativeblocks/column", key = "v1", visibility = visible, version = 1)
		block(keyType = "nativeblocks/column", key = "v2", visibility = visible, version = 2)
		block(keyType = "nativeblocks/column", key = "v10", visibility = visible, version = 10)
	}
}
`

	l := lexer.NewLexer(dsl)
	p := parser.NewParser(l, dsl)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatal("Failed to parse DSL")
	}

	errorCollector := p.ErrorCollector()
	if errorCollector.HasErrors() {
		t.Fatalf("Parse errors: %s", errorCollector.FormatAll())
	}

	rootBlock := frame.Blocks[0]
	if len(rootBlock.Blocks) != 3 {
		t.Fatalf("Expected 3 blocks, got %d", len(rootBlock.Blocks))
	}

	expectedVersions := []int{1, 2, 10}
	for i, block := range rootBlock.Blocks {
		if block.IntegrationVersion != expectedVersions[i] {
			t.Errorf("Block %d: expected version %d, got %d", i, expectedVersions[i], block.IntegrationVersion)
		}
	}
}

func TestErrorCases(t *testing.T) {
	testCases := []struct {
		name          string
		dsl           string
		expectError   bool
		errorContains string
	}{
		{
			name: "Duplicate variable",
			dsl: `
frame(name = "test", route = "/test") {
	var count: INT = 0
	var count: INT = 1
}
`,
			expectError:   true,
			errorContains: "Duplicate declaration",
		},
		{
			name: "Undefined variable reference",
			dsl: `
frame(name = "test", route = "/test") {
	var count: INT = 0
	block(keyType = "ROOT", key = "root", visibility = undefined_var)
}
`,
			expectError:   true,
			errorContains: "Undefined variable",
		},
		{
			name: "Invalid variable type",
			dsl: `
frame(name = "test", route = "/test") {
	var count: INVALID_TYPE = 0
}
`,
			expectError:   true,
			errorContains: "Unknown type",
		},
		{
			name: "Invalid variable value",
			dsl: `
frame(name = "test", route = "/test") {
	var count: INT = not_a_number
}
`,
			expectError:   true,
			errorContains: "Invalid initial value",
		},
		{
			name: "Duplicate block key",
			dsl: `
frame(name = "test", route = "/test") {
	block(keyType = "ROOT", key = "root")
	block(keyType = "ROOT", key = "root")
}
`,
			expectError:   true,
			errorContains: "Duplicate block key",
		},
		{
			name: "Type mismatch - INT assigned to STRING",
			dsl: `
frame(name = "test", route = "/test") {
	var text: STRING = 42
}
`,
			expectError:   true,
			errorContains: "numeric value",
		},
		{
			name: "Type mismatch - BOOLEAN assigned to STRING",
			dsl: `
frame(name = "test", route = "/test") {
	var text: STRING = true
}
`,
			expectError:   true,
			errorContains: "boolean literal",
		},
		{
			name: "Type mismatch - STRING assigned to INT",
			dsl: `
frame(name = "test", route = "/test") {
	var count: INT = "hello"
}
`,
			expectError:   true,
			errorContains: "not a valid integer",
		},
		{
			name: "Type mismatch - INT assigned to FLOAT (no implicit conversion)",
			dsl: `
frame(name = "test", route = "/test") {
	var value: FLOAT = 42
}
`,
			expectError:   true,
			errorContains: "integer. FLOAT requires decimal point",
		},
		{
			name: "Type mismatch - INT assigned to DOUBLE (no implicit conversion)",
			dsl: `
frame(name = "test", route = "/test") {
	var value: DOUBLE = 123
}
`,
			expectError:   true,
			errorContains: "integer. DOUBLE requires decimal point",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := lexer.NewLexer(tc.dsl)
			p := parser.NewParser(l, tc.dsl)
			frame := p.ParseNBX()

			errorCollector := p.ErrorCollector()

			hasErrors := errorCollector.HasErrors()
			errorMsg := ""

			if frame != nil {
				validatorCollector, _ := validator.ValidateWithSource(frame, tc.dsl)
				if validatorCollector.HasErrors() {
					hasErrors = true
					errorMsg = validatorCollector.FormatAll()
				}
			} else {
				errorMsg = errorCollector.FormatAll()
			}

			if tc.expectError && !hasErrors {
				t.Errorf("Expected error containing '%s', but got none", tc.errorContains)
			}

			if tc.expectError && !strings.Contains(errorMsg, tc.errorContains) {
				t.Errorf("Expected error containing '%s', got: %s", tc.errorContains, errorMsg)
			}

			if tc.expectError && hasErrors {
				if tc.name != "Undefined variable reference" && !strings.Contains(errorMsg, "line") {
					t.Error("Error message should contain line information")
				}
				t.Logf("Error with line tracking:\n%s", errorMsg)
			}
		})
	}
}
