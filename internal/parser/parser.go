package parser

import (
	"fmt"
	"github.com/nativeblocks/nbx/internal/lexer"
	"github.com/nativeblocks/nbx/internal/model"
	"strconv"
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
		errors: []string{},
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
	var dataList []model.BlockDataDSLModel

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

	return dataList
}

// _parseTriggerData parses a trigger .data() declaration and returns a list of data items
func (p *Parser) _parseTriggerData() []model.TriggerDataDSLModel {
	var dataList []model.TriggerDataDSLModel

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

	return dataList
}

// _parseBlockProperty parses a block .prop() declaration and returns a list of properties
func (p *Parser) _parseBlockProperty() []model.BlockPropertyDSLModel {
	var propList []model.BlockPropertyDSLModel

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return propList
	}

	// Handle empty prop()
	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return propList
	}

	// Parse all property declarations
	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errors = append(p.errors, "Expected identifier in prop declaration")
			return propList
		}
		key := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return propList
		}

		if !p._expectPeek(lexer.TOKEN_LPAREN) {
			return propList
		}

		// Parse property values
		var valueMobile, valueTablet, valueDesktop string
		propValues := p._parseKeyValuePairs()

		// Handle different property value formats
		if val, ok := propValues["value"]; ok {
			valueMobile = val
			valueTablet = val
			valueDesktop = val
		}
		if val, ok := propValues["valueMobile"]; ok {
			valueMobile = val
		}
		if val, ok := propValues["mobile"]; ok {
			valueMobile = val
		}
		if val, ok := propValues["valueTablet"]; ok {
			valueTablet = val
		}
		if val, ok := propValues["tablet"]; ok {
			valueTablet = val
		}
		if val, ok := propValues["valueDesktop"]; ok {
			valueDesktop = val
		}
		if val, ok := propValues["desktop"]; ok {
			valueDesktop = val
		}

		propList = append(propList, model.BlockPropertyDSLModel{
			Key:          key,
			ValueMobile:  valueMobile,
			ValueTablet:  valueTablet,
			ValueDesktop: valueDesktop,
			Type:         "PLACEHOLDER",
		})

		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
		} else if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	return propList
}

// _parseTriggerProperty parses a trigger .prop() declaration and returns a list of properties
func (p *Parser) _parseTriggerProperty() []model.TriggerPropertyDSLModel {
	var propList []model.TriggerPropertyDSLModel

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return propList
	}

	// Handle empty prop()
	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return propList
	}

	// Parse all property declarations
	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errors = append(p.errors, "Expected identifier in prop declaration")
			return propList
		}
		key := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return propList
		}

		if !p._expectPeek(lexer.TOKEN_LPAREN) {
			return propList
		}

		// Parse the nested value
		if !p._expectPeek(lexer.TOKEN_IDENT) {
			return propList
		}
		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return propList
		}
		p._nextToken()
		value := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_RPAREN) {
			return propList
		}

		propList = append(propList, model.TriggerPropertyDSLModel{
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

	// Parse child blocks in the slot
	for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
		if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "block" {
			child := p._parseBlock()
			if child != nil {
				child.Slot = slotName
				block.Blocks = append(block.Blocks, *child)
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
			}
		}
		p._nextToken()
	}
}

func (p *Parser) _parseFrame() *model.FrameDSLModel {
	frame := &model.FrameDSLModel{
		Type:      "FRAME",
		Variables: []model.VariableDSLModel{},
		Blocks:    []model.BlockDSLModel{},
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
				}
			} else if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "block" {
				block := p._parseBlock()
				if block != nil {
					frame.Blocks = append(frame.Blocks, *block)
				}
			}
			p._nextToken()
		}
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
		Data:       []model.BlockDataDSLModel{},
		Properties: []model.BlockPropertyDSLModel{},
		Slots:      []model.BlockSlotDSLModel{},
		Blocks:     []model.BlockDSLModel{},
		Actions:    []model.ActionDSLModel{},
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
			case "prop":
				propItems := p._parseBlockProperty()
				block.Properties = append(block.Properties, propItems...)
			case "slot":
				p._parseSlot(block)
			case "action":
				block.Actions = append(block.Actions, p._parseAction())
			}
		}
	}
	return block
}

func (p *Parser) _parseAction() model.ActionDSLModel {
	action := model.ActionDSLModel{
		Triggers: []model.ActionTriggerDSLModel{},
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
			}
		}
		p._nextToken()
	}
	return action
}

func (p *Parser) _parseTriggerWithContext(defaultThen string) *model.ActionTriggerDSLModel {
	trigger := &model.ActionTriggerDSLModel{
		Properties: []model.TriggerPropertyDSLModel{},
		Data:       []model.TriggerDataDSLModel{},
		Triggers:   []model.ActionTriggerDSLModel{},
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
			case "prop":
				propItems := p._parseTriggerProperty()
				trigger.Properties = append(trigger.Properties, propItems...)
			case "then":
				p._parseThen(trigger)
			}
		}
	}
	return trigger
}
