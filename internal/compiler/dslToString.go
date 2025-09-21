package compiler

import (
	"github.com/nativeblocks/nbx/internal/formatter"
	"github.com/nativeblocks/nbx/internal/model"
)

// ToString converts a FrameDSLModel back to the original NBX DSL string format
// using the consistent formatter
func ToString(frame model.FrameDSLModel) string {
	return formatter.FormatFrameDSL(frame)
}
