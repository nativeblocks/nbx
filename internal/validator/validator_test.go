package validator

import (
	"testing"

	"github.com/nativeblocks/nbx/internal/model"
)

func TestValidateVariables(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Variables: []model.VariableDSLModel{
			{Key: "count", Type: "INT", Value: "0"},
			{Key: "name", Type: "STRING", Value: "test"},
			{Key: "enabled", Type: "BOOLEAN", Value: "true"},
		},
		Blocks: []model.BlockDSLModel{},
	}

	collector, err := Validate(frame)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if collector.HasErrors() {
		t.Errorf("Expected no errors, got: %v", collector.Errors())
	}
}

func TestDuplicateVariables(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Variables: []model.VariableDSLModel{
			{Key: "count", Type: "INT", Value: "0"},
			{Key: "count", Type: "INT", Value: "1"}, // Duplicate
		},
		Blocks: []model.BlockDSLModel{},
	}

	collector, _ := Validate(frame)

	if !collector.HasErrors() {
		t.Error("Expected error for duplicate variable")
	}

	if len(collector.Errors()) != 1 {
		t.Errorf("Expected 1 error, got %d", len(collector.Errors()))
	}
}

func TestInvalidVariableType(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Variables: []model.VariableDSLModel{
			{Key: "count", Type: "INVALID_TYPE", Value: "0"},
		},
		Blocks: []model.BlockDSLModel{},
	}

	collector, _ := Validate(frame)

	if !collector.HasErrors() {
		t.Error("Expected error for invalid type")
	}
}

func TestInvalidVariableValue(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Variables: []model.VariableDSLModel{
			{Key: "count", Type: "INT", Value: "not_a_number"},
		},
		Blocks: []model.BlockDSLModel{},
	}

	collector, _ := Validate(frame)

	if !collector.HasErrors() {
		t.Error("Expected error for invalid value")
	}
}

func TestUndefinedVariableReference(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Variables: []model.VariableDSLModel{
			{Key: "visible", Type: "BOOLEAN", Value: "true"},
		},
		Blocks: []model.BlockDSLModel{
			{
				KeyType:       "ROOT",
				Key:           "root",
				VisibilityKey: "invisble", // Typo - should be "visible"
			},
		},
	}

	collector, _ := Validate(frame)

	if !collector.HasErrors() {
		t.Error("Expected error for undefined variable reference")
	}

	// Check if error suggests the correct variable
	errors := collector.Errors()
	if len(errors) > 0 {
		// The error should suggest "visible" as similar
		// This is a basic check - could be more sophisticated
		if errors[0].Suggestion == "" {
			t.Error("Expected suggestion for undefined variable")
		}
	}
}

func TestDuplicateBlockKeys(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Blocks: []model.BlockDSLModel{
			{KeyType: "ROOT", Key: "root"},
			{KeyType: "COLUMN", Key: "root"}, // Duplicate key
		},
	}

	collector, _ := Validate(frame)

	if !collector.HasErrors() {
		t.Error("Expected error for duplicate block key")
	}
}

func TestMissingFrameAttributes(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "", // Missing
		Route: "/test",
		Type:  "FRAME",
	}

	collector, _ := Validate(frame)

	if !collector.HasErrors() {
		t.Error("Expected error for missing frame name")
	}
}

func TestUnusedVariable(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Variables: []model.VariableDSLModel{
			{Key: "unused", Type: "INT", Value: "0"},
			{Key: "used", Type: "BOOLEAN", Value: "true"},
		},
		Blocks: []model.BlockDSLModel{
			{
				KeyType:       "ROOT",
				Key:           "root",
				VisibilityKey: "used", // Only "used" is referenced
			},
		},
	}

	collector, _ := Validate(frame)

	if !collector.HasWarnings() {
		t.Error("Expected warning for unused variable")
	}

	warnings := collector.Warnings()
	if len(warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(warnings))
	}
}

func TestValidateDataBinding(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Variables: []model.VariableDSLModel{
			{Key: "message", Type: "STRING", Value: "hello"},
		},
		Blocks: []model.BlockDSLModel{
			{
				KeyType: "TEXT",
				Key:     "text1",
				Data: []model.BlockDataDSLModel{
					{Key: "text", Value: "message"}, // Variable reference
				},
			},
		},
	}

	collector, _ := Validate(frame)

	if collector.HasErrors() {
		t.Errorf("Expected no errors, got: %v", collector.Errors())
	}
}

func TestValidateTrigger(t *testing.T) {
	frame := &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Variables: []model.VariableDSLModel{
			{Key: "count", Type: "INT", Value: "0"},
		},
		Blocks: []model.BlockDSLModel{
			{
				KeyType: "BUTTON",
				Key:     "btn1",
				Actions: []model.ActionDSLModel{
					{
						Event: "onClick",
						Triggers: []model.ActionTriggerDSLModel{
							{
								KeyType: "CHANGE_VARIABLE",
								Name:    "increment",
								Then:    "NEXT",
								Data: []model.TriggerDataDSLModel{
									{Key: "variableKey", Value: "count"},
								},
							},
						},
					},
				},
			},
		},
	}

	collector, _ := Validate(frame)

	if collector.HasErrors() {
		t.Errorf("Expected no errors, got: %v", collector.Errors())
	}
}

func TestIsVariableName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"count", true},
		{"myVariable", true},
		{"var_name", true},
		{"123", false},
		{"1var", false},
		{"", false},
		{"my-var", false},
		{"my var", false},
	}

	for _, tt := range tests {
		result := _isVariableName(tt.input)
		if result != tt.expected {
			t.Errorf("_isVariableName(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
