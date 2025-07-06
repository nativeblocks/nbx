package compiler

import (
	"fmt"
	"strings"
)

type nbxValidationError struct {
	Description string      `json:"description"`
	Field       string      `json:"field"`
	Type        string      `json:"type"`
	Value       interface{} `json:"value"`
}

func formatErrors(errs []nbxValidationError) []string {
	out := make([]string, 0, len(errs))
	for _, e := range errs {
		suffix := e.Description
		if strings.HasPrefix(suffix, e.Field) {
			suffix = strings.TrimPrefix(suffix, e.Field)
		} else if idx := strings.Index(suffix, " "); idx != -1 {
			suffix = suffix[idx:]
		}
		out = append(out, fmt.Sprintf("%s%s", e.Value, suffix))
	}
	return out
}
