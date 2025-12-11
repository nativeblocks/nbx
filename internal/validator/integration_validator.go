package validator

import (
	"fmt"
	"strings"

	"github.com/nativeblocks/nbx/internal/model"
)

type IntegrationValidator struct {
	registry *IntegrationRegistry
}

func NewIntegrationValidator(registry *IntegrationRegistry) *IntegrationValidator {
	return &IntegrationValidator{
		registry: registry,
	}
}

// ValidateFrame validates all blocks and actions in a frame against the integration registry
func (iv *IntegrationValidator) ValidateFrame(frame *model.FrameDSLModel) error {
	var errors []string

	blockErrors := iv._validateBlocks(frame.Blocks)
	errors = append(errors, blockErrors...)

	if len(errors) > 0 {
		return fmt.Errorf("integration validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

func (iv *IntegrationValidator) _validateBlocks(blocks []model.BlockDSLModel) []string {
	var errors []string

	for _, block := range blocks {
		if block.KeyType == "ROOT" {
			errors = append(errors, iv._validateBlocks(block.Blocks)...)
			continue
		}

		integration, exists := iv.registry.GetBlock(block.KeyType)
		if !exists {
			errors = append(errors, fmt.Sprintf("block '%s' uses unknown integration '%s'", block.Key, block.KeyType))
			errors = append(errors, iv._validateBlocks(block.Blocks)...)
			continue
		}

		propErrors := iv._validateBlockProperties(block, integration)
		errors = append(errors, propErrors...)

		dataErrors := iv._validateBlockData(block, integration)
		errors = append(errors, dataErrors...)

		slotErrors := iv._validateBlockSlots(block, integration)
		errors = append(errors, slotErrors...)

		eventErrors := iv._validateBlockEvents(block, integration)
		errors = append(errors, eventErrors...)

		errors = append(errors, iv._validateBlocks(block.Blocks)...)
	}

	return errors
}

func (iv *IntegrationValidator) _validateBlockProperties(block model.BlockDSLModel, integration BlockIntegration) []string {
	var errors []string

	validProps := make(map[string]PropertyDefinition)
	for _, prop := range integration.Properties {
		validProps[prop.Key] = prop
	}

	for _, prop := range block.Properties {
		if _, exists := validProps[prop.Key]; !exists {
			availableProps := make([]string, 0, len(validProps))
			for key := range validProps {
				availableProps = append(availableProps, key)
			}
			errors = append(errors, fmt.Sprintf(
				"block '%s' uses invalid property '%s' for integration '%s'. Available properties: [%s]",
				block.Key, prop.Key, block.KeyType, strings.Join(availableProps, ", "),
			))
		}
	}

	return errors
}

func (iv *IntegrationValidator) _validateBlockData(block model.BlockDSLModel, integration BlockIntegration) []string {
	var errors []string

	validData := make(map[string]DataDefinition)
	for _, data := range integration.Data {
		validData[data.Key] = data
	}

	for _, data := range block.Data {
		if _, exists := validData[data.Key]; !exists {
			availableData := make([]string, 0, len(validData))
			for key := range validData {
				availableData = append(availableData, key)
			}
			errors = append(errors, fmt.Sprintf(
				"block '%s' uses invalid data key '%s' for integration '%s'. Available data keys: [%s]",
				block.Key, data.Key, block.KeyType, strings.Join(availableData, ", "),
			))
		}
	}

	return errors
}

func (iv *IntegrationValidator) _validateBlockSlots(block model.BlockDSLModel, integration BlockIntegration) []string {
	var errors []string

	validSlots := make(map[string]SlotDefinition)
	for _, slot := range integration.Slots {
		validSlots[slot.Slot] = slot
	}

	for _, slot := range block.Slots {
		if _, exists := validSlots[slot.Slot]; !exists {
			availableSlots := make([]string, 0, len(validSlots))
			for key := range validSlots {
				availableSlots = append(availableSlots, key)
			}
			errors = append(errors, fmt.Sprintf(
				"block '%s' uses invalid slot '%s' for integration '%s'. Available slots: [%s]",
				block.Key, slot.Slot, block.KeyType, strings.Join(availableSlots, ", "),
			))
		}
	}

	return errors
}

func (iv *IntegrationValidator) _validateBlockEvents(block model.BlockDSLModel, integration BlockIntegration) []string {
	var errors []string

	validEvents := make(map[string]EventDefinition)
	for _, event := range integration.Events {
		validEvents[event.Event] = event
	}

	for _, action := range block.Actions {
		if _, exists := validEvents[action.Event]; !exists {
			availableEvents := make([]string, 0, len(validEvents))
			for key := range validEvents {
				availableEvents = append(availableEvents, key)
			}
			errors = append(errors, fmt.Sprintf(
				"block '%s' uses invalid event '%s' for integration '%s'. Available events: [%s]",
				block.Key, action.Event, block.KeyType, strings.Join(availableEvents, ", "),
			))
		}

		triggerErrors := iv._validateTriggers(action.Triggers, block.Key)
		errors = append(errors, triggerErrors...)
	}

	return errors
}

func (iv *IntegrationValidator) _validateTriggers(triggers []model.ActionTriggerDSLModel, blockKey string) []string {
	var errors []string

	for _, trigger := range triggers {
		integration, exists := iv.registry.GetAction(trigger.KeyType)
		if !exists {
			errors = append(errors, fmt.Sprintf(
				"block '%s' uses unknown action integration '%s' in trigger '%s'",
				blockKey, trigger.KeyType, trigger.Name,
			))
			errors = append(errors, iv._validateTriggers(trigger.Triggers, blockKey)...)
			continue
		}

		propErrors := iv._validateTriggerProperties(trigger, integration, blockKey)
		errors = append(errors, propErrors...)

		dataErrors := iv._validateTriggerData(trigger, integration, blockKey)
		errors = append(errors, dataErrors...)

		errors = append(errors, iv._validateTriggers(trigger.Triggers, blockKey)...)
	}

	return errors
}

func (iv *IntegrationValidator) _validateTriggerProperties(trigger model.ActionTriggerDSLModel, integration ActionIntegration, blockKey string) []string {
	var errors []string

	validProps := make(map[string]PropertyDefinition)
	for _, prop := range integration.Properties {
		validProps[prop.Key] = prop
	}

	for _, prop := range trigger.Properties {
		if _, exists := validProps[prop.Key]; !exists {
			availableProps := make([]string, 0, len(validProps))
			for key := range validProps {
				availableProps = append(availableProps, key)
			}
			errors = append(errors, fmt.Sprintf(
				"block '%s' trigger '%s' uses invalid property '%s' for action integration '%s'. Available properties: [%s]",
				blockKey, trigger.Name, prop.Key, trigger.KeyType, strings.Join(availableProps, ", "),
			))
		}
	}

	return errors
}

func (iv *IntegrationValidator) _validateTriggerData(trigger model.ActionTriggerDSLModel, integration ActionIntegration, blockKey string) []string {
	var errors []string

	validData := make(map[string]DataDefinition)
	for _, data := range integration.Data {
		validData[data.Key] = data
	}

	for _, data := range trigger.Data {
		if _, exists := validData[data.Key]; !exists {
			availableData := make([]string, 0, len(validData))
			for key := range validData {
				availableData = append(availableData, key)
			}
			errors = append(errors, fmt.Sprintf(
				"block '%s' trigger '%s' uses invalid data key '%s' for action integration '%s'. Available data keys: [%s]",
				blockKey, trigger.Name, data.Key, trigger.KeyType, strings.Join(availableData, ", "),
			))
		}
	}

	return errors
}
