package nbx

import (
	"errors"
	"strings"

	"github.com/nativeblocks/nbx/internal/compiler"
	"github.com/nativeblocks/nbx/internal/lexer"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/nativeblocks/nbx/internal/parser"
)

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

// ToDSL converts a FrameJson to a FrameDSLModel.
func ToDSL(frame FrameJson) FrameDSLModel {
	return compiler.ToDsl(frame)
}

// ToJSON converts a FrameDSLModel to a FrameJson, validating with the given schema and frameID.
func ToJSON(frameDSL FrameDSLModel, schema string, frameID string) (FrameJson, error) {
	return compiler.ToJson(frameDSL, schema, frameID)
}

// Parse parses a stringify DSL and returns a FrameDSLModel.
// It returns an error if parsing fails, including parser errors joined by semicolons.
func Parse(stringifyDsl string) (FrameDSLModel, error) {
	l := lexer.NewLexer(stringifyDsl)
	p := parser.NewParser(l)
	frame := p.ParseNBX()
	if frame == nil {
		return FrameDSLModel{}, errors.New(strings.Join(p.Errors(), "; "))
	}
	return *frame, nil
}

// ToString converts a FrameDSLModel back to the original NBX DSL string format.
func ToString(frameDSL FrameDSLModel) string {
	return compiler.ToString(frameDSL)
}
