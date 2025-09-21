package formatter

import (
	"strings"
	"testing"

	"github.com/nativeblocks/nbx/internal/model"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple frame with basic block",
			input: `frame(
    name = "test",
    route = "/test"
) {
    block(keyType = "ROOT", key = "root")
}`,
			expected: `frame(
    name = "test",
    route = "/test"
) {
    block(keyType = "ROOT", key = "root")
}`,
		},
		{
			name: "frame with variables",
			input: `frame(
    name = "login",
    route = "/login"
) {
    var visible: BOOLEAN = true
    var username: STRING = ""
    var count: INT = 42

    block(keyType = "ROOT", key = "root")
}`,
			expected: `frame(
    name = "login",
    route = "/login"
) {
    var visible: BOOLEAN = true
    var username: STRING = ""
    var count: INT = 42

    block(keyType = "ROOT", key = "root")
}`,
		},
		{
			name: "block with properties",
			input: `frame(
    name = "test",
    route = "/test"
) {
    block(keyType = "TEXT", key = "title")
    .prop(
        text = "Hello World",
        color = "#000000"
    )
}`,
			expected: `frame(
    name = "test",
    route = "/test"
) {
    block(keyType = "TEXT", key = "title")
    .prop(
        text = "Hello World",
        color = "#000000"
    )
}`,
		},
		{
			name: "block with visibility and version",
			input: `frame(
    name = "test",
    route = "/test"
) {
    block(keyType = "BUTTON", key = "btn", visibility = visible, version = 2)
    .prop(
        title = "Click me"
    )
}`,
			expected: `frame(
    name = "test",
    route = "/test"
) {
    block(keyType = "BUTTON", key = "btn", visibility = visible, version = 2)
    .prop(
        title = "Click me"
    )
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Format(tt.input)
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			// Normalize whitespace for comparison
			expected := strings.TrimSpace(tt.expected)
			actual := strings.TrimSpace(result)

			if actual != expected {
				t.Errorf("Format() mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
			}
		})
	}
}

func TestFormat_Error(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid syntax",
			input: "invalid dsl syntax",
		},
		{
			name:  "incomplete frame",
			input: "frame(",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Format(tt.input)
			if err == nil {
				t.Errorf("Format() expected error for invalid input, got nil")
			}
		})
	}
}

func TestFormatFrameDSL(t *testing.T) {
	frame := model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Variables: []model.VariableDSLModel{
			{Key: "visible", Type: "BOOLEAN", Value: "true"},
			{Key: "title", Type: "STRING", Value: "Hello"},
		},
		Blocks: []model.BlockDSLModel{
			{
				KeyType: "ROOT",
				Key:     "root",
			},
		},
	}

	result := FormatFrameDSL(frame)
	expected := `frame(
    name = "test",
    route = "/test"
) {
    var visible: BOOLEAN = true
    var title: STRING = "Hello"

    block(keyType = "ROOT", key = "root")
}`

	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("FormatFrameDSL() mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, result)
	}
}

func TestFormatVariableValueConsistent(t *testing.T) {
	tests := []struct {
		value     string
		valueType string
		expected  string
	}{
		{"hello", "STRING", `"hello"`},
		{"true", "BOOLEAN", "true"},
		{"42", "INT", "42"},
		{"100", "LONG", "100"},
		{"3.14", "FLOAT", "3.14"},
		{"2.718", "DOUBLE", "2.718"},
		{"custom", "CUSTOM", `"custom"`},
	}

	for _, tt := range tests {
		t.Run(tt.valueType, func(t *testing.T) {
			result := formatVariableValueConsistent(tt.value, tt.valueType)
			if result != tt.expected {
				t.Errorf("formatVariableValueConsistent(%q, %q) = %q, expected %q",
					tt.value, tt.valueType, result, tt.expected)
			}
		})
	}
}

func TestGetSinglePropertyValueConsistent(t *testing.T) {
	tests := []struct {
		name     string
		prop     model.BlockPropertyDSLModel
		expected string
	}{
		{
			name: "mobile value first",
			prop: model.BlockPropertyDSLModel{
				ValueMobile:  "mobile",
				ValueTablet:  "tablet",
				ValueDesktop: "desktop",
			},
			expected: "mobile",
		},
		{
			name: "tablet value when mobile empty",
			prop: model.BlockPropertyDSLModel{
				ValueMobile:  "",
				ValueTablet:  "tablet",
				ValueDesktop: "desktop",
			},
			expected: "tablet",
		},
		{
			name: "desktop value when mobile and tablet empty",
			prop: model.BlockPropertyDSLModel{
				ValueMobile:  "",
				ValueTablet:  "",
				ValueDesktop: "desktop",
			},
			expected: "desktop",
		},
		{
			name: "empty when all values empty",
			prop: model.BlockPropertyDSLModel{
				ValueMobile:  "",
				ValueTablet:  "",
				ValueDesktop: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSinglePropertyValueConsistent(tt.prop)
			if result != tt.expected {
				t.Errorf("getSinglePropertyValueConsistent() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestFormatScriptBlock(t *testing.T) {
	tests := []struct {
		name            string
		script          string
		baseIndentLevel int
		expected        string
	}{
		{
			name:            "no script block",
			script:          "regular text",
			baseIndentLevel: 1,
			expected:        "regular text",
		},
		{
			name:            "simple script block",
			script:          "#SCRIPT console.log('hello'); #ENDSCRIPT",
			baseIndentLevel: 1,
			expected:        "#SCRIPT\n    console.log('hello');\n    #ENDSCRIPT",
		},
		{
			name:            "script with whitespace",
			script:          "#SCRIPT\n  console.log('hello');\n  \n#ENDSCRIPT",
			baseIndentLevel: 1,
			expected:        "#SCRIPT\n    console.log('hello');\n    #ENDSCRIPT",
		},
		{
			name:            "script with base indent level 2",
			script:          "#SCRIPT console.log('test'); #ENDSCRIPT",
			baseIndentLevel: 2,
			expected:        "#SCRIPT\n        console.log('test');\n        #ENDSCRIPT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatScriptBlock(tt.script, tt.baseIndentLevel)
			if result != tt.expected {
				t.Errorf("formatScriptBlock() mismatch:\nExpected:\n%s\n\nActual:\n%s", tt.expected, result)
			}
		})
	}
}

func TestComplexFrameFormatting(t *testing.T) {
	input := `frame(
    name = "welcome",
    route = "/welcome"
) {
    var visible: BOOLEAN = true
    var enable: BOOLEAN = true
    var count: INT = 0
    var decreaseButton: STRING = "-"
    var increaseButton: STRING = "+"
    var welcomeText: STRING = "Welcome to Nativeblocks"
    var logo: STRING = "https://nativeblocks.io/nativeblocks_logo.png"

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "nativeblocks/vstack", key = "outerVStack", visibility = visible, version = 1)
        .prop(
            width = "fill",
            height = "fill"
        )
        .slot("content") {
            block(keyType = "nativeblocks/spacer", key = "spacer1", visibility = visible, version = 1)
            block(keyType = "nativeblocks/vstack", key = "logoVStack", visibility = visible, version = 1)
            .prop(
                spacing = "16",
                alignmentHorizontal = "center"
            )
            .slot("content") {
                block(keyType = "nativeblocks/image", key = "logoImage", visibility = visible, version = 1)
                .prop(
                    width = "128",
                    height = "128",
                    contentMode = "fit"
                )
                .data(imageUrl = logo)
                block(keyType = "nativeblocks/text", key = "welcomeText", visibility = visible, version = 1)
                .prop(
                    fontSize = "24",
                    multilineTextAlignment = "center",
                    width = "fill"
                )
                .data(text = welcomeText)
            }
            block(keyType = "nativeblocks/hstack", key = "controlsHStack", visibility = visible, version = 1)
            .prop(
                spacing = "16",
                width = "fill",
                alignmentHorizontal = "center",
                paddingTop = "12",
                paddingBottom = "32"
            )
            .slot("content") {
                block(keyType = "nativeblocks/spacer", key = "spacer_control_1", visibility = visible, version = 1)
                block(keyType = "nativeblocks/button", key = "decreaseBtn", visibility = visible, version = 1)
                .prop(
                    width = "56",
                    height = "40",
                    backgroundColor = "#2563EB",
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

                block(keyType = "nativeblocks/text", key = "counterText", visibility = visible, version = 1)
                .prop(
                    fontSize = "18",
                    width = "128",
                    height = "40",
                    multilineTextAlignment = "center"
                )
                .data(text = count)
                block(keyType = "nativeblocks/button", key = "increaseBtn", visibility = visible, version = 1)
                .prop(
                    width = "56",
                    height = "40",
                    backgroundColor = "#2563EB",
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

                block(keyType = "nativeblocks/spacer", key = "spacer_control_2", visibility = visible, version = 1)
            }
            block(keyType = "nativeblocks/spacer", key = "spacer2", visibility = visible, version = 1)
        }
    }
}`

	result, err := Format(input)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	if !strings.Contains(result, `name = "welcome"`) {
		t.Error("Expected frame name to be preserved")
	}
	if !strings.Contains(result, `var visible: BOOLEAN = true`) {
		t.Error("Expected variable formatting to be preserved")
	}
	if !strings.Contains(result, `.slot("content")`) {
		t.Error("Expected slot formatting to be preserved")
	}
	if !strings.Contains(result, `.action(event = "onClick")`) {
		t.Error("Expected action formatting to be preserved")
	}
}
