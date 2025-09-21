package compiler

import (
	"fmt"
	"strings"

	"github.com/nativeblocks/nbx/internal/model"
)

// ToString converts a FrameDSLModel back to the original NBX DSL string format
func ToString(frame model.FrameDSLModel) string {
	var builder strings.Builder

	builder.WriteString("frame(\n")
	builder.WriteString(fmt.Sprintf("    name = \"%s\",\n", frame.Name))
	builder.WriteString(fmt.Sprintf("    route = \"%s\"\n", frame.Route))
	builder.WriteString(") {\n")

	for _, variable := range frame.Variables {
		builder.WriteString(fmt.Sprintf("    var %s: %s = %s\n",
			variable.Key,
			variable.Type,
			formatVariableValue(variable.Value, variable.Type)))
	}

	if len(frame.Variables) > 0 {
		builder.WriteString("\n")
	}

	for _, block := range frame.Blocks {
		formatBlock(&builder, block, 1)
	}

	builder.WriteString("}")
	return builder.String()
}

func formatVariableValue(value, valueType string) string {
	switch valueType {
	case "STRING":
		return fmt.Sprintf("\"%s\"", value)
	case "BOOLEAN", "INT", "LONG", "FLOAT", "DOUBLE":
		return value
	default:
		return fmt.Sprintf("\"%s\"", value)
	}
}

func formatBlock(builder *strings.Builder, block model.BlockDSLModel, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	builder.WriteString(fmt.Sprintf("%sblock(keyType = \"%s\", key = \"%s\"", indent, block.KeyType, block.Key))
	if block.VisibilityKey != "" {
		builder.WriteString(fmt.Sprintf(", visibility = %s", block.VisibilityKey))
	}
	if block.IntegrationVersion > 0 {
		builder.WriteString(fmt.Sprintf(", version = %d", block.IntegrationVersion))
	}

	builder.WriteString(")")

	if len(block.Properties) > 0 {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("%s.prop(\n", indent))

		for i, prop := range block.Properties {
			propIndent := strings.Repeat("    ", indentLevel+1)

			if hasMultipleDeviceValues(prop) {
				builder.WriteString(fmt.Sprintf("%s%s = \"%s\"", propIndent, prop.Key, prop.ValueMobile))
			} else {
				value := getSinglePropertyValue(prop)
				builder.WriteString(fmt.Sprintf("%s%s = \"%s\"", propIndent, prop.Key, value))
			}

			if i < len(block.Properties)-1 {
				builder.WriteString(",")
			}
			builder.WriteString("\n")
		}

		builder.WriteString(fmt.Sprintf("%s)", indent))
	}

	if len(block.Data) > 0 {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("%s.data(", indent))

		for i, data := range block.Data {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("%s = %s", data.Key, data.Value))
		}

		builder.WriteString(")")
	}

	for _, action := range block.Actions {
		formatAction(builder, action, indentLevel)
	}

	if len(block.Blocks) > 0 || len(block.Slots) > 0 {
		slotBlocks := make(map[string][]model.BlockDSLModel)

		for _, childBlock := range block.Blocks {
			slotName := childBlock.Slot
			if (slotName == "" || slotName == "null") && childBlock.KeyType != "ROOT" {
				slotName = "content"
			}
			slotBlocks[slotName] = append(slotBlocks[slotName], childBlock)
		}

		for _, slot := range block.Slots {
			if _, exists := slotBlocks[slot.Slot]; !exists {
				slotBlocks[slot.Slot] = []model.BlockDSLModel{}
			}
		}

		for slotName, blocks := range slotBlocks {
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("%s.slot(\"%s\") {\n", indent, slotName))

			for _, childBlock := range blocks {
				formatBlock(builder, childBlock, indentLevel+1)
			}

			builder.WriteString(fmt.Sprintf("%s}\n", indent))
		}
	} else {
		builder.WriteString("\n")
	}
}

func hasMultipleDeviceValues(prop model.BlockPropertyDSLModel) bool {
	count := 0
	if prop.ValueMobile != "" {
		count++
	}
	if prop.ValueTablet != "" {
		count++
	}
	if prop.ValueDesktop != "" {
		count++
	}
	return count > 1
}

func getSinglePropertyValue(prop model.BlockPropertyDSLModel) string {
	if prop.ValueMobile != "" {
		return prop.ValueMobile
	}
	if prop.ValueTablet != "" {
		return prop.ValueTablet
	}
	if prop.ValueDesktop != "" {
		return prop.ValueDesktop
	}
	return ""
}

func formatAction(builder *strings.Builder, action model.ActionDSLModel, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("%s.action(event = \"%s\") {\n", indent, action.Event))

	for _, trigger := range action.Triggers {
		formatTrigger(builder, trigger, indentLevel+1)
	}

	builder.WriteString(fmt.Sprintf("%s}\n", indent))
}

func formatTrigger(builder *strings.Builder, trigger model.ActionTriggerDSLModel, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	builder.WriteString(fmt.Sprintf("%strigger(keyType = \"%s\", name = \"%s\"",
		indent, trigger.KeyType, trigger.Name))

	if trigger.IntegrationVersion > 0 {
		builder.WriteString(fmt.Sprintf(", version = %d", trigger.IntegrationVersion))
	}

	builder.WriteString(")")

	if len(trigger.Properties) > 0 {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("%s.prop(", indent))

		for i, prop := range trigger.Properties {
			if i > 0 {
				builder.WriteString(", ")
			}

			if strings.Contains(prop.Value, "#SCRIPT") {
				builder.WriteString(fmt.Sprintf("%s = \"%s\"", prop.Key, prop.Value))
			} else {
				builder.WriteString(fmt.Sprintf("%s = \"%s\"", prop.Key, prop.Value))
			}
		}

		builder.WriteString(")")
	}

	if len(trigger.Data) > 0 {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("%s.data(", indent))

		for i, data := range trigger.Data {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("%s = %s", data.Key, data.Value))
		}

		builder.WriteString(")")
	}

	if len(trigger.Triggers) > 0 {
		builder.WriteString("\n")
		for _, nestedTrigger := range trigger.Triggers {
			formatTrigger(builder, nestedTrigger, indentLevel+1)
		}
	}

	builder.WriteString("\n")
}
