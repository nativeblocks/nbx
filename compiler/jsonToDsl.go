package compiler

import "github.com/nativeblocks/nbx/model"

func ToDsl(frame model.FrameJson) model.FrameDSLModel {
	variables := make([]model.VariableDSLModel, len(frame.Variables))
	for i, variable := range frame.Variables {
		variables[i] = _mapVariableModelToDSL(variable)
	}
	return model.FrameDSLModel{
		Name:      frame.Name,
		Route:     frame.Route,
		Type:      frame.Type,
		Variables: variables,
		Blocks:    _buildBlockTreeWithActions(frame.Blocks, frame.Actions),
	}
}

func _findActionTriggerChildren(triggers []model.ActionTriggerJson, parentId string) []model.ActionTriggerDSLModel {
	var children []model.ActionTriggerDSLModel

	for _, trigger := range triggers {
		if trigger.ParentId == parentId {
			child := model.ActionTriggerDSLModel{
				KeyType:            trigger.KeyType,
				Then:               trigger.Then,
				Name:               trigger.Name,
				IntegrationVersion: trigger.IntegrationVersion,
				Properties:         make([]model.TriggerPropertyDSLModel, len(trigger.Properties)),
				Data:               make([]model.TriggerDataDSLModel, len(trigger.Data)),
			}

			for i, prop := range trigger.Properties {
				child.Properties[i] = _mapTriggerPropertyModelToDSL(prop)
			}

			for i, data := range trigger.Data {
				child.Data[i] = _mapTriggerDataModelToDSL(data)
			}

			child.Triggers = _findActionTriggerChildren(triggers, trigger.Id)
			children = append(children, child)
		}
	}

	if children == nil {
		return make([]model.ActionTriggerDSLModel, 0)
	} else {
		return children
	}
}

func _buildActionTriggerTree(triggers []model.ActionTriggerJson) []model.ActionTriggerDSLModel {
	var roots []model.ActionTriggerDSLModel

	for _, trigger := range triggers {
		if trigger.ParentId == "" {
			root := model.ActionTriggerDSLModel{
				KeyType:            trigger.KeyType,
				Then:               trigger.Then,
				Name:               trigger.Name,
				IntegrationVersion: trigger.IntegrationVersion,
				Properties:         make([]model.TriggerPropertyDSLModel, len(trigger.Properties)),
				Data:               make([]model.TriggerDataDSLModel, len(trigger.Data)),
			}

			for i, prop := range trigger.Properties {
				root.Properties[i] = _mapTriggerPropertyModelToDSL(prop)
			}

			for i, data := range trigger.Data {
				root.Data[i] = _mapTriggerDataModelToDSL(data)
			}

			root.Triggers = _findActionTriggerChildren(triggers, trigger.Id)
			roots = append(roots, root)
		}
	}

	if roots == nil {
		return make([]model.ActionTriggerDSLModel, 0)
	} else {
		return roots
	}
}

func _findBlockChildren(blocks []model.BlockJson, parentId string, actions []model.ActionJson) []model.BlockDSLModel {
	var children []model.BlockDSLModel

	for _, block := range blocks {
		if block.ParentId == parentId {
			child := model.BlockDSLModel{
				KeyType:            block.KeyType,
				Key:                block.Key,
				VisibilityKey:      block.VisibilityKey,
				Slot:               block.Slot,
				IntegrationVersion: block.IntegrationVersion,
				Data:               make([]model.BlockDataDSLModel, len(block.Data)),
				Properties:         make([]model.BlockPropertyDSLModel, len(block.Properties)),
				Slots:              make([]model.BlockSlotDSLModel, len(block.Slots)),
			}

			for i, data := range block.Data {
				child.Data[i] = _mapBlockDataModelToDSL(data)
			}

			for i, prop := range block.Properties {
				child.Properties[i] = _mapBlockPropertyModelToDSL(prop)
			}

			for i, slot := range block.Slots {
				child.Slots[i] = _mapBlockSlotModelToDSL(slot)
			}

			for _, action := range actions {
				if action.Key == block.Key {
					child.Actions = append(child.Actions, _mapActionModelToDSL(action))
				}
			}

			if child.Actions == nil {
				child.Actions = make([]model.ActionDSLModel, 0)
			}

			child.Blocks = _findBlockChildren(blocks, block.Id, actions)
			children = append(children, child)
		}
	}

	if children == nil {
		return make([]model.BlockDSLModel, 0)
	} else {
		return children
	}
}

func _buildBlockTreeWithActions(blocks []model.BlockJson, actions []model.ActionJson) []model.BlockDSLModel {
	var dslBlocks []model.BlockDSLModel
	for _, block := range blocks {
		if block.ParentId == "" {
			root := model.BlockDSLModel{
				KeyType:            block.KeyType,
				Key:                block.Key,
				VisibilityKey:      block.VisibilityKey,
				Slot:               "null",
				IntegrationVersion: block.IntegrationVersion,
				Data:               make([]model.BlockDataDSLModel, len(block.Data)),
				Properties:         make([]model.BlockPropertyDSLModel, len(block.Properties)),
				Slots:              make([]model.BlockSlotDSLModel, len(block.Slots)),
			}

			for i, data := range block.Data {
				root.Data[i] = _mapBlockDataModelToDSL(data)
			}

			for i, prop := range block.Properties {
				root.Properties[i] = _mapBlockPropertyModelToDSL(prop)
			}

			for i, slot := range block.Slots {
				root.Slots[i] = _mapBlockSlotModelToDSL(slot)
			}

			for _, action := range actions {
				if action.Key == block.Key {
					root.Actions = append(root.Actions, _mapActionModelToDSL(action))
				}
			}

			if root.Actions == nil {
				root.Actions = make([]model.ActionDSLModel, 0)
			}

			root.Blocks = _findBlockChildren(blocks, block.Id, actions)
			dslBlocks = append(dslBlocks, root)
		}
	}

	if dslBlocks == nil {
		return make([]model.BlockDSLModel, 0)
	} else {
		return dslBlocks
	}
}

func _mapVariableModelToDSL(variable model.VariableJson) model.VariableDSLModel {
	return model.VariableDSLModel{
		Key:   variable.Key,
		Value: variable.Value,
		Type:  variable.Type,
	}
}

func _mapBlockDataModelToDSL(data model.BlockDataJson) model.BlockDataDSLModel {
	return model.BlockDataDSLModel{
		Key:   data.Key,
		Value: data.Value,
		Type:  data.Type,
	}
}

func _mapBlockPropertyModelToDSL(property model.BlockPropertyJson) model.BlockPropertyDSLModel {
	return model.BlockPropertyDSLModel{
		Key:          property.Key,
		ValueMobile:  property.ValueMobile,
		ValueTablet:  property.ValueTablet,
		ValueDesktop: property.ValueDesktop,
		Type:         property.Type,
	}
}

func _mapBlockSlotModelToDSL(slot model.BlockSlotJson) model.BlockSlotDSLModel {
	return model.BlockSlotDSLModel{
		Slot: slot.Slot,
	}
}

func _mapActionModelToDSL(action model.ActionJson) model.ActionDSLModel {
	return model.ActionDSLModel{
		Key:      action.Key,
		Event:    action.Event,
		Triggers: _buildActionTriggerTree(action.Triggers),
	}
}

func _mapTriggerPropertyModelToDSL(property model.TriggerPropertyJson) model.TriggerPropertyDSLModel {
	return model.TriggerPropertyDSLModel{
		Key:   property.Key,
		Value: property.Value,
		Type:  property.Type,
	}
}

func _mapTriggerDataModelToDSL(data model.TriggerDataJson) model.TriggerDataDSLModel {
	return model.TriggerDataDSLModel{
		Key:   data.Key,
		Value: data.Value,
		Type:  data.Type,
	}
}
