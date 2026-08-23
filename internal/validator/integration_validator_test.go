package validator

import (
	"strings"
	"testing"

	"github.com/nativeblocks/nbx/internal/model"
)

const scopedBlocksJSON = `{
  "schema-version": "test",
  "nativeblocks/lazy_column": {
    "keyType": "nativeblocks/lazy_column",
    "version": 1,
    "properties": [],
    "data": [],
    "events": [{"event": "onClick", "scope": "lazyListScope"}],
    "slots": [{"slot": "content", "scope": "lazyListScope"}]
  },
  "nativeblocks/column": {
    "keyType": "nativeblocks/column",
    "version": 1,
    "properties": [],
    "data": [],
    "events": [{"event": "onClick"}],
    "slots": [{"slot": "content"}]
  },
  "nativeblocks/lazy_item": {
    "keyType": "nativeblocks/lazy_item",
    "version": 1,
    "scope": "lazyListScope",
    "properties": [],
    "data": [],
    "events": [],
    "slots": [{"slot": "content"}]
  }
}`

const scopedActionsJSON = `{
  "schema-version": "test",
  "SCRIPT": {
    "keyType": "SCRIPT",
    "version": 1,
    "properties": [],
    "data": [],
    "events": [{"event": "SUCCESS", "scope": "resultScope"}, {"event": "FAILURE"}]
  },
  "READ_RESULT": {
    "keyType": "READ_RESULT",
    "version": 1,
    "scope": "resultScope",
    "properties": [],
    "data": [],
    "events": []
  }
}`

func newScopedValidator(t *testing.T) *IntegrationValidator {
	t.Helper()

	registry, err := LoadIntegrations(scopedBlocksJSON, scopedActionsJSON)
	if err != nil {
		t.Fatalf("Failed to load integrations: %v", err)
	}

	return NewIntegrationValidator(registry)
}

func newFrame(blocks ...model.BlockDSLModel) *model.FrameDSLModel {
	return &model.FrameDSLModel{
		Name:  "test",
		Route: "/test",
		Type:  "FRAME",
		Blocks: []model.BlockDSLModel{
			{
				KeyType: "ROOT",
				Key:     "root",
				Slots:   []model.BlockSlotDSLModel{{Slot: "content"}},
				Blocks:  blocks,
			},
		},
	}
}

func TestLoadIntegrationsParsesScope(t *testing.T) {
	registry, err := LoadIntegrations(scopedBlocksJSON, scopedActionsJSON)
	if err != nil {
		t.Fatalf("Failed to load integrations: %v", err)
	}

	block, _ := registry.GetBlock("nativeblocks/lazy_item")
	if block.Scope != "lazyListScope" {
		t.Errorf("Expected block scope 'lazyListScope', got '%s'", block.Scope)
	}

	container, _ := registry.GetBlock("nativeblocks/lazy_column")
	if container.Slots[0].Scope != "lazyListScope" {
		t.Errorf("Expected slot scope 'lazyListScope', got '%s'", container.Slots[0].Scope)
	}
	if container.Events[0].Scope != "lazyListScope" {
		t.Errorf("Expected event scope 'lazyListScope', got '%s'", container.Events[0].Scope)
	}

	action, _ := registry.GetAction("READ_RESULT")
	if action.Scope != "resultScope" {
		t.Errorf("Expected action scope 'resultScope', got '%s'", action.Scope)
	}
}

func TestBlockScopeSatisfiedBySlot(t *testing.T) {
	frame := newFrame(model.BlockDSLModel{
		KeyType: "nativeblocks/lazy_column",
		Key:     "list",
		Slot:    "content",
		Slots:   []model.BlockSlotDSLModel{{Slot: "content"}},
		Blocks: []model.BlockDSLModel{
			{KeyType: "nativeblocks/lazy_item", Key: "item", Slot: "content"},
		},
	})

	if err := newScopedValidator(t).ValidateFrame(frame); err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestBlockScopeMissingOutsideSlot(t *testing.T) {
	frame := newFrame(model.BlockDSLModel{
		KeyType: "nativeblocks/lazy_item",
		Key:     "item",
		Slot:    "content",
	})

	err := newScopedValidator(t).ValidateFrame(frame)
	if err == nil {
		t.Fatal("Expected error for block used outside its scope")
	}
	if !strings.Contains(err.Error(), "block 'item' requires scope 'lazyListScope'") {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestBlockScopeNotProvidedByUnscopedSlot(t *testing.T) {
	frame := newFrame(model.BlockDSLModel{
		KeyType: "nativeblocks/column",
		Key:     "column",
		Slot:    "content",
		Slots:   []model.BlockSlotDSLModel{{Slot: "content"}},
		Blocks: []model.BlockDSLModel{
			{KeyType: "nativeblocks/lazy_item", Key: "item", Slot: "content"},
		},
	})

	err := newScopedValidator(t).ValidateFrame(frame)
	if err == nil {
		t.Fatal("Expected error for block placed in a slot that provides no scope")
	}
	if !strings.Contains(err.Error(), "not placed in a slot that provides a scope") {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestBlockScopeDoesNotLeakToGrandChildren(t *testing.T) {
	frame := newFrame(model.BlockDSLModel{
		KeyType: "nativeblocks/lazy_column",
		Key:     "list",
		Slot:    "content",
		Slots:   []model.BlockSlotDSLModel{{Slot: "content"}},
		Blocks: []model.BlockDSLModel{
			{
				KeyType: "nativeblocks/lazy_item",
				Key:     "item",
				Slot:    "content",
				Slots:   []model.BlockSlotDSLModel{{Slot: "content"}},
				Blocks: []model.BlockDSLModel{
					{KeyType: "nativeblocks/lazy_item", Key: "nestedItem", Slot: "content"},
				},
			},
		},
	})

	err := newScopedValidator(t).ValidateFrame(frame)
	if err == nil {
		t.Fatal("Expected error for scoped block nested below the scoped slot")
	}
	if !strings.Contains(err.Error(), "block 'nestedItem' requires scope 'lazyListScope'") {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestTriggerScopeSatisfiedByEvent(t *testing.T) {
	frame := newFrame(model.BlockDSLModel{
		KeyType: "nativeblocks/column",
		Key:     "column",
		Slot:    "content",
		Actions: []model.ActionDSLModel{
			{
				Key:   "column",
				Event: "onClick",
				Triggers: []model.ActionTriggerDSLModel{
					{
						KeyType: "SCRIPT",
						Name:    "run",
						Then:    "NEXT",
						Triggers: []model.ActionTriggerDSLModel{
							{KeyType: "READ_RESULT", Name: "read", Then: "SUCCESS"},
						},
					},
				},
			},
		},
	})

	if err := newScopedValidator(t).ValidateFrame(frame); err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestTriggerScopeMissingOnUnscopedEvent(t *testing.T) {
	frame := newFrame(model.BlockDSLModel{
		KeyType: "nativeblocks/column",
		Key:     "column",
		Slot:    "content",
		Actions: []model.ActionDSLModel{
			{
				Key:   "column",
				Event: "onClick",
				Triggers: []model.ActionTriggerDSLModel{
					{
						KeyType: "SCRIPT",
						Name:    "run",
						Then:    "NEXT",
						Triggers: []model.ActionTriggerDSLModel{
							{KeyType: "READ_RESULT", Name: "read", Then: "FAILURE"},
						},
					},
				},
			},
		},
	})

	err := newScopedValidator(t).ValidateFrame(frame)
	if err == nil {
		t.Fatal("Expected error for trigger used outside its scope")
	}
	if !strings.Contains(err.Error(), "trigger 'read' requires scope 'resultScope'") {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestNestedTriggerThenAgainstParentEvents(t *testing.T) {
	frame := newFrame(model.BlockDSLModel{
		KeyType: "nativeblocks/column",
		Key:     "column",
		Slot:    "content",
		Actions: []model.ActionDSLModel{
			{
				Key:   "column",
				Event: "onClick",
				Triggers: []model.ActionTriggerDSLModel{
					{
						KeyType: "SCRIPT",
						Name:    "run",
						Then:    "NEXT",
						Triggers: []model.ActionTriggerDSLModel{
							{KeyType: "SCRIPT", Name: "after", Then: "ALWAYS"},
						},
					},
				},
			},
		},
	})

	err := newScopedValidator(t).ValidateFrame(frame)
	if err == nil {
		t.Fatal("Expected error for then value the parent action integration does not declare")
	}
	if !strings.Contains(err.Error(), "trigger 'after' uses invalid then 'ALWAYS'") {
		t.Errorf("Unexpected error: %v", err)
	}
}
