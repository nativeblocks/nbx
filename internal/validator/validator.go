package validator

import (
	"fmt"
	"strings"

	"github.com/nativeblocks/nbx/internal/errors"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/nativeblocks/nbx/internal/types"
)

type Validator struct {
	frame          *model.FrameDSLModel
	errorCollector *errors.ErrorCollector
	variables      map[string]variableInfo
	blockKeys      map[string]int
	actionKeys     map[string]int
	slotNames      map[string]bool
}

type variableInfo struct {
	varType types.Type
	line    int
	used    bool
}

func NewValidator(frame *model.FrameDSLModel, source string) *Validator {
	return &Validator{
		frame:          frame,
		errorCollector: errors.NewErrorCollector(source),
		variables:      make(map[string]variableInfo),
		blockKeys:      make(map[string]int),
		actionKeys:     make(map[string]int),
		slotNames:      make(map[string]bool),
	}
}

func (v *Validator) _validate() (*errors.ErrorCollector, error) {
	v._collectVariables()
	v._collectBlockKeys()
	v._collectSlots()

	v._validateFrame()
	v._validateBlocks(v.frame.Blocks)

	v._checkUnusedVariables()

	return v.errorCollector, nil
}

func (v *Validator) _collectVariables() {
	for _, variable := range v.frame.Variables {
		if existing, exists := v.variables[variable.Key]; exists {
			v.errorCollector.AddError(errors.DuplicateDeclarationError(
				variable.Key, variable.Line, variable.Column, existing.line,
			))
			continue
		}

		varType, err := types.FromString(variable.Type)
		if err != nil {
			v.errorCollector.AddSimpleError(
				fmt.Sprintf("Unknown type '%s' for variable '%s'", variable.Type, variable.Key),
				variable.Line, variable.Column,
			)
			varType = types.TypeUnknown
		}

		if varType != types.TypeUnknown {
			if valid, msg := types.ValidateValue(variable.Value, varType); !valid {
				v.errorCollector.AddSimpleError(
					fmt.Sprintf("Invalid initial value for variable '%s': %s", variable.Key, msg),
					variable.Line, variable.Column,
				)
			}
		}

		v.variables[variable.Key] = variableInfo{
			varType: varType,
			line:    variable.Line,
			used:    false,
		}
	}
}

func (v *Validator) _collectBlockKeys() {
	v._collectBlockKeysRecursive(v.frame.Blocks)
}

func (v *Validator) _collectBlockKeysRecursive(blocks []model.BlockDSLModel) {
	for _, block := range blocks {
		if firstLine, exists := v.blockKeys[block.Key]; exists {
			v.errorCollector.AddSimpleError(
				fmt.Sprintf("Duplicate block key '%s' (first declared at line %d)", block.Key, firstLine),
				block.Line, block.Column,
			)
		} else {
			v.blockKeys[block.Key] = block.Line
		}

		v._collectBlockKeysRecursive(block.Blocks)
	}
}

func (v *Validator) _collectSlots() {
	v._collectSlotsRecursive(v.frame.Blocks)
}

func (v *Validator) _collectSlotsRecursive(blocks []model.BlockDSLModel) {
	for _, block := range blocks {
		for _, slot := range block.Slots {
			v.slotNames[slot.Slot] = true
		}

		v._collectSlotsRecursive(block.Blocks)
	}
}

func (v *Validator) _validateFrame() {
	if v.frame.Name == "" {
		v.errorCollector.AddSimpleError(
			"Frame 'name' attribute is required",
			v.frame.Line, v.frame.Column,
		)
	}

	if v.frame.Route == "" {
		v.errorCollector.AddSimpleError(
			"Frame 'route' attribute is required",
			v.frame.Line, v.frame.Column,
		)
	}

	if v.frame.Type != "FRAME" && v.frame.Type != "BOTTOM_SHEET" && v.frame.Type != "DIALOG" {
		v.errorCollector.AddWarning(
			fmt.Sprintf("Unexpected frame type '%s'. Valid types: FRAME, BOTTOM_SHEET, DIALOG", v.frame.Type),
			v.frame.Line, v.frame.Column,
			"Consider using one of the standard frame types",
		)
	}
}

func (v *Validator) _validateBlocks(blocks []model.BlockDSLModel) {
	for _, block := range blocks {
		v._validateBlock(&block)
	}
}

func (v *Validator) _validateBlock(block *model.BlockDSLModel) {
	if block.KeyType == "" {
		v.errorCollector.AddSimpleError(
			fmt.Sprintf("Block '%s' is missing required 'keyType' attribute", block.Key),
			block.Line, block.Column,
		)
	}

	if block.Key == "" {
		v.errorCollector.AddSimpleError(
			"Block is missing required 'key' attribute",
			block.Line, block.Column,
		)
	}

	if block.VisibilityKey != "" {
		v._validateVariableReference(block.VisibilityKey, 0, 0)
	}

	for _, data := range block.Data {
		v._validateDataBinding(data.Value, data.Line, data.Column)
	}

	for _, action := range block.Actions {
		v._validateAction(&action, block.Key)
	}

	v._validateBlocks(block.Blocks)
}

func (v *Validator) _validateAction(action *model.ActionDSLModel, blockKey string) {
	if action.Event == "" {
		v.errorCollector.AddSimpleError(
			fmt.Sprintf("Action in block '%s' is missing required 'event' attribute", blockKey),
			action.Line, action.Column,
		)
	}

	for _, trigger := range action.Triggers {
		v._validateTrigger(&trigger)
	}
}

func (v *Validator) _validateTrigger(trigger *model.ActionTriggerDSLModel) {
	if trigger.KeyType == "" {
		v.errorCollector.AddSimpleError(
			fmt.Sprintf("Trigger '%s' is missing required 'keyType' attribute", trigger.Name),
			trigger.Line, trigger.Column,
		)
	}

	validThenValues := []string{"NEXT", "SUCCESS", "FAILURE", "ALWAYS"}
	if trigger.Then != "" {
		valid := false
		for _, validValue := range validThenValues {
			if trigger.Then == validValue {
				valid = true
				break
			}
		}
		if !valid {
			v.errorCollector.AddSimpleError(
				fmt.Sprintf("Trigger '%s' has unexpected 'then' value '%s'", trigger.Name, trigger.Then),
				trigger.Line, trigger.Column,
			)
		}
	}

	for _, data := range trigger.Data {
		v._validateDataBinding(data.Value, data.Line, data.Column)
	}

	for _, nestedTrigger := range trigger.Triggers {
		v._validateTrigger(&nestedTrigger)
	}
}

func (v *Validator) _validateVariableReference(varName string, line int, column int) {
	if varInfo, exists := v.variables[varName]; exists {
		varInfo.used = true
		v.variables[varName] = varInfo
	} else {
		availableVars := make([]string, 0, len(v.variables))
		for varName := range v.variables {
			availableVars = append(availableVars, varName)
		}
		// Note: Variable references in data bindings don't have their own line/column
		// We'd need to track where the reference occurs, not where the variable is declared
		v.errorCollector.AddError(errors.UndefinedVariableError(varName, line, column, availableVars))
	}
}

func (v *Validator) _validateDataBinding(value string, line int, column int) {
	trimmed := strings.TrimSpace(value)
	if _isVariableName(trimmed) {
		v._validateVariableReference(trimmed, line, column)
	}
}

func (v *Validator) _checkUnusedVariables() {
	for varName, info := range v.variables {
		if !info.used {
			v.errorCollector.AddWarning(
				fmt.Sprintf("Variable '%s' is declared but never used", varName),
				info.line, 0,
				fmt.Sprintf("Remove the unused variable or use it in your blocks"),
			)
		}
	}
}

func _isVariableName(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, ch := range s {
		if i == 0 {
			if !_isLetter(ch) {
				return false
			}
		} else {
			if !_isLetter(ch) && !_isDigit(ch) && ch != '_' {
				return false
			}
		}
	}

	return true
}

func _isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func _isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func ValidateWithSource(frame *model.FrameDSLModel, source string) (*errors.ErrorCollector, error) {
	validator := NewValidator(frame, source)
	return validator._validate()
}

func Validate(frame *model.FrameDSLModel) (*errors.ErrorCollector, error) {
	return ValidateWithSource(frame, "")
}
