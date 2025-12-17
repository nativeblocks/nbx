package compiler

import (
	"github.com/nativeblocks/nbx/internal/formatter"
	"github.com/nativeblocks/nbx/internal/model"
)

// ToXML converts a FrameDSLModel to XML string format
func ToXML(frame model.FrameDSLModel) string {
	return formatter.FormatFrameXML(frame)
}
