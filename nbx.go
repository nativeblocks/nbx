package nbx

import (
	"github.com/nativeblocks/nbx/internal/compiler"
	"github.com/nativeblocks/nbx/internal/model"
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
