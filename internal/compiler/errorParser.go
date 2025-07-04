package compiler

import (
	"fmt"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/xeipuuv/gojsonschema"
	"strconv"
	"strings"
)

func _formatValidationErrors(result *gojsonschema.Result, frameJson model.FrameJson) []string {
	var humanReadableErrors []string

	for _, err := range result.Errors() {
		errorPath := err.Field()
		errorMessage := err.Description()

		pathSegments := strings.Split(errorPath, ".")

		if len(pathSegments) < 3 {
			humanReadableErrors = append(humanReadableErrors, fmt.Sprintf("Validation error in %s: %s", errorPath, errorMessage))
			continue
		}

		switch pathSegments[0] {
		case "blocks":
			formattedError, err := _formatBlockError(pathSegments, errorMessage, frameJson.Blocks)
			if err != nil {
				humanReadableErrors = append(humanReadableErrors, fmt.Sprintf("Error processing validation error: %s", err.Error()))
			} else {
				humanReadableErrors = append(humanReadableErrors, formattedError)
			}
		case "actions":
			formattedError, err := _formatActionError(pathSegments, errorMessage, frameJson.Actions)
			if err != nil {
				humanReadableErrors = append(humanReadableErrors, fmt.Sprintf("Error processing validation error: %s", err.Error()))
			} else {
				humanReadableErrors = append(humanReadableErrors, formattedError)
			}
		case "variables":
			variableIndex := -1
			if len(pathSegments) > 1 {
				variableIndex, _ = strconv.Atoi(pathSegments[1])
			}

			variableKey := "unknown"
			if variableIndex >= 0 && variableIndex < len(frameJson.Variables) {
				variableKey = frameJson.Variables[variableIndex].Key
			}

			humanReadableErrors = append(humanReadableErrors,
				fmt.Sprintf("Variable \"%s\" has an invalid value: %s", variableKey, errorMessage))
		default:
			humanReadableErrors = append(humanReadableErrors, fmt.Sprintf("Validation error in %s: %s", errorPath, errorMessage))
		}
	}

	return humanReadableErrors
}

func _formatBlockError(pathSegments []string, errorMessage string, blocks []model.BlockJson) (string, error) {
	section, key, err := _getBlockSectionAndKey(pathSegments, blocks)
	if err != nil {
		return "", err
	}

	blockKey := _getBlockKeyFromPath(pathSegments, blocks)
	if blockKey == "" {
		return fmt.Sprintf("Block has an invalid %s for key \"%s\": %s", section, key, errorMessage), nil
	}

	var readableSection string
	switch section {
	case "property":
		readableSection = "property"
	case "data":
		readableSection = "data field"
	case "slot":
		readableSection = "slot"
	default:
		readableSection = section
	}

	return fmt.Sprintf("Block \"%s\" has an invalid %s for key \"%s\": %s", blockKey, readableSection, key, errorMessage), nil
}

func _formatActionError(pathSegments []string, errorMessage string, actions []model.ActionJson) (string, error) {
	section, key, err := _getActionSectionAndKey(pathSegments, actions)
	if err != nil {
		return "", err
	}

	actionKey := _getActionKeyFromPath(pathSegments, actions)
	if actionKey == "" {
		return fmt.Sprintf("Action has an invalid %s for key \"%s\": %s", section, key, errorMessage), nil
	}

	var readableSection string
	switch section {
	case "trigger":
		readableSection = "trigger"
	case "property":
		readableSection = "property"
	case "data":
		readableSection = "data field"
	default:
		readableSection = section
	}

	return fmt.Sprintf("Action \"%s\" has an invalid %s for key \"%s\": %s", actionKey, readableSection, key, errorMessage), nil
}

func _getBlockKeyFromPath(pathSegments []string, blocks []model.BlockJson) string {
	block, found := _findDeepestBlockByPath(pathSegments, blocks)
	if !found {
		return ""
	}
	return block.Key
}

func _getActionKeyFromPath(pathSegments []string, actions []model.ActionJson) string {
	if len(pathSegments) < 2 || pathSegments[0] != "actions" {
		return ""
	}

	indexStr := pathSegments[1]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return ""
	}

	if index < 0 || index >= len(actions) {
		return ""
	}

	return actions[index].Key
}

func _getBlockSectionAndKey(pathSegments []string, blocks []model.BlockJson) (section string, key string, err error) {
	section = "property"
	key = "unknown"

	for i := 0; i < len(pathSegments); i++ {
		if pathSegments[i] == "properties" || pathSegments[i] == "data" || pathSegments[i] == "slots" {
			if pathSegments[i] == "properties" {
				section = "property"
			} else if pathSegments[i] == "data" {
				section = "data"
			} else if pathSegments[i] == "slots" {
				section = "slot"
			}

			if i+2 < len(pathSegments) && pathSegments[i+2] == "key" {
				indexStr := pathSegments[i+1]
				index, err := strconv.Atoi(indexStr)
				if err == nil {
					block, found := _findBlockByPath(pathSegments, blocks)
					if found {
						switch section {
						case "property":
							if index >= 0 && index < len(block.Properties) {
								key = block.Properties[index].Key
							}
						case "data":
							if index >= 0 && index < len(block.Data) {
								key = block.Data[index].Key
							}
						case "slot":
							if index >= 0 && index < len(block.Slots) {
								key = block.Slots[index].Slot
							}
						}
					}
				}
			} else if i+2 < len(pathSegments) && pathSegments[i+2] == "value" {
				indexStr := pathSegments[i+1]
				index, err := strconv.Atoi(indexStr)
				if err == nil {
					block, found := _findBlockByPath(pathSegments, blocks)
					if found {
						switch section {
						case "property":
							if index >= 0 && index < len(block.Properties) {
								key = block.Properties[index].Key
							}
						case "data":
							if index >= 0 && index < len(block.Data) {
								key = block.Data[index].Key
							}
						case "slot":
							if index >= 0 && index < len(block.Slots) {
								key = block.Slots[index].Slot
							}
						}
					}
				}
			} else if i+1 < len(pathSegments) {
				key = pathSegments[i+1]
			}
			break
		}
	}

	return section, key, nil
}

func _getActionSectionAndKey(pathSegments []string, actions []model.ActionJson) (section string, key string, err error) {
	section = "trigger"
	key = "unknown"

	for i := 0; i < len(pathSegments); i++ {
		if pathSegments[i] == "triggers" || pathSegments[i] == "properties" || pathSegments[i] == "data" {
			if pathSegments[i] == "triggers" {
				section = "trigger"
			} else if pathSegments[i] == "properties" {
				section = "property"
			} else if pathSegments[i] == "data" {
				section = "data"
			}

			if i+2 < len(pathSegments) && pathSegments[i+2] == "key" {
				indexStr := pathSegments[i+1]
				_, err := strconv.Atoi(indexStr)
				if err == nil {
					if len(pathSegments) >= 2 && pathSegments[0] == "actions" {
						actionIndex, err := strconv.Atoi(pathSegments[1])
						if err == nil && actionIndex >= 0 && actionIndex < len(actions) {
							action := actions[actionIndex]

							switch section {
							case "trigger":
								if i+1 < len(pathSegments) {
									triggerIndex, err := strconv.Atoi(pathSegments[i+1])
									if err == nil && triggerIndex >= 0 && triggerIndex < len(action.Triggers) {
										key = action.Triggers[triggerIndex].Name
									}
								}
							case "property":
								key = "property"
							case "data":
								key = "data"
							}
						}
					}
				}
			} else if i+2 < len(pathSegments) && pathSegments[i+2] == "value" {
				indexStr := pathSegments[i+1]
				index, err := strconv.Atoi(indexStr)
				if err == nil {
					if len(pathSegments) >= 2 && pathSegments[0] == "actions" {
						actionIndex, err := strconv.Atoi(pathSegments[1])
						if err == nil && actionIndex >= 0 && actionIndex < len(actions) {
							action := actions[actionIndex]

							if section == "trigger" {
							} else if i >= 4 && pathSegments[i-2] == "triggers" {
								triggerIndex, err := strconv.Atoi(pathSegments[i-1])
								if err == nil && triggerIndex >= 0 && triggerIndex < len(action.Triggers) {
									trigger := action.Triggers[triggerIndex]

									switch section {
									case "property":
										if index >= 0 && index < len(trigger.Properties) {
											key = trigger.Properties[index].Key
										}
									case "data":
										if index >= 0 && index < len(trigger.Data) {
											key = trigger.Data[index].Key
										}
									}
								}
							}
						}
					}
				}
			} else if i+1 < len(pathSegments) {
				key = pathSegments[i+1]
			}
			break
		}
	}

	return section, key, nil
}

func _findBlockByPath(pathSegments []string, blocks []model.BlockJson) (model.BlockJson, bool) {
	return _findDeepestBlockByPath(pathSegments, blocks)
}

func _findDeepestBlockByPath(pathSegments []string, allBlocks []model.BlockJson) (model.BlockJson, bool) {
	var emptyBlock model.BlockJson
	var blockMap = make(map[string]model.BlockJson)
	var parentChildMap = make(map[string][]string)

	for _, block := range allBlocks {
		blockMap[block.Id] = block
		if block.ParentId != "" {
			parentChildMap[block.ParentId] = append(parentChildMap[block.ParentId], block.Id)
		}
	}

	var blockPath []int
	for i := 0; i < len(pathSegments)-1; i++ {
		if pathSegments[i] == "blocks" {
			index, err := strconv.Atoi(pathSegments[i+1])
			if err == nil {
				blockPath = append(blockPath, index)
			}
		}
	}

	if len(blockPath) == 0 {
		return emptyBlock, false
	}

	currentBlocks := allBlocks
	var currentBlock model.BlockJson

	for depth, blockIndex := range blockPath {
		if depth == 0 {
			if blockIndex < 0 || blockIndex >= len(currentBlocks) {
				return emptyBlock, false
			}
			currentBlock = currentBlocks[blockIndex]
		} else {
			childBlockIds := parentChildMap[currentBlock.Id]
			var childBlocks []model.BlockJson

			for _, childId := range childBlockIds {
				childBlocks = append(childBlocks, blockMap[childId])
			}

			if blockIndex < 0 || blockIndex >= len(childBlocks) {
				return currentBlock, true
			}

			currentBlock = childBlocks[blockIndex]
		}
	}

	return currentBlock, true
}
