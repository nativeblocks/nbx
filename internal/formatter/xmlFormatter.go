package formatter

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/nativeblocks/nbx/internal/errors"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/nativeblocks/nbx/internal/parser"
)

func FormatXML(xmlString string) (string, []*errors.Error) {
	frame, errs := parser.ParseXML(xmlString)
	if len(errs) > 0 {
		return "", errs
	}

	return FormatFrameXML(frame), nil
}

func FormatFrameXML(frame model.FrameDSLModel) string {
	var builder strings.Builder

	builder.WriteString(xml.Header)

	builder.WriteString(fmt.Sprintf("<frame name=%q route=%q", frame.Name, frame.Route))
	if frame.Type != "" && frame.Type != "FRAME" {
		builder.WriteString(fmt.Sprintf(" type=%q", frame.Type))
	}
	builder.WriteString(">\n")

	for _, v := range frame.Variables {
		builder.WriteString(fmt.Sprintf("  <var key=%q type=%q value=%q />\n",
			v.Key, v.Type, _escapeXML(v.Value)))
	}

	if len(frame.Variables) > 0 && len(frame.Blocks) > 0 {
		builder.WriteString("\n")
	}

	for _, b := range frame.Blocks {
		_formatBlock(&builder, b, 1)
	}

	builder.WriteString("</frame>\n")
	return builder.String()
}

func _formatBlock(builder *strings.Builder, block model.BlockDSLModel, indent int) {
	ind := strings.Repeat("  ", indent)

	builder.WriteString(fmt.Sprintf("%s<block keyType=%q key=%q",
		ind, _escapeXML(block.KeyType), _escapeXML(block.Key)))

	if block.VisibilityKey != "" && block.VisibilityKey != "null" {
		builder.WriteString(fmt.Sprintf(" visibility=%q", _escapeXML(block.VisibilityKey)))
	}
	if block.IntegrationVersion > 0 {
		builder.WriteString(fmt.Sprintf(" version=%q", fmt.Sprint(block.IntegrationVersion)))
	}
	builder.WriteString(">\n")

	for _, p := range block.Properties {
		_formatProperty(builder, p, indent+1)
	}

	for _, d := range block.Data {
		builder.WriteString(fmt.Sprintf("%s  <data key=%q value=%q />\n",
			ind, _escapeXML(d.Key), _escapeXML(d.Value)))
	}

	for _, a := range block.Actions {
		_formatAction(builder, a, indent+1)
	}

	slotMap := make(map[string][]model.BlockDSLModel)
	for _, child := range block.Blocks {
		slotName := child.Slot
		if slotName == "" || slotName == "null" {
			slotName = "content"
		}
		slotMap[slotName] = append(slotMap[slotName], child)
	}

	for _, slot := range block.Slots {
		if blocks, ok := slotMap[slot.Slot]; ok {
			builder.WriteString(fmt.Sprintf("%s  <slot name=%q>\n", ind, _escapeXML(slot.Slot)))
			for _, child := range blocks {
				_formatBlock(builder, child, indent+2)
			}
			builder.WriteString(fmt.Sprintf("%s  </slot>\n", ind))
		}
	}

	builder.WriteString(fmt.Sprintf("%s</block>\n", ind))
}

func _formatProperty(builder *strings.Builder, prop model.BlockPropertyDSLModel, indent int) {
	ind := strings.Repeat("  ", indent)

	mobile := prop.ValueMobile
	tablet := prop.ValueTablet
	desktop := prop.ValueDesktop

	// If all values are the same, use single value attribute
	if mobile == tablet && tablet == desktop {
		// Check if value contains newlines (multiline)
		if strings.Contains(mobile, "\n") {
			builder.WriteString(fmt.Sprintf("%s<prop key=\"%s\"\n%s      value=\"%s\" />\n",
				ind, _escapeXML(prop.Key), ind, _escapeXML(mobile)))
		} else {
			builder.WriteString(fmt.Sprintf("%s<prop key=\"%s\" value=\"%s\" />\n",
				ind, _escapeXML(prop.Key), _escapeXML(mobile)))
		}
	} else {
		// Use device-specific attributes
		builder.WriteString(fmt.Sprintf("%s<prop key=\"%s\"", ind, _escapeXML(prop.Key)))

		if mobile != "" {
			if strings.Contains(mobile, "\n") {
				builder.WriteString(fmt.Sprintf("\n%s      mobile=\"%s\"", ind, _escapeXML(mobile)))
			} else {
				builder.WriteString(fmt.Sprintf(" mobile=\"%s\"", _escapeXML(mobile)))
			}
		}
		if tablet != "" {
			if strings.Contains(tablet, "\n") {
				builder.WriteString(fmt.Sprintf("\n%s      tablet=\"%s\"", ind, _escapeXML(tablet)))
			} else {
				builder.WriteString(fmt.Sprintf(" tablet=\"%s\"", _escapeXML(tablet)))
			}
		}
		if desktop != "" {
			if strings.Contains(desktop, "\n") {
				builder.WriteString(fmt.Sprintf("\n%s      desktop=\"%s\"", ind, _escapeXML(desktop)))
			} else {
				builder.WriteString(fmt.Sprintf(" desktop=\"%s\"", _escapeXML(desktop)))
			}
		}

		builder.WriteString(" />\n")
	}
}

func _formatAction(builder *strings.Builder, action model.ActionDSLModel, indent int) {
	ind := strings.Repeat("  ", indent)

	builder.WriteString(fmt.Sprintf("%s<action event=%q>\n", ind, _escapeXML(action.Event)))

	for _, trigger := range action.Triggers {
		_formatTrigger(builder, trigger, indent+1, "NEXT")
	}

	builder.WriteString(fmt.Sprintf("%s</action>\n", ind))
}

func _formatTrigger(builder *strings.Builder, trigger model.ActionTriggerDSLModel, indent int, defaultThen string) {
	ind := strings.Repeat("  ", indent)

	builder.WriteString(fmt.Sprintf("%s<trigger keyType=%q name=%q",
		ind, _escapeXML(trigger.KeyType), _escapeXML(trigger.Name)))

	if trigger.IntegrationVersion > 0 {
		builder.WriteString(fmt.Sprintf(" version=%q", fmt.Sprint(trigger.IntegrationVersion)))
	}
	builder.WriteString(">\n")

	// Properties
	for _, p := range trigger.Properties {
		if strings.Contains(p.Value, "\n") {
			builder.WriteString(fmt.Sprintf("%s  <prop key=\"%s\"\n%s        value=\"%s\" />\n",
				ind, _escapeXML(p.Key), ind, _escapeXML(p.Value)))
		} else {
			builder.WriteString(fmt.Sprintf("%s  <prop key=\"%s\" value=\"%s\" />\n",
				ind, _escapeXML(p.Key), _escapeXML(p.Value)))
		}
	}

	for _, d := range trigger.Data {
		builder.WriteString(fmt.Sprintf("%s  <data key=%q value=%q />\n",
			ind, _escapeXML(d.Key), _escapeXML(d.Value)))
	}

	thenMap := make(map[string][]model.ActionTriggerDSLModel)
	for _, nested := range trigger.Triggers {
		thenMap[nested.Then] = append(thenMap[nested.Then], nested)
	}

	for thenValue, nestedTriggers := range thenMap {
		if thenValue != "" && thenValue != defaultThen {
			builder.WriteString(fmt.Sprintf("%s  <then value=%q>\n", ind, thenValue))
			for _, nested := range nestedTriggers {
				_formatTrigger(builder, nested, indent+2, thenValue)
			}
			builder.WriteString(fmt.Sprintf("%s  </then>\n", ind))
		} else {
			for _, nested := range nestedTriggers {
				_formatTrigger(builder, nested, indent+1, thenValue)
			}
		}
	}

	builder.WriteString(fmt.Sprintf("%s</trigger>\n", ind))
}

func _escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
