package parser

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/nativeblocks/nbx/internal/errors"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/nativeblocks/nbx/internal/types"
)

// ParseXML parses XML format NBX content into a FrameDSLModel
func ParseXML(xmlString string) (model.FrameDSLModel, []*errors.Error) {
	if xmlString == "" {
		return model.FrameDSLModel{}, []*errors.Error{{
			Severity: errors.SeverityError,
			Message:  "XML content is empty",
			Line:     0,
			Column:   0,
		}}
	}

	posTracker := NewPositionTracker(xmlString)
	errorCollector := errors.NewErrorCollector(xmlString)

	var xmlFrame model.XMLFrame
	decoder := xml.NewDecoder(strings.NewReader(xmlString))

	if err := decoder.Decode(&xmlFrame); err != nil {
		errorCollector.AddError(&errors.Error{
			Severity: errors.SeverityError,
			Message:  fmt.Sprintf("Failed to parse XML: %v", err),
			Line:     0,
			Column:   0,
		})
		return model.FrameDSLModel{}, errorCollector.AllIssues()
	}

	// Validate basic structure
	if xmlFrame.Name == "" {
		errorCollector.AddError(&errors.Error{
			Severity: errors.SeverityError,
			Message:  "Frame name is required",
			Line:     posTracker.FindElementPosition("frame", "").Line,
			Column:   posTracker.FindElementPosition("frame", "").Column,
		})
	}

	if xmlFrame.Route == "" {
		errorCollector.AddError(&errors.Error{
			Severity: errors.SeverityError,
			Message:  "Frame route is required",
			Line:     posTracker.FindElementPosition("frame", "").Line,
			Column:   posTracker.FindElementPosition("frame", "").Column,
		})
	}

	if errorCollector.HasErrors() {
		return model.FrameDSLModel{}, errorCollector.AllIssues()
	}

	frame := _toFrameDSLModel(xmlFrame, posTracker)

	return frame, errorCollector.AllIssues()
}

func _toFrameDSLModel(xf model.XMLFrame, tracker *PositionTracker) model.FrameDSLModel {
	pos := tracker.FindElementPosition("frame", xf.Name)

	frame := model.FrameDSLModel{
		Name:      xf.Name,
		Route:     xf.Route,
		Type:      xf.Type,
		Variables: make([]model.VariableDSLModel, 0, len(xf.Variables)),
		Blocks:    make([]model.BlockDSLModel, 0, len(xf.Blocks)),
		Line:      pos.Line,
		Column:    pos.Column,
	}

	if frame.Type == "" {
		frame.Type = "FRAME"
	}

	for _, xv := range xf.Variables {
		varPos := tracker.FindElementPosition("var", xv.Key)
		frame.Variables = append(frame.Variables, model.VariableDSLModel{
			Key:    xv.Key,
			Type:   strings.ToUpper(xv.Type),
			Value:  xv.Value,
			Line:   varPos.Line,
			Column: varPos.Column,
		})
	}

	for _, xb := range xf.Blocks {
		frame.Blocks = append(frame.Blocks, _toBlockDSLModel(xb, tracker))
	}

	return frame
}

func _toBlockDSLModel(xb model.XMLBlock, tracker *PositionTracker) model.BlockDSLModel {
	pos := tracker.FindElementPosition("block", xb.Key)

	block := model.BlockDSLModel{
		KeyType:            xb.KeyType,
		Key:                xb.Key,
		VisibilityKey:      xb.Visibility,
		IntegrationVersion: xb.Version,
		Properties:         make([]model.BlockPropertyDSLModel, 0),
		Data:               make([]model.BlockDataDSLModel, 0),
		Slots:              make([]model.BlockSlotDSLModel, 0),
		Blocks:             make([]model.BlockDSLModel, 0),
		Actions:            make([]model.ActionDSLModel, 0),
		Line:               pos.Line,
		Column:             pos.Column,
	}

	for _, xp := range xb.Properties {
		propPos := tracker.FindElementPosition("prop", xp.Key)

		mobile := xp.Mobile
		tablet := xp.Tablet
		desktop := xp.Desktop

		if xp.Value != "" {
			mobile = xp.Value
			tablet = xp.Value
			desktop = xp.Value
		}

		inferValue := mobile
		if inferValue == "" {
			inferValue = tablet
		}
		if inferValue == "" {
			inferValue = desktop
		}

		block.Properties = append(block.Properties, model.BlockPropertyDSLModel{
			Key:          xp.Key,
			ValueMobile:  mobile,
			ValueTablet:  tablet,
			ValueDesktop: desktop,
			Type:         types.InferType(inferValue).Name(),
			Line:         propPos.Line,
			Column:       propPos.Column,
		})
	}

	for _, xd := range xb.Data {
		dataPos := tracker.FindElementPosition("data", xd.Key)
		block.Data = append(block.Data, model.BlockDataDSLModel{
			Key:    xd.Key,
			Value:  xd.Value,
			Type:   types.InferType(xd.Value).Name(),
			Line:   dataPos.Line,
			Column: dataPos.Column,
		})
	}

	for _, xs := range xb.Slots {
		slotPos := tracker.FindElementPosition("slot", xs.Name)
		block.Slots = append(block.Slots, model.BlockSlotDSLModel{
			Slot:   xs.Name,
			Line:   slotPos.Line,
			Column: slotPos.Column,
		})

		for _, childBlock := range xs.Blocks {
			childModel := _toBlockDSLModel(childBlock, tracker)
			childModel.Slot = xs.Name
			block.Blocks = append(block.Blocks, childModel)
		}
	}

	for _, xa := range xb.Actions {
		actionPos := tracker.FindElementPosition("action", xa.Event)
		action := model.ActionDSLModel{
			Key:      xb.Key,
			Event:    xa.Event,
			Triggers: make([]model.ActionTriggerDSLModel, 0),
			Line:     actionPos.Line,
			Column:   actionPos.Column,
		}

		for _, xt := range xa.Triggers {
			action.Triggers = append(action.Triggers, _toTriggerDSLModel(xt, tracker, "NEXT"))
		}

		block.Actions = append(block.Actions, action)
	}

	return block
}

func _toTriggerDSLModel(xt model.XMLTrigger, tracker *PositionTracker, defaultThen string) model.ActionTriggerDSLModel {
	pos := tracker.FindElementPosition("trigger", xt.Name)

	trigger := model.ActionTriggerDSLModel{
		KeyType:            xt.KeyType,
		Name:               xt.Name,
		Then:               defaultThen,
		IntegrationVersion: xt.Version,
		Properties:         make([]model.TriggerPropertyDSLModel, 0),
		Data:               make([]model.TriggerDataDSLModel, 0),
		Triggers:           make([]model.ActionTriggerDSLModel, 0),
		Line:               pos.Line,
		Column:             pos.Column,
	}

	for _, xp := range xt.Properties {
		propPos := tracker.FindElementPosition("prop", xp.Key)
		value := xp.Value

		trigger.Properties = append(trigger.Properties, model.TriggerPropertyDSLModel{
			Key:    xp.Key,
			Value:  value,
			Type:   types.InferType(value).Name(),
			Line:   propPos.Line,
			Column: propPos.Column,
		})
	}

	for _, xd := range xt.Data {
		dataPos := tracker.FindElementPosition("data", xd.Key)
		trigger.Data = append(trigger.Data, model.TriggerDataDSLModel{
			Key:    xd.Key,
			Value:  xd.Value,
			Type:   types.InferType(xd.Value).Name(),
			Line:   dataPos.Line,
			Column: dataPos.Column,
		})
	}

	for _, th := range xt.Then {
		for _, nestedTrigger := range th.Triggers {
			nested := _toTriggerDSLModel(nestedTrigger, tracker, th.Value)
			trigger.Triggers = append(trigger.Triggers, nested)
		}
	}

	if len(trigger.Triggers) == 0 && defaultThen == "" {
		trigger.Then = "NEXT"
	}

	return trigger
}
