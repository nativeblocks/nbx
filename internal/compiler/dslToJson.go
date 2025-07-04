package compiler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/xeipuuv/gojsonschema"
	"regexp"
	"strings"
)

func ToJson(frameDSL model.FrameDSLModel, schema string, frameID string) (model.FrameJson, error) {
	initialFrameId := frameID
	if initialFrameId == "" {
		initialFrameId = _generateId()
	}

	if len(frameDSL.Blocks) > 0 && frameDSL.Blocks[0].KeyType != "ROOT" {
		return model.FrameJson{}, errors.New("first block's keyType must be 'ROOT'")
	}

	result, err := _validateDSL(frameDSL, schema)
	if err != nil {
		return model.FrameJson{}, err
	}

	variables := _createVariables(frameDSL.Variables, initialFrameId)
	var actions []model.ActionJson
	blocks, blockErr := _processBlocks(initialFrameId, frameDSL.Blocks, "", []model.BlockSlotJson{}, variables, func(blockActions []model.ActionJson) {
		actions = append(actions, blockActions...)
	})

	var frameJson model.FrameJson
	if result != nil && !result.Valid() {
		frameJson = _createFrame(initialFrameId, frameDSL, variables, blocks, actions)
		validationErrors := _formatValidationErrors(result, frameJson)
		return model.FrameJson{}, fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	if blockErr != nil {
		return model.FrameJson{}, blockErr
	}

	if duplicateKeys := _findDuplicateKeys(blocks); len(duplicateKeys) > 0 {
		return model.FrameJson{}, errors.New("duplicate block keys found: " + strings.Join(duplicateKeys, ","))
	}

	if frameJson.Actions == nil {
		frameJson.Actions = []model.ActionJson{}
	}
	if frameJson.Blocks == nil {
		frameJson.Blocks = []model.BlockJson{}
	}

	return frameJson, nil
}

func _validateDSL(frameDSL model.FrameDSLModel, schema string) (*gojsonschema.Result, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	if len(frameDSL.Blocks) > 0 && frameDSL.Blocks[0].KeyType == "ROOT" && frameDSL.Blocks[0].Slot == "" {
		// because the root has no parent, we need to pass a null slot to it
		validationDSL := frameDSL
		validationDSL.Blocks = make([]model.BlockDSLModel, len(frameDSL.Blocks))
		copy(validationDSL.Blocks, frameDSL.Blocks)
		validationDSL.Blocks[0].Slot = "null"
		return gojsonschema.Validate(schemaLoader, gojsonschema.NewGoLoader(validationDSL))
	}
	return gojsonschema.Validate(schemaLoader, gojsonschema.NewGoLoader(frameDSL))
}

func _createVariables(variableDSLs []model.VariableDSLModel, frameId string) []model.VariableJson {
	variables := make([]model.VariableJson, 0, len(variableDSLs))
	for _, variable := range variableDSLs {
		variables = append(variables, model.VariableJson{
			Id:      _generateId(),
			FrameId: frameId,
			Key:     variable.Key,
			Value:   variable.Value,
			Type:    variable.Type,
		})
	}
	return variables
}

func _createFrame(frameId string, dsl model.FrameDSLModel, variables []model.VariableJson, blocks []model.BlockJson, actions []model.ActionJson) model.FrameJson {
	return model.FrameJson{
		Id:             frameId,
		Name:           dsl.Name,
		Route:          dsl.Route,
		RouteArguments: _convertRouteArguments(dsl.Route),
		Type:           dsl.Type,
		Variables:      variables,
		Blocks:         blocks,
		Actions:        actions,
	}
}

func _processActions(frameId, key string, inputActions []model.ActionDSLModel, variables []model.VariableJson) ([]model.ActionJson, error) {
	var actions []model.ActionJson

	for _, inputAction := range inputActions {
		actionId := _generateId()
		subTriggers, err := _processTriggers(actionId, inputAction.Triggers, "", variables)
		if err != nil {
			return nil, err
		}

		newAction := model.ActionJson{
			Id:       actionId,
			FrameId:  frameId,
			Key:      key,
			Event:    inputAction.Event,
			Triggers: subTriggers,
		}
		actions = append(actions, newAction)
	}

	return actions, nil
}

func _processTriggers(actionId string, triggers []model.ActionTriggerDSLModel, parentId string, variables []model.VariableJson) ([]model.ActionTriggerJson, error) {
	var flatTriggers []model.ActionTriggerJson

	for _, trigger := range triggers {
		newTrigger := model.ActionTriggerJson{
			Id:                 _generateId(),
			ActionId:           actionId,
			ParentId:           parentId,
			KeyType:            trigger.KeyType,
			Then:               trigger.Then,
			Name:               trigger.Name,
			IntegrationVersion: trigger.IntegrationVersion,
			Properties:         []model.TriggerPropertyJson{},
			Data:               []model.TriggerDataJson{},
		}

		if newTrigger.Then == "END" && len(trigger.Triggers) > 0 {
			return nil, errors.New("The " + newTrigger.Name + " can not have a subTrigger because it defines with \"END\" then ")
		}

		for _, property := range trigger.Properties {
			newProperty := model.TriggerPropertyJson{
				ActionTriggerId:    newTrigger.Id,
				Key:                property.Key,
				Type:               property.Type,
				Value:              property.Value,
				Description:        "",
				ValuePicker:        "",
				ValuePickerGroup:   "",
				ValuePickerOptions: "",
			}
			newTrigger.Properties = append(newTrigger.Properties, newProperty)
		}

		for _, dataItem := range trigger.Data {
			newData := model.TriggerDataJson{
				ActionTriggerId: newTrigger.Id,
				Key:             dataItem.Key,
				Value:           dataItem.Value,
				Type:            dataItem.Type,
				Description:     "",
			}
			newTrigger.Data = append(newTrigger.Data, newData)
		}

		err := _findTriggerVariable(variables, newTrigger.Data, newTrigger.Name)
		if err != nil {
			return nil, err
		}

		flatTriggers = append(flatTriggers, newTrigger)

		if len(trigger.Triggers) > 0 {
			subTriggers, err := _processTriggers(actionId, trigger.Triggers, newTrigger.Id, variables)
			if err != nil {
				return nil, err
			}
			flatTriggers = append(flatTriggers, subTriggers...)
		}
	}

	if flatTriggers == nil {
		flatTriggers = []model.ActionTriggerJson{}
	}

	return flatTriggers, nil
}

func _processBlocks(frameId string, blocks []model.BlockDSLModel, parentId string, parentSlots []model.BlockSlotJson, variables []model.VariableJson, onNewAction func([]model.ActionJson)) ([]model.BlockJson, error) {
	var flatBlocks []model.BlockJson

	for index, block := range blocks {
		newBlock := model.BlockJson{
			Id:                 _generateId(),
			FrameId:            frameId,
			KeyType:            block.KeyType,
			Key:                block.Key,
			VisibilityKey:      block.VisibilityKey,
			Position:           index,
			Slot:               block.Slot,
			IntegrationVersion: block.IntegrationVersion,
			ParentId:           parentId,
			Data:               []model.BlockDataJson{},
			Properties:         []model.BlockPropertyJson{},
			Slots:              []model.BlockSlotJson{},
		}

		if newBlock.Slot == "null" {
			emptySlot := ""
			newBlock.Slot = emptySlot
		}

		if len(parentSlots) > 0 {
			contain := _containsSlot(parentSlots, newBlock.Slot)
			if !contain {
				return nil, errors.New("The " + newBlock.Key + " used in a wrong slot")
			}
		}

		processedActions, err := _processActions(frameId, block.Key, block.Actions, variables)
		if err != nil {
			return nil, err
		}
		onNewAction(processedActions)

		for _, property := range block.Properties {
			newProperty := model.BlockPropertyJson{
				BlockId:            newBlock.Id,
				Key:                property.Key,
				Type:               property.Type,
				ValueMobile:        property.ValueMobile,
				ValueTablet:        property.ValueTablet,
				ValueDesktop:       property.ValueDesktop,
				Description:        "",
				ValuePicker:        "",
				ValuePickerGroup:   "",
				ValuePickerOptions: "",
			}
			newBlock.Properties = append(newBlock.Properties, newProperty)
		}

		for _, dataItem := range block.Data {
			newData := model.BlockDataJson{
				BlockId:     newBlock.Id,
				Key:         dataItem.Key,
				Value:       dataItem.Value,
				Type:        dataItem.Type,
				Description: "",
			}
			newBlock.Data = append(newBlock.Data, newData)
		}

		for _, slotItem := range block.Slots {
			newSlot := model.BlockSlotJson{
				BlockId:     newBlock.Id,
				Slot:        slotItem.Slot,
				Description: "",
			}
			newBlock.Slots = append(newBlock.Slots, newSlot)
		}

		err = _findBlockVariable(variables, newBlock.Data, newBlock.Key)
		if err != nil {
			return nil, err
		}

		flatBlocks = append(flatBlocks, newBlock)

		if len(block.Blocks) > 0 {
			subBlocks, err := _processBlocks(frameId, block.Blocks, newBlock.Id, newBlock.Slots, variables, onNewAction)
			if err != nil {
				return nil, err
			}
			flatBlocks = append(flatBlocks, subBlocks...)
		}
	}

	return flatBlocks, nil
}

func _getWordsBetweenCurly(text string) []string {
	re := regexp.MustCompile(`\{(.*?)}`)
	matches := re.FindAllStringSubmatch(text, -1)

	var result []string
	for _, match := range matches {
		if len(match) > 1 {
			result = append(result, match[1])
		}
	}
	return result
}

func _convertRouteArguments(route string) []model.RouteArgumentJson {
	args := _getWordsBetweenCurly(route)
	routeArguments := make([]model.RouteArgumentJson, len(args))

	for i, arg := range args {
		routeArguments[i] = model.RouteArgumentJson{Name: arg}
	}
	return routeArguments
}

func _findBlockVariable(variables []model.VariableJson, data []model.BlockDataJson, blockKey string) error {
	for _, dataEntry := range data {
		found := false
		for _, variable := range variables {
			if variable.Key == dataEntry.Value {
				found = true
				break
			}
		}
		if !found && dataEntry.Value != "" {
			return fmt.Errorf("no matching variable found for %s block in data entry with key: %s", blockKey, dataEntry.Key)
		}
	}
	return nil
}

func _findTriggerVariable(variables []model.VariableJson, data []model.TriggerDataJson, triggerName string) error {
	for _, dataEntry := range data {
		found := false
		for _, variable := range variables {
			if variable.Key == dataEntry.Value {
				found = true
				break
			}
		}
		if !found && dataEntry.Value != "" {
			return fmt.Errorf("no matching variable found for %s trigger in data entry with key: %s", triggerName, dataEntry.Key)
		}
	}
	return nil
}

func _containsSlot(slots []model.BlockSlotJson, key string) bool {
	for _, slot := range slots {
		if slot.Slot == key {
			return true
		}
	}
	return false
}

func _findDuplicateKeys(blocks []model.BlockJson) []string {
	keyCount := make(map[string]int)
	var duplicates []string

	for _, block := range blocks {
		keyCount[block.Key]++
	}

	for key, count := range keyCount {
		if count > 1 {
			duplicates = append(duplicates, key)
		}
	}
	return duplicates
}

func _generateId() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}
