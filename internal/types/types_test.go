package types

import (
	"testing"
)

func TestInferType(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected Type
	}{
		{"boolean true", "true", TypeBoolean},
		{"boolean false", "false", TypeBoolean},
		{"small integer", "42", TypeInt},
		{"large integer", "1234567890", TypeLong},
		{"float", "3.14", TypeFloat},
		{"double high precision", "3.14159265358979", TypeDouble},
		{"string", "hello", TypeString},
		{"quoted string", "\"hello\"", TypeString},
		{"empty string", "", TypeString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InferType(tt.value)
			if result.Name() != tt.expected.Name() {
				t.Errorf("InferType(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestTypeCompatibility(t *testing.T) {
	tests := []struct {
		name       string
		sourceType Type
		targetType Type
		compatible bool
	}{
		{"INT to INT", TypeInt, TypeInt, true},
		{"INT to LONG", TypeInt, TypeLong, true},
		{"LONG to INT", TypeLong, TypeInt, false},
		{"FLOAT to DOUBLE", TypeFloat, TypeDouble, true},
		{"DOUBLE to FLOAT", TypeDouble, TypeFloat, false},
		{"STRING to STRING", TypeString, TypeString, true},
		{"BOOLEAN to STRING", TypeBoolean, TypeString, false},
		{"INT to FLOAT", TypeInt, TypeFloat, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sourceType.IsCompatible(tt.targetType)
			if result != tt.compatible {
				t.Errorf("%s.IsCompatible(%s) = %v, want %v",
					tt.sourceType, tt.targetType, result, tt.compatible)
			}
		})
	}
}

func TestValidateValue(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		expectedType Type
		valid        bool
	}{
		{"valid boolean true", "true", TypeBoolean, true},
		{"valid boolean false", "false", TypeBoolean, true},
		{"invalid boolean", "yes", TypeBoolean, false},
		{"valid int", "42", TypeInt, true},
		{"invalid int", "abc", TypeInt, false},
		{"int too large", "12345678901", TypeInt, false},
		{"valid long", "12345678901", TypeLong, true},
		{"invalid long", "abc", TypeLong, false},
		{"valid float", "3.14", TypeFloat, true},
		{"invalid float integer (strict typing)", "42", TypeFloat, false}, // Strict: no implicit INT to FLOAT conversion
		{"invalid float", "abc", TypeFloat, false},
		{"valid double", "3.14159", TypeDouble, true},
		{"invalid double integer (strict typing)", "123", TypeDouble, false}, // Strict: no implicit INT to DOUBLE conversion
		{"invalid double", "abc", TypeDouble, false},
		{"valid string", "hello", TypeString, true},
		{"empty string", "", TypeString, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, _ := ValidateValue(tt.value, tt.expectedType)
			if valid != tt.valid {
				t.Errorf("ValidateValue(%q, %s) = %v, want %v",
					tt.value, tt.expectedType, valid, tt.valid)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name        string
		typeName    string
		expectedErr bool
	}{
		{"valid STRING", "STRING", false},
		{"valid INT", "INT", false},
		{"valid LONG", "LONG", false},
		{"valid FLOAT", "FLOAT", false},
		{"valid DOUBLE", "DOUBLE", false},
		{"valid BOOLEAN", "BOOLEAN", false},
		{"valid lowercase", "string", false},
		{"invalid type", "INVALID", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromString(tt.typeName)
			if (err != nil) != tt.expectedErr {
				t.Errorf("FromString(%q) error = %v, expectedErr %v", tt.typeName, err, tt.expectedErr)
			}
		})
	}
}
