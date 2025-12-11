package validator

import (
	"encoding/json"
	"fmt"
)

type BlockIntegration struct {
	KeyType    string               `json:"keyType"`
	Version    int                  `json:"version"`
	Properties []PropertyDefinition `json:"properties"`
	Data       []DataDefinition     `json:"data"`
	Events     []EventDefinition    `json:"events"`
	Slots      []SlotDefinition     `json:"slots"`
}

type ActionIntegration struct {
	KeyType    string               `json:"keyType"`
	Version    int                  `json:"version"`
	Properties []PropertyDefinition `json:"properties"`
	Data       []DataDefinition     `json:"data"`
	Events     []EventDefinition    `json:"events"`
}

type PropertyDefinition struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type DataDefinition struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type EventDefinition struct {
	Event string `json:"event"`
}

type SlotDefinition struct {
	Slot string `json:"slot"`
}

type IntegrationRegistry struct {
	Blocks  map[string]BlockIntegration
	Actions map[string]ActionIntegration
}

// LoadIntegrations creates an IntegrationRegistry from JSON strings.
// blocksJSON and actionsJSON should contain the integration definitions in the expected format.
func LoadIntegrations(blocksJSON, actionsJSON string) (*IntegrationRegistry, error) {
	registry := &IntegrationRegistry{
		Blocks:  make(map[string]BlockIntegration),
		Actions: make(map[string]ActionIntegration),
	}

	if err := registry._parseBlocks(blocksJSON); err != nil {
		return nil, fmt.Errorf("failed to parse blocks JSON: %w", err)
	}

	if err := registry._parseActions(actionsJSON); err != nil {
		return nil, fmt.Errorf("failed to parse actions JSON: %w", err)
	}

	return registry, nil
}

func (r *IntegrationRegistry) _parseBlocks(blocksJSON string) error {
	var blocksMap map[string]interface{}
	if err := json.Unmarshal([]byte(blocksJSON), &blocksMap); err != nil {
		return fmt.Errorf("invalid blocks JSON: %w", err)
	}

	for keyType, value := range blocksMap {
		if keyType == "schema-version" {
			continue
		}

		valueBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal block %s: %w", keyType, err)
		}

		var block BlockIntegration
		if err := json.Unmarshal(valueBytes, &block); err != nil {
			return fmt.Errorf("failed to unmarshal block %s: %w", keyType, err)
		}

		block.KeyType = keyType
		r.Blocks[keyType] = block
	}

	return nil
}

func (r *IntegrationRegistry) _parseActions(actionsJSON string) error {
	var actionsMap map[string]interface{}
	if err := json.Unmarshal([]byte(actionsJSON), &actionsMap); err != nil {
		return fmt.Errorf("invalid actions JSON: %w", err)
	}

	for keyType, value := range actionsMap {
		if keyType == "schema-version" {
			continue
		}

		valueBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal action %s: %w", keyType, err)
		}

		var action ActionIntegration
		if err := json.Unmarshal(valueBytes, &action); err != nil {
			return fmt.Errorf("failed to unmarshal action %s: %w", keyType, err)
		}

		action.KeyType = keyType
		r.Actions[keyType] = action
	}

	return nil
}

func (r *IntegrationRegistry) GetBlock(keyType string) (BlockIntegration, bool) {
	block, exists := r.Blocks[keyType]
	return block, exists
}

func (r *IntegrationRegistry) GetAction(keyType string) (ActionIntegration, bool) {
	action, exists := r.Actions[keyType]
	return action, exists
}
