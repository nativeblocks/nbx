package nbx

import (
	"github.com/nativeblocks/nbx/internal/compiler"
	"github.com/nativeblocks/nbx/internal/errors"
	"github.com/nativeblocks/nbx/internal/formatter"
	"github.com/nativeblocks/nbx/internal/lexer"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/nativeblocks/nbx/internal/parser"
	"github.com/nativeblocks/nbx/internal/validator"
)

type Error = errors.Error
type FrameJson = model.FrameJson
type FrameDSLModel = model.FrameDSLModel

// DSL model types
type VariableDSLModel = model.VariableDSLModel
type BlockDSLModel = model.BlockDSLModel
type BlockPropertyDSLModel = model.BlockPropertyDSLModel
type BlockDataDSLModel = model.BlockDataDSLModel
type BlockSlotDSLModel = model.BlockSlotDSLModel
type ActionDSLModel = model.ActionDSLModel
type ActionTriggerDSLModel = model.ActionTriggerDSLModel
type TriggerPropertyDSLModel = model.TriggerPropertyDSLModel
type TriggerDataDSLModel = model.TriggerDataDSLModel

// JSON model types
type RouteArgumentJson = model.RouteArgumentJson
type VariableJson = model.VariableJson
type BlockJson = model.BlockJson
type BlockPropertyJson = model.BlockPropertyJson
type BlockDataJson = model.BlockDataJson
type BlockSlotJson = model.BlockSlotJson
type ActionJson = model.ActionJson
type ActionTriggerJson = model.ActionTriggerJson
type TriggerPropertyJson = model.TriggerPropertyJson
type TriggerDataJson = model.TriggerDataJson

func Parse(stringifyDsl string) (FrameDSLModel, []Error) {
	l := lexer.NewLexer(stringifyDsl)
	p := parser.NewParser(l, stringifyDsl)
	frame := p.ParseNBX()

	errorCollector := p.ErrorCollector()

	if frame == nil || errorCollector.HasErrors() {
		return FrameDSLModel{}, _errorValueOf(errorCollector.Errors())
	}

	collector, _ := validator.ValidateWithSource(frame, stringifyDsl)

	var all []Error
	if errorCollector != nil {
		all = append(all, _errorValueOf(errorCollector.Errors())...)
		all = append(all, _errorValueOf(errorCollector.Warnings())...)
	}
	if collector != nil {
		all = append(all, _errorValueOf(collector.Errors())...)
		all = append(all, _errorValueOf(collector.Warnings())...)
	}

	return *frame, all
}

// ToDSL converts a FrameJson to a FrameDSLModel.
func ToDSL(frame FrameJson) FrameDSLModel {
	return compiler.ToDsl(frame)
}

// ToJSON converts a FrameDSLModel to a FrameJson with integration validation.
// blocksJSON and actionsJSON must contain the integration definitions in JSON format.
// frameID can be empty to auto-generate an ID.
func ToJSON(frameDSL FrameDSLModel, blocksJSON, actionsJSON, frameID string) (FrameJson, error) {
	return compiler.ToJson(frameDSL, blocksJSON, actionsJSON, frameID)
}

// ToString converts a FrameDSLModel back to DSL string format.
func ToString(frameDSL FrameDSLModel) string {
	return compiler.ToString(frameDSL)
}

// Format takes a DSL string and returns a properly formatted version.
// It parses the DSL to ensure validity and then formats it with consistent indentation and spacing.
func Format(dslString string) (string, error) {
	return formatter.Format(dslString)
}

func _errorValueOf(items []*Error) []Error {
	out := make([]Error, 0, len(items))
	for _, e := range items {
		if e == nil {
			continue
		}
		out = append(out, *e)
	}
	return out
}
