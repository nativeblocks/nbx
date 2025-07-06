package compiler

import (
	"fmt"
	"github.com/nativeblocks/nbx/internal/model"
	"strconv"
	"strings"
)

type nbxValidationError struct {
	Field       string
	Description string
	Value       interface{}
	Type        string
}

func formatValidationErrors(errs []nbxValidationError, frame model.FrameDSLModel) []string {
	var out []string

	for _, e := range errs {
		desc := e.Description
		if strings.HasPrefix(desc, e.Field) {
			desc = strings.TrimPrefix(desc, e.Field)
		} else if idx := strings.Index(desc, " "); idx != -1 {
			desc = desc[idx:]
		}

		ctx := buildContext(strings.Split(e.Field, "."), frame)
		if ctx != "" {
			out = append(out, fmt.Sprintf("%s: %v%s", ctx, e.Value, desc))
		} else {
			out = append(out, fmt.Sprintf("%v%s", e.Value, desc))
		}
	}
	return out
}

func buildContext(segments []string, frame model.FrameDSLModel) string {
	var (
		currentBlock  model.BlockDSLModel
		currentAction model.ActionDSLModel
		currentTrig   model.ActionTriggerDSLModel

		blkName  string
		actKey   string
		trigName string
		propKey  string
		dataKey  string
		slotKey  string
		blocks   = frame.Blocks
	)

	for i := 0; i < len(segments)-1; i++ {
		switch segments[i] {
		case "blocks":
			idx, err := strconv.Atoi(segments[i+1])
			if err != nil || idx < 0 || idx >= len(blocks) {
				return ""
			}
			currentBlock = blocks[idx]
			blkName = currentBlock.Key
			blocks = currentBlock.Blocks
			i++

		case "actions":
			idx, err := strconv.Atoi(segments[i+1])
			if err != nil || idx < 0 || idx >= len(currentBlock.Actions) {
				break
			}
			currentAction = currentBlock.Actions[idx]
			actKey = currentAction.Key
			i++

		case "triggers":
			idx, err := strconv.Atoi(segments[i+1])
			if err != nil || idx < 0 || idx >= len(currentAction.Triggers) {
				break
			}
			currentTrig = currentAction.Triggers[idx]
			trigName = currentTrig.Name
			i++

		case "properties":
			idx, err := strconv.Atoi(segments[i+1])
			if err != nil || idx < 0 || idx >= len(currentBlock.Properties) {
				break
			}
			propKey = currentBlock.Properties[idx].Key
			i++

		case "data":
			idx, err := strconv.Atoi(segments[i+1])
			if err != nil {
				break
			}
			if currentTrig.KeyType != "" && idx < len(currentTrig.Data) {
				dataKey = currentTrig.Data[idx].Key
			} else if idx < len(currentBlock.Data) {
				dataKey = currentBlock.Data[idx].Key
			}
			i++

		case "slots":
			idx, err := strconv.Atoi(segments[i+1])
			if err != nil || idx < 0 || idx >= len(currentBlock.Slots) {
				break
			}
			slotKey = currentBlock.Slots[idx].Slot
			i++
		}
	}

	switch {
	case trigName != "":
		return fmt.Sprintf(`trigger "%s"`, trigName)
	case actKey != "":
		return fmt.Sprintf(`action "%s"`, actKey)
	case propKey != "":
		return fmt.Sprintf(`property "%s" in block "%s"`, propKey, blkName)
	case dataKey != "":
		return fmt.Sprintf(`data "%s" in block "%s"`, dataKey, blkName)
	case slotKey != "":
		return fmt.Sprintf(`slot "%s" in block "%s"`, slotKey, blkName)
	case blkName != "":
		return fmt.Sprintf(`block "%s"`, blkName)
	default:
		return ""
	}
}
