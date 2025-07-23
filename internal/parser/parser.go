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

func (p *Parser) _parseFrame() *model.FrameDSLModel {
	frame := &model.FrameDSLModel{
		Type:      "FRAME",
		Variables: []model.VariableDSLModel{},
		Blocks:    []model.BlockDSLModel{},
	}

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	for !p._curTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errors = append(p.errors, "Expected identifier in frame header")
			return nil
		}
		key := p.curToken.Literal
		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return nil
		}
		p._nextToken()
		switch key {
		case "name":
			frame.Name = p.curToken.Literal
		case "route":
			frame.Route = p.curToken.Literal
		default:
			p.errors = append(p.errors, fmt.Sprintf("Unknown frame field: %s", key))
		}
		p._nextToken()
		if p._curTokenIs(lexer.TOKEN_COMMA) {
			continue
		} else if p._curTokenIs(lexer.TOKEN_RPAREN) {
			break
		}
	}

	if !p._curTokenIs(lexer.TOKEN_RPAREN) {
		p.errors = append(p.errors, "Expected ')' to close frame header")
		return nil
	}

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
	var value string
	if p._curTokenIs(lexer.TOKEN_BOOLEAN) {
		value = p.curToken.Literal
	} else if p._curTokenIs(lexer.TOKEN_STRING) {
		value = p.curToken.Literal
	} else if p._curTokenIs(lexer.TOKEN_INT) || p._curTokenIs(lexer.TOKEN_LONG) || p._curTokenIs(lexer.TOKEN_FLOAT) || p._curTokenIs(lexer.TOKEN_DOUBLE) {
		value = p.curToken.Literal
	} else {
		value = p.curToken.Literal
	}

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
	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		key := p.curToken.Literal
		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return nil
		}
		p._nextToken()
		switch key {
		case "keyType":
			block.KeyType = p.curToken.Literal
		case "key":
			block.Key = p.curToken.Literal
		case "visibility":
			block.VisibilityKey = p.curToken.Literal
		case "version":
			block.IntegrationVersion, _ = strconv.Atoi(p.curToken.Literal)
		}
		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
		} else if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	for p._peekTokenIs(lexer.TOKEN_DOT) {
		p._nextToken()
		if p._expectPeek(lexer.TOKEN_KEYWORD) {
			switch p.curToken.Literal {
			case "data":
				block.Data = append(block.Data, p._parseBlockData())
				p._expectPeek(lexer.TOKEN_RPAREN)
			case "prop":
				block.Properties = append(block.Properties, p._parseBlockProperty())
			case "slot":
				slot := p._parseBlockSlot()
				block.Slots = append(block.Slots, slot)

				// Parse child blocks in the slot
				for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
					if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "block" {
						child := p._parseBlock()
						if child != nil {
							child.Slot = slot.Slot
							block.Blocks = append(block.Blocks, *child)
						}
					}
					p._nextToken()
				}
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

	p._expectPeek(lexer.TOKEN_LPAREN)
	p._expectPeek(lexer.TOKEN_IDENT)
	p._expectPeek(lexer.TOKEN_ASSIGN)
	p._nextToken()
	action.Event = p.curToken.Literal
	p._expectPeek(lexer.TOKEN_RPAREN)
	p._expectPeek(lexer.TOKEN_LBRACE)
	p._nextToken() // Move to first token inside the action block

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

func (p *Parser) _parseTrigger() *model.ActionTriggerDSLModel {
	return p._parseTriggerWithContext("")
}

func (p *Parser) _parseTriggerWithContext(defaultThen string) *model.ActionTriggerDSLModel {
	trigger := &model.ActionTriggerDSLModel{
		Properties: []model.TriggerPropertyDSLModel{},
		Data:       []model.TriggerDataDSLModel{},
		Triggers:   []model.ActionTriggerDSLModel{},
		Then:       defaultThen,
	}

	p._expectPeek(lexer.TOKEN_LPAREN)
	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		key := p.curToken.Literal
		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return nil
		}
		p._nextToken()
		switch key {
		case "keyType":
			trigger.KeyType = p.curToken.Literal
		case "name":
			trigger.Name = p.curToken.Literal
		case "then":
			trigger.Then = p.curToken.Literal
		case "version":
			trigger.IntegrationVersion, _ = strconv.Atoi(p.curToken.Literal)
		default:
			p.errors = append(p.errors, fmt.Sprintf("Unknown trigger field: %s", key))
		}
		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
		} else if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	for p._peekTokenIs(lexer.TOKEN_DOT) {
		p._nextToken()
		if p._expectPeek(lexer.TOKEN_KEYWORD) {
			switch p.curToken.Literal {
			case "data":
				trigger.Data = append(trigger.Data, p._parseTriggerData())
				p._expectPeek(lexer.TOKEN_RPAREN)

			case "prop":
				property := p._parseTriggerProperty()
				trigger.Properties = append(trigger.Properties, property)

			case "then":
				p._expectPeek(lexer.TOKEN_LPAREN)
				p._expectPeek(lexer.TOKEN_STRING)
				thenValue := p.curToken.Literal
				p._expectPeek(lexer.TOKEN_RPAREN)
				p._expectPeek(lexer.TOKEN_LBRACE)
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
		}
	}
	return trigger
}

func (p *Parser) _parseBlockProperty() model.BlockPropertyDSLModel {
	p._expectPeek(lexer.TOKEN_LPAREN)
	p._expectPeek(lexer.TOKEN_IDENT)
	key := p.curToken.Literal
	p._expectPeek(lexer.TOKEN_ASSIGN)
	p._expectPeek(lexer.TOKEN_LPAREN)

	// Parse property values for different device types
	var valueMobile, valueTablet, valueDesktop string
	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		propKey := p.curToken.Literal
		p._expectPeek(lexer.TOKEN_ASSIGN)
		p._nextToken()
		propValue := p.curToken.Literal

		switch propKey {
		case "mobile":
			valueMobile = propValue
		case "tablet":
			valueTablet = propValue
		case "desktop":
			valueDesktop = propValue
		case "value":
			valueMobile = propValue
			valueTablet = propValue
			valueDesktop = propValue
		}

		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
		} else {
			break
		}
	}

	p._expectPeek(lexer.TOKEN_RPAREN)
	p._expectPeek(lexer.TOKEN_RPAREN)
	return model.BlockPropertyDSLModel{
		Key:          key,
		ValueMobile:  valueMobile,
		ValueTablet:  valueTablet,
		ValueDesktop: valueDesktop,
		Type:         "PLACEHOLDER",
	}
}

func (p *Parser) _parseBlockData() model.BlockDataDSLModel {
	p._expectPeek(lexer.TOKEN_LPAREN)
	p._expectPeek(lexer.TOKEN_IDENT)
	key := p.curToken.Literal
	p._expectPeek(lexer.TOKEN_ASSIGN)
	p._nextToken()
	val := p.curToken.Literal
	return model.BlockDataDSLModel{
		Key:   key,
		Value: val,
		Type:  "PLACEHOLDER",
	}
}

func (p *Parser) _parseBlockSlot() model.BlockSlotDSLModel {
	p._expectPeek(lexer.TOKEN_LPAREN)
	p._expectPeek(lexer.TOKEN_STRING)
	slotName := p.curToken.Literal
	p._expectPeek(lexer.TOKEN_RPAREN)
	p._expectPeek(lexer.TOKEN_LBRACE)
	p._nextToken() // Move to first token inside the slot

	return model.BlockSlotDSLModel{Slot: slotName}
}

func (p *Parser) _parseTriggerProperty() model.TriggerPropertyDSLModel {
	p._expectPeek(lexer.TOKEN_LPAREN)
	p._expectPeek(lexer.TOKEN_IDENT)
	key := p.curToken.Literal
	p._expectPeek(lexer.TOKEN_ASSIGN)
	p._expectPeek(lexer.TOKEN_LPAREN)
	p._expectPeek(lexer.TOKEN_IDENT)
	p._expectPeek(lexer.TOKEN_ASSIGN)
	p._nextToken()
	val := p.curToken.Literal
	p._expectPeek(lexer.TOKEN_RPAREN)
	p._expectPeek(lexer.TOKEN_RPAREN)
	return model.TriggerPropertyDSLModel{
		Key:   key,
		Value: val,
		Type:  "PLACEHOLDER",
	}
}

func (p *Parser) _parseTriggerData() model.TriggerDataDSLModel {
	p._expectPeek(lexer.TOKEN_LPAREN)
	p._expectPeek(lexer.TOKEN_IDENT)
	key := p.curToken.Literal
	p._expectPeek(lexer.TOKEN_ASSIGN)
	p._nextToken()
	val := p.curToken.Literal
	return model.TriggerDataDSLModel{
		Key:   key,
		Value: val,
		Type:  "PLACEHOLDER",
	}
}
