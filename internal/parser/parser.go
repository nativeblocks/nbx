package parser

import (
	"fmt"
	"strconv"

	"github.com/nativeblocks/nbx/internal/lexer"
	"github.com/nativeblocks/nbx/internal/model"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: make([]string, 0, 0),
	}
	p._nextToken()
	p._nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseNBX() *model.FrameDSLModel {
	if !p._curTokenIs(lexer.TOKEN_KEYWORD) || p.curToken.Literal != "frame" {
		p.errors = append(p.errors, "Program must start with a frame declaration")
		return nil
	}
	return p._parseFrame()
}

func (p *Parser) _nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) _expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p._nextToken()
		return true
	} else {
		p._peekError(t)
		return false
	}
}

func (p *Parser) _peekError(t lexer.TokenType) {
	err := fmt.Sprintf("Line %d, Column %d: expected next token to be %v, got %v",
		p.peekToken.Line, p.peekToken.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) _curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) _peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

// _parseKeyValuePairs parses key=value pairs within parentheses
func (p *Parser) _parseKeyValuePairs() map[string]string {
	pairs := make(map[string]string)

	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errors = append(p.errors, "Expected identifier")
			return pairs
		}
		key := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return pairs
		}
		p._nextToken()
		pairs[key] = p.curToken.Literal

		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
		} else if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	return pairs
}

// _parseBlockData parses a .data() declaration and returns a list of data items
func (p *Parser) _parseBlockData() []model.BlockDataDSLModel {
	dataList := make([]model.BlockDataDSLModel, 0, 0)

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return dataList
	}

	// Handle empty data()
	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return dataList
	}

	// Parse all key=value pairs
	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errors = append(p.errors, "Expected identifier in data declaration")
			return dataList
		}
		key := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return dataList
		}
		p._nextToken()
		value := p.curToken.Literal

		dataList = append(dataList, model.BlockDataDSLModel{
			Key:   key,
			Value: value,
			Type:  "PLACEHOLDER",
		})

		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
		} else if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	// Enforce that len == cap for dataList
	if cap(dataList) != len(dataList) {
		tmp := make([]model.BlockDataDSLModel, len(dataList))
		copy(tmp, dataList)
		dataList = tmp
	}

	return dataList
}

// _parseTriggerData parses a trigger .data() declaration and returns a list of data items
func (p *Parser) _parseTriggerData() []model.TriggerDataDSLModel {
	dataList := make([]model.TriggerDataDSLModel, 0, 0)

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return dataList
	}

	// Handle empty data()
	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return dataList
	}

	// Parse all key=value pairs
	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errors = append(p.errors, "Expected identifier in data declaration")
			return dataList
		}
		key := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return dataList
		}
		p._nextToken()
		value := p.curToken.Literal

		dataList = append(dataList, model.TriggerDataDSLModel{
			Key:   key,
			Value: value,
			Type:  "PLACEHOLDER",
		})

		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
		} else if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	// Enforce that len == cap for dataList
	if cap(dataList) != len(dataList) {
		tmp := make([]model.TriggerDataDSLModel, len(dataList))
		copy(tmp, dataList)
		dataList = tmp
	}

	return dataList
}

// _parseBlockProperty parses a block .prop() declaration and returns a list of properties
func (p *Parser) _parseBlockProperty() []model.BlockPropertyDSLModel {
	propList := make([]model.BlockPropertyDSLModel, 0, 0)

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return propList
	}

	// Handle empty prop()
	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return propList
	}

	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()

		// Handle trailing comma - if we hit closing paren after a comma, break
		if p._curTokenIs(lexer.TOKEN_RPAREN) {
			break
		}

		// Property key
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errors = append(p.errors, "Expected identifier in prop declaration")
			return propList
		}
		propKey := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return propList
		}

		// Value: either a group (for multi-device) or a single value
		p._nextToken()
		if p.curToken.Type == lexer.TOKEN_LPAREN {
			// Handle per-device or grouped values
			deviceValues := make(map[string]string)
			for {
				p._nextToken()

				if p.curToken.Type == lexer.TOKEN_RPAREN {
					break
				}

				if p.curToken.Type != lexer.TOKEN_IDENT {
					p.errors = append(p.errors, "Expected identifier inside property value parenthesis")
					return propList
				}

				deviceKey := p.curToken.Literal
				if !p._expectPeek(lexer.TOKEN_ASSIGN) {
					return propList
				}
				p._nextToken()
				deviceValue := p.curToken.Literal
				deviceValues[deviceKey] = deviceValue

				if p._peekTokenIs(lexer.TOKEN_COMMA) {
					p._nextToken()
					continue
				} else if p._peekTokenIs(lexer.TOKEN_RPAREN) {
					// Not advancing here; for will step and break
					continue
				}
			}

			// Assign values with priority order
			mobile := ""
			tablet := ""
			desktop := ""

			if v, ok := deviceValues["value"]; ok {
				mobile = v
				tablet = v
				desktop = v
			}
			if v, ok := deviceValues["valueMobile"]; ok {
				mobile = v
			}
			if v, ok := deviceValues["mobile"]; ok {
				mobile = v
			}
			if v, ok := deviceValues["valueTablet"]; ok {
				tablet = v
			}
			if v, ok := deviceValues["tablet"]; ok {
				tablet = v
			}
			if v, ok := deviceValues["valueDesktop"]; ok {
				desktop = v
			}
			if v, ok := deviceValues["desktop"]; ok {
				desktop = v
			}

			propList = append(propList, model.BlockPropertyDSLModel{
				Key:          propKey,
				ValueMobile:  mobile,
				ValueTablet:  tablet,
				ValueDesktop: desktop,
				Type:         "PLACEHOLDER",
			})
		} else {
			// Single value, possibly multiline, parse everything up to next ',' or ')'
			value := p.curToken.Literal
			for {
				// If at end or comma (i.e end of prop), break.
				if p._peekTokenIs(lexer.TOKEN_COMMA) || p._peekTokenIs(lexer.TOKEN_RPAREN) || p._peekTokenIs(lexer.TOKEN_EOF) {
					break
				}
				// Otherwise, accumulate new tokens as part of the value (including newlines, etc)
				p._nextToken()
				// Separate by space or newline character for clarity. If token type was NEWLINE (optional), add \n, else add a space.
				// But for now, we just add Literal with a newline.
				// Go lexers rarely make a newline token so we just use the literal.
				value += "\n" + p.curToken.Literal
			}
			propList = append(propList, model.BlockPropertyDSLModel{
				Key:          propKey,
				ValueMobile:  value,
				ValueTablet:  value,
				ValueDesktop: value,
				Type:         "PLACEHOLDER",
			})
		}

		// Now step to comma or paren (prop separator)
		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
			continue
		}
		if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	// Enforce that len == cap for propList
	if cap(propList) != len(propList) {
		tmp := make([]model.BlockPropertyDSLModel, len(propList))
		copy(tmp, propList)
		propList = tmp
	}

	return propList
}

// _parseTriggerProperty parses a trigger .prop() declaration and returns a list of properties
func (p *Parser) _parseTriggerProperty() []model.TriggerPropertyDSLModel {
	propList := make([]model.TriggerPropertyDSLModel, 0, 0)

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return propList
	}

	// Handle empty prop()
	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return propList
	}

	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()

		// Handle trailing comma - if we hit closing paren after a comma, break
		if p._curTokenIs(lexer.TOKEN_RPAREN) {
			break
		}

		// Property key
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errors = append(p.errors, "Expected identifier in prop declaration")
			return propList
		}
		propKey := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return propList
		}

		// Only single value is supported
		p._nextToken()
		value := p.curToken.Literal
		for {
			// If at end or comma (i.e end of prop), break.
			if p._peekTokenIs(lexer.TOKEN_COMMA) || p._peekTokenIs(lexer.TOKEN_RPAREN) || p._peekTokenIs(lexer.TOKEN_EOF) {
				break
			}
			// Otherwise, accumulate new tokens as part of the value (including newlines, etc)
			p._nextToken()
			value += "\n" + p.curToken.Literal
		}
		propList = append(propList, model.TriggerPropertyDSLModel{
			Key:   propKey,
			Value: value,
			Type:  "PLACEHOLDER",
		})

		// Now step to comma or paren (prop separator)
		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
			continue
		}
		if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	// Enforce that len == cap for propList
	if cap(propList) != len(propList) {
		tmp := make([]model.TriggerPropertyDSLModel, len(propList))
		copy(tmp, propList)
		propList = tmp
	}

	return propList
}

// _parseSlot parses a .slot() declaration
func (p *Parser) _parseSlot(block *model.BlockDSLModel) {
	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return
	}
	if !p._expectPeek(lexer.TOKEN_STRING) {
		return
	}
	slotName := p.curToken.Literal
	if !p._expectPeek(lexer.TOKEN_RPAREN) {
		return
	}
	if !p._expectPeek(lexer.TOKEN_LBRACE) {
		return
	}
	p._nextToken() // Move to first token inside the slot

	// Create a slot
	slot := model.BlockSlotDSLModel{Slot: slotName}
	block.Slots = append(block.Slots, slot)
	// Enforce that len == cap for block.Slots
	if cap(block.Slots) != len(block.Slots) {
		tmp := make([]model.BlockSlotDSLModel, len(block.Slots))
		copy(tmp, block.Slots)
		block.Slots = tmp
	}

	// Parse child blocks in the slot
	for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
		if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "block" {
			child := p._parseBlock()
			if child != nil {
				child.Slot = slotName
				block.Blocks = append(block.Blocks, *child)
				// Enforce that len == cap for block.Blocks
				if cap(block.Blocks) != len(block.Blocks) {
					tmp := make([]model.BlockDSLModel, len(block.Blocks))
					copy(tmp, block.Blocks)
					block.Blocks = tmp
				}
			}
		}
		p._nextToken()
	}
}

// _parseThen parses a .then() declaration for triggers
func (p *Parser) _parseThen(trigger *model.ActionTriggerDSLModel) {
	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return
	}
	if !p._expectPeek(lexer.TOKEN_STRING) {
		return
	}
	thenValue := p.curToken.Literal
	if !p._expectPeek(lexer.TOKEN_RPAREN) {
		return
	}
	if !p._expectPeek(lexer.TOKEN_LBRACE) {
		return
	}
	p._nextToken() // Move to first token inside the then block

	// Parse nested triggers inside the then block
	for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
		if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "trigger" {
			nestedTrigger := p._parseTriggerWithContext(thenValue)
			if nestedTrigger != nil {
				trigger.Triggers = append(trigger.Triggers, *nestedTrigger)
				// Enforce that len == cap for trigger.Triggers
				if cap(trigger.Triggers) != len(trigger.Triggers) {
					tmp := make([]model.ActionTriggerDSLModel, len(trigger.Triggers))
					copy(tmp, trigger.Triggers)
					trigger.Triggers = tmp
				}
			}
		}
		p._nextToken()
	}
}

func (p *Parser) _parseFrame() *model.FrameDSLModel {
	frame := &model.FrameDSLModel{
		Type:      "FRAME",
		Variables: make([]model.VariableDSLModel, 0, 0),
		Blocks:    make([]model.BlockDSLModel, 0, 0),
	}

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	// Parse frame header
	frameAttrs := p._parseKeyValuePairs()
	frame.Name = frameAttrs["name"]
	frame.Route = frameAttrs["route"]

	// Check for unknown fields
	for key := range frameAttrs {
		if key != "name" && key != "route" {
			p.errors = append(p.errors, fmt.Sprintf("Unknown frame field: %s", key))
		}
	}

	if !p._curTokenIs(lexer.TOKEN_RPAREN) {
		p.errors = append(p.errors, "Expected ')' to close frame header")
		return nil
	}

	// Parse frame body
	if p._peekTokenIs(lexer.TOKEN_LBRACE) {
		p._nextToken()
		p._nextToken() // Move to first token inside the block
		for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
			if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "var" {
				varDecl := p._parseVariable()
				if varDecl != nil {
					frame.Variables = append(frame.Variables, *varDecl)
					// Enforce that len == cap for frame.Variables
					if cap(frame.Variables) != len(frame.Variables) {
						tmp := make([]model.VariableDSLModel, len(frame.Variables))
						copy(tmp, frame.Variables)
						frame.Variables = tmp
					}
				}
			} else if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "block" {
				block := p._parseBlock()
				if block != nil {
					frame.Blocks = append(frame.Blocks, *block)
					// Enforce that len == cap for frame.Blocks
					if cap(frame.Blocks) != len(frame.Blocks) {
						tmp := make([]model.BlockDSLModel, len(frame.Blocks))
						copy(tmp, frame.Blocks)
						frame.Blocks = tmp
					}
				}
			} else {
				// Unexpected token in frame body
				p.errors = append(p.errors, fmt.Sprintf("Line %d, Column %d: unexpected token '%s' in frame body, expected 'var' or 'block'",
					p.curToken.Line, p.curToken.Column, p.curToken.Literal))
			}
			p._nextToken()
		}
	}

	// Return nil if there were any parsing errors
	if len(p.errors) > 0 {
		return nil
	}

	return frame
}

func (p *Parser) _parseVariable() *model.VariableDSLModel {
	if !p._expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	key := p.curToken.Literal
	if !p._expectPeek(lexer.TOKEN_COLON) {
		return nil
	}
	if !p._expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	typ := p.curToken.Literal
	if !p._expectPeek(lexer.TOKEN_ASSIGN) {
		return nil
	}
	p._nextToken()

	value := p.curToken.Literal

	return &model.VariableDSLModel{
		Key:   key,
		Type:  typ,
		Value: value,
	}
}

func (p *Parser) _parseBlock() *model.BlockDSLModel {
	block := &model.BlockDSLModel{
		Data:       make([]model.BlockDataDSLModel, 0, 0),
		Properties: make([]model.BlockPropertyDSLModel, 0, 0),
		Slots:      make([]model.BlockSlotDSLModel, 0, 0),
		Blocks:     make([]model.BlockDSLModel, 0, 0),
		Actions:    make([]model.ActionDSLModel, 0, 0),
	}

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	// Parse block header
	blockAttrs := p._parseKeyValuePairs()
	block.KeyType = blockAttrs["keyType"]
	block.Key = blockAttrs["key"]
	block.VisibilityKey = blockAttrs["visibility"]
	if version, ok := blockAttrs["version"]; ok {
		block.IntegrationVersion, _ = strconv.Atoi(version)
	}

	// Parse block methods
	for p._peekTokenIs(lexer.TOKEN_DOT) {
		p._nextToken()
		if p._expectPeek(lexer.TOKEN_KEYWORD) {
			switch p.curToken.Literal {
			case "data":
				dataItems := p._parseBlockData()
				block.Data = append(block.Data, dataItems...)
				// Enforce that len == cap for block.Data
				if cap(block.Data) != len(block.Data) {
					tmp := make([]model.BlockDataDSLModel, len(block.Data))
					copy(tmp, block.Data)
					block.Data = tmp
				}
			case "prop":
				propItems := p._parseBlockProperty()
				block.Properties = append(block.Properties, propItems...)
				// Enforce that len == cap for block.Properties
				if cap(block.Properties) != len(block.Properties) {
					tmp := make([]model.BlockPropertyDSLModel, len(block.Properties))
					copy(tmp, block.Properties)
					block.Properties = tmp
				}
			case "slot":
				p._parseSlot(block)
			case "action":
				action := p._parseAction()
				block.Actions = append(block.Actions, action)
				// Enforce that len == cap for block.Actions
				if cap(block.Actions) != len(block.Actions) {
					tmp := make([]model.ActionDSLModel, len(block.Actions))
					copy(tmp, block.Actions)
					block.Actions = tmp
				}
			}
		}
	}
	return block
}

func (p *Parser) _parseAction() model.ActionDSLModel {
	action := model.ActionDSLModel{
		Triggers: make([]model.ActionTriggerDSLModel, 0, 0),
	}

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return action
	}

	// Parse action header
	actionAttrs := p._parseKeyValuePairs()
	action.Event = actionAttrs["event"]

	if !p._expectPeek(lexer.TOKEN_LBRACE) {
		return action
	}
	p._nextToken() // Move to first token inside the action block

	// Parse triggers inside action
	for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
		if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "trigger" {
			trigger := p._parseTriggerWithContext("NEXT")
			if trigger != nil {
				action.Triggers = append(action.Triggers, *trigger)
				// Enforce that len == cap for action.Triggers
				if cap(action.Triggers) != len(action.Triggers) {
					tmp := make([]model.ActionTriggerDSLModel, len(action.Triggers))
					copy(tmp, action.Triggers)
					action.Triggers = tmp
				}
			}
		}
		p._nextToken()
	}
	return action
}

func (p *Parser) _parseTriggerWithContext(defaultThen string) *model.ActionTriggerDSLModel {
	trigger := &model.ActionTriggerDSLModel{
		Properties: make([]model.TriggerPropertyDSLModel, 0, 0),
		Data:       make([]model.TriggerDataDSLModel, 0, 0),
		Triggers:   make([]model.ActionTriggerDSLModel, 0, 0),
		Then:       defaultThen,
	}

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	// Parse trigger header
	triggerAttrs := p._parseKeyValuePairs()
	trigger.KeyType = triggerAttrs["keyType"]
	trigger.Name = triggerAttrs["name"]
	if then, ok := triggerAttrs["then"]; ok {
		trigger.Then = then
	}
	if version, ok := triggerAttrs["version"]; ok {
		trigger.IntegrationVersion, _ = strconv.Atoi(version)
	}

	// Check for unknown fields
	for key := range triggerAttrs {
		if key != "keyType" && key != "name" && key != "then" && key != "version" {
			p.errors = append(p.errors, fmt.Sprintf("Unknown trigger field: %s", key))
		}
	}

	// Parse trigger methods
	for p._peekTokenIs(lexer.TOKEN_DOT) {
		p._nextToken()
		if p._expectPeek(lexer.TOKEN_KEYWORD) {
			switch p.curToken.Literal {
			case "data":
				dataItems := p._parseTriggerData()
				trigger.Data = append(trigger.Data, dataItems...)
				// Enforce that len == cap for trigger.Data
				if cap(trigger.Data) != len(trigger.Data) {
					tmp := make([]model.TriggerDataDSLModel, len(trigger.Data))
					copy(tmp, trigger.Data)
					trigger.Data = tmp
				}
			case "prop":
				propItems := p._parseTriggerProperty()
				trigger.Properties = append(trigger.Properties, propItems...)
				// Enforce that len == cap for trigger.Properties
				if cap(trigger.Properties) != len(trigger.Properties) {
					tmp := make([]model.TriggerPropertyDSLModel, len(trigger.Properties))
					copy(tmp, trigger.Properties)
					trigger.Properties = tmp
				}
			case "then":
				p._parseThen(trigger)
			}
		}
	}
	return trigger
}
