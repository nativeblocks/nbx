package nbx

import (
	"errors"
	"github.com/nativeblocks/nbx/internal/compiler"
	"github.com/nativeblocks/nbx/internal/lexer"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/nativeblocks/nbx/internal/parser"
	"strings"
)

type FrameJson = model.FrameJson

type FrameDSLModel = model.FrameDSLModel

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
	frame := p.ParseSDUI()
	if frame == nil {
		return FrameDSLModel{}, errors.New(strings.Join(p.Errors(), "; "))
	}
	return *frame, nil
}
