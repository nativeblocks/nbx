package compiler

import (
	"github.com/nativeblocks/nbx/internal/model"
	"testing"
)

func TestFormatBlockError(t *testing.T) {
	blocks := []model.BlockJson{
		{
			Id:      "root-1",
			Key:     "rootBlock",
			KeyType: "ROOT",
			Properties: []model.BlockPropertyJson{
				{
					Key: "paddingStart",
				},
				{
					Key: "paddingTop",
				},
			},
			Data: []model.BlockDataJson{
				{
					Key: "text",
				},
			},
		},
		{
			Id:       "child-1",
			ParentId: "root-1",
			Key:      "childBlock1",
			Properties: []model.BlockPropertyJson{
				{
					Key: "height",
				},
			},
		},
	}

	testCases := []struct {
		name         string
		pathSegments []string
		errorMessage string
		expected     string
	}{
		{
			name:         "Simple block path",
			pathSegments: []string{"blocks", "0", "properties", "1", "key"},
			errorMessage: "must be one of: paddingTop, paddingBottom",
			expected:     "Block \"rootBlock\" has an invalid property for key \"paddingTop\": must be one of: paddingTop, paddingBottom",
		},
		{
			name:         "Child block path",
			pathSegments: []string{"blocks", "1", "properties", "0", "key"},
			errorMessage: "must be one of: height, width",
			expected:     "Block \"childBlock1\" has an invalid property for key \"height\": must be one of: height, width",
		},
		{
			name:         "Block data path",
			pathSegments: []string{"blocks", "0", "data", "0", "value"},
			errorMessage: "cannot be empty",
			expected:     "Block \"rootBlock\" has an invalid data field for key \"text\": cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := _formatBlockError(tc.pathSegments, tc.errorMessage, blocks)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected: \n%s\nGot: \n%s", tc.expected, result)
			}
		})
	}
}

func TestGetBlockKeyFromPath(t *testing.T) {
	blocks := []model.BlockJson{
		{
			Id:  "root-id",
			Key: "rootBlock",
		},
	}

	path := []string{"blocks", "0", "properties", "1", "key"}
	key := _getBlockKeyFromPath(path, blocks)

	if key != "rootBlock" {
		t.Errorf("Expected 'rootBlock', got '%s'", key)
	}
}

func TestFindDeepestBlockByPath(t *testing.T) {
	blocks := []model.BlockJson{
		{
			Id:      "root-1",
			Key:     "rootBlock",
			KeyType: "ROOT",
		},
		{
			Id:       "child-1",
			ParentId: "root-1",
			Key:      "childBlock",
			Position: 0,
		},
	}

	path := []string{"blocks", "0", "properties", "0", "key"}
	block, found := _findDeepestBlockByPath(path, blocks)

	if !found {
		t.Errorf("Expected to find block but didn't")
		return
	}

	if block.Key != "rootBlock" {
		t.Errorf("Expected 'rootBlock', got '%s'", block.Key)
	}
}

func TestGetBlockSectionAndKey(t *testing.T) {
	blocks := []model.BlockJson{
		{
			Id:  "root-1",
			Key: "rootBlock",
			Properties: []model.BlockPropertyJson{
				{
					Key: "width",
				},
			},
		},
	}

	path := []string{"blocks", "0", "properties", "0", "key"}
	section, key, err := _getBlockSectionAndKey(path, blocks)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if section != "property" {
		t.Errorf("Expected section 'property', got '%s'", section)
	}

	if key != "width" {
		t.Errorf("Expected key 'width', got '%s'", key)
	}
}

func TestFormatActionError(t *testing.T) {
	actions := []model.ActionJson{
		{
			Id:    "action-1",
			Key:   "testAction",
			Event: "click",
			Triggers: []model.ActionTriggerJson{
				{
					Id:       "trigger-1",
					ActionId: "action-1",
					Name:     "testTrigger",
					Properties: []model.TriggerPropertyJson{
						{
							Key: "delay",
						},
					},
					Data: []model.TriggerDataJson{
						{
							Key: "targetBlock",
						},
					},
				},
			},
		},
	}

	testCases := []struct {
		name         string
		pathSegments []string
		errorMessage string
		expected     string
	}{
		{
			name:         "Action trigger property error",
			pathSegments: []string{"actions", "0", "properties", "0", "key"},
			errorMessage: "must be one of: delay, duration",
			expected:     "Action \"testAction\" has an invalid property for key \"property\": must be one of: delay, duration",
		},
		{
			name:         "Action trigger data error",
			pathSegments: []string{"actions", "0", "data", "0", "value"},
			errorMessage: "cannot be empty",
			expected:     "Action \"testAction\" has an invalid data field for key \"unknown\": cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := _formatActionError(tc.pathSegments, tc.errorMessage, actions)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected: \n%s\nGot: \n%s", tc.expected, result)
			}
		})
	}
}
