package formatter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/nativeblocks/nbx/internal/lexer"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/nativeblocks/nbx/internal/parser"
)

func Format(dslString string) (string, error) {
	frame, err := _parseToFrameDSL(dslString)
	if err != nil {
		return "", err
	}

	return FormatFrameDSL(frame), nil
}

func FormatFrameDSL(frame model.FrameDSLModel) string {
	var builder strings.Builder

	builder.WriteString("frame(\n")
	builder.WriteString(fmt.Sprintf("    name = \"%s\",\n", frame.Name))
	builder.WriteString(fmt.Sprintf("    route = \"%s\"\n", frame.Route))
	builder.WriteString(") {\n")

	for _, variable := range frame.Variables {
		builder.WriteString(fmt.Sprintf("    var %s: %s = %s\n",
			variable.Key,
			variable.Type,
			_formatVariableValueConsistent(variable.Value, variable.Type)))
	}

	if len(frame.Variables) > 0 {
		builder.WriteString("\n")
	}

	for _, block := range frame.Blocks {
		_formatBlockConsistent(&builder, block, 1)
	}

	builder.WriteString("}")
	return builder.String()
}

func _formatVariableValueConsistent(value, valueType string) string {
	switch valueType {
	case "STRING":
		return fmt.Sprintf("\"%s\"", value)
	case "BOOLEAN", "INT", "LONG", "FLOAT", "DOUBLE":
		return value
	default:
		return fmt.Sprintf("\"%s\"", value)
	}
}

func _formatBlockConsistent(builder *strings.Builder, block model.BlockDSLModel, indentLevel int) {
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
			value := _getSinglePropertyValueConsistent(prop)
			builder.WriteString(fmt.Sprintf("%s%s = \"%s\"", propIndent, prop.Key, value))

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
		_formatActionConsistent(builder, action, indentLevel)
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
			if slotName != "" && slotName != "null" {
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("%s.slot(\"%s\") {\n", indent, slotName))

				for _, childBlock := range blocks {
					_formatBlockConsistent(builder, childBlock, indentLevel+1)
				}

				builder.WriteString(fmt.Sprintf("%s}\n", indent))
			}
		}
	} else {
		builder.WriteString("\n")
	}
}

func _getSinglePropertyValueConsistent(prop model.BlockPropertyDSLModel) string {
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

func _formatActionConsistent(builder *strings.Builder, action model.ActionDSLModel, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("%s.action(event = \"%s\") {\n", indent, action.Event))

	for _, trigger := range action.Triggers {
		_formatTriggerConsistent(builder, trigger, indentLevel+1)
	}

	builder.WriteString(fmt.Sprintf("%s}\n", indent))
}

func _formatTriggerConsistent(builder *strings.Builder, trigger model.ActionTriggerDSLModel, indentLevel int) {
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
				formattedScript := _formatScriptBlock(prop.Value, indentLevel+1)
				builder.WriteString(fmt.Sprintf("%s = \"%s\"", prop.Key, formattedScript))
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
			_formatTriggerConsistent(builder, nestedTrigger, indentLevel+1)
		}
	}

	builder.WriteString("\n")
}

func _formatScriptBlock(script string, baseIndentLevel int) string {
	if !strings.Contains(script, "#SCRIPT") {
		return script
	}

	scriptRegex := regexp.MustCompile(`#SCRIPT\s*(.*?)\s*#ENDSCRIPT`)
	match := scriptRegex.FindStringSubmatch(script)
	if len(match) < 2 {
		return script
	}

	scriptContent := strings.TrimSpace(match[1])
	lines := strings.Split(scriptContent, "\n")

	var formattedLines []string
	scriptIndent := strings.Repeat("    ", baseIndentLevel)

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			formattedLines = append(formattedLines, scriptIndent+trimmedLine)
		}
	}

	return fmt.Sprintf("#SCRIPT\n%s\n%s#ENDSCRIPT",
		strings.Join(formattedLines, "\n"),
		strings.Repeat("    ", baseIndentLevel))
}

func _parseToFrameDSL(dslString string) (model.FrameDSLModel, error) {
	l := lexer.NewLexer(dslString)
	p := parser.NewParser(l, dslString)
	frame := p.ParseNBX()

	errorCollector := p.ErrorCollector()
	if frame == nil || errorCollector.HasErrors() {
		if errorCollector.HasErrors() {
			return model.FrameDSLModel{}, errors.New(errorCollector.FormatAll())
		}
		return model.FrameDSLModel{}, errors.New("failed to parse DSL")
	}
	return *frame, nil
}
