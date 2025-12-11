package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Type interface {
	Name() string
	IsCompatible(other Type) bool
}

type PrimitiveType struct {
	name string
}

func (t PrimitiveType) Name() string {
	return t.name
}

func (t PrimitiveType) IsCompatible(other Type) bool {
	if t.name == other.Name() {
		return true
	}

	switch t.name {
	case "INT":
		return other.Name() == "LONG"
	case "FLOAT":
		return other.Name() == "DOUBLE"
	case "LONG":
		return false
	case "DOUBLE":
		return false
	default:
		return false
	}
}

var (
	TypeString  = PrimitiveType{"STRING"}
	TypeInt     = PrimitiveType{"INT"}
	TypeLong    = PrimitiveType{"LONG"}
	TypeFloat   = PrimitiveType{"FLOAT"}
	TypeDouble  = PrimitiveType{"DOUBLE"}
	TypeBoolean = PrimitiveType{"BOOLEAN"}
	TypeUnknown = PrimitiveType{"UNKNOWN"}
)

func FromString(typeName string) (Type, error) {
	switch strings.ToUpper(typeName) {
	case "STRING":
		return TypeString, nil
	case "INT":
		return TypeInt, nil
	case "LONG":
		return TypeLong, nil
	case "FLOAT":
		return TypeFloat, nil
	case "DOUBLE":
		return TypeDouble, nil
	case "BOOLEAN":
		return TypeBoolean, nil
	default:
		return TypeUnknown, fmt.Errorf("unknown type: %s", typeName)
	}
}

func InferType(value string) Type {
	value = strings.TrimSpace(value)

	if value == "true" || value == "false" {
		return TypeBoolean
	}

	if _isInteger(value) {
		if len(value) > 9 {
			return TypeLong
		}
		return TypeInt
	}

	if _isFloat(value) {
		parts := strings.Split(value, ".")
		if len(parts) == 2 && len(parts[1]) > 6 {
			return TypeDouble
		}
		return TypeFloat
	}

	return TypeString
}

func ValidateValue(value string, expectedType Type) (bool, string) {
	value = strings.TrimSpace(value)

	switch expectedType.Name() {
	case "BOOLEAN":
		if value != "true" && value != "false" {
			return false, fmt.Sprintf("'%s' is not a valid boolean. Expected 'true' or 'false'", value)
		}
		return true, ""

	case "INT":
		if !_isInteger(value) {
			return false, fmt.Sprintf("'%s' is not a valid integer", value)
		}
		if _, err := strconv.ParseInt(value, 10, 32); err != nil {
			return false, fmt.Sprintf("'%s' is out of range for INT. Consider using LONG", value)
		}
		if len(value) > 9 {
			return false, fmt.Sprintf("'%s' might be too large for INT (max 9 digits). Consider using LONG", value)
		}
		return true, ""

	case "LONG":
		if !_isInteger(value) {
			return false, fmt.Sprintf("'%s' is not a valid integer", value)
		}
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return false, fmt.Sprintf("'%s' is out of range for LONG", value)
		}
		return true, ""

	case "FLOAT":
		if !_isFloat(value) {
			if _isInteger(value) {
				return false, fmt.Sprintf("'%s' is an integer. FLOAT requires decimal point (e.g., %s.0)", value, value)
			}
			return false, fmt.Sprintf("'%s' is not a valid FLOAT value. FLOAT requires decimal point", value)
		}
		if _, err := strconv.ParseFloat(value, 32); err != nil {
			return false, fmt.Sprintf("'%s' is out of range for FLOAT", value)
		}
		return true, ""

	case "DOUBLE":
		if !_isFloat(value) {
			if _isInteger(value) {
				return false, fmt.Sprintf("'%s' is an integer. DOUBLE requires decimal point (e.g., %s.0)", value, value)
			}
			return false, fmt.Sprintf("'%s' is not a valid DOUBLE value. DOUBLE requires decimal point", value)
		}
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return false, fmt.Sprintf("'%s' is out of range for DOUBLE", value)
		}
		return true, ""

	case "STRING":
		if _isInteger(value) {
			return false, fmt.Sprintf("'%s' is a numeric value. Cannot assign to STRING type. Use quotes for string values", value)
		}
		if _isFloat(value) {
			return false, fmt.Sprintf("'%s' is a numeric value. Cannot assign to STRING type. Use quotes for string values", value)
		}
		if value == "true" || value == "false" {
			return false, fmt.Sprintf("'%s' is a boolean literal. Cannot assign to STRING type. Use quotes for string values", value)
		}
		return true, ""

	default:
		return false, fmt.Sprintf("unknown type: %s", expectedType.Name())
	}
}

var (
	integerRegex = regexp.MustCompile(`^-?\d+$`)
	floatRegex   = regexp.MustCompile(`^-?\d+\.\d+$`)
)

func _isInteger(s string) bool {
	return integerRegex.MatchString(s)
}

func _isFloat(s string) bool {
	return floatRegex.MatchString(s)
}
