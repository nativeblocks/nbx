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
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t lexer.TokenType) {
	err := fmt.Sprintf("Line %d, Column %d: expected next token to be %v, got %v",
		p.peekToken.Line, p.peekToken.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) ParseSDUI() *model.FrameDSLModel {
	if !p.curTokenIs(lexer.TOKEN_KEYWORD) || p.curToken.Literal != "frame" {
		p.errors = append(p.errors, "Program must start with a frame declaration")
		return nil
	}
	return p.parseFrame()
}

func (p *Parser) parseFrame() *model.FrameDSLModel {
	frame := &model.FrameDSLModel{
		Type:      "FRAME",
		Variables: []model.VariableDSLModel{},
		Blocks:    []model.BlockDSLModel{},
	}

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	for !p.curTokenIs(lexer.TOKEN_RPAREN) {
		p.nextToken()
		if !p.curTokenIs(lexer.TOKEN_IDENT) {
			p.errors = append(p.errors, "Expected identifier in frame header")
			return nil
		}
		key := p.curToken.Literal
		if !p.expectPeek(lexer.TOKEN_ASSIGN) {
			return nil
		}
		p.nextToken()
		switch key {
		case "name":
			frame.Name = p.curToken.Literal
		case "route":
			frame.Route = p.curToken.Literal
		default:
			p.errors = append(p.errors, fmt.Sprintf("Unknown frame field: %s", key))
		}
		p.nextToken()
		if p.curTokenIs(lexer.TOKEN_COMMA) {
			continue
		} else if p.curTokenIs(lexer.TOKEN_RPAREN) {
			break
		}
	}

	if !p.curTokenIs(lexer.TOKEN_RPAREN) {
		p.errors = append(p.errors, "Expected ')' to close frame header")
		return nil
	}

	if p.peekTokenIs(lexer.TOKEN_LBRACE) {
		p.nextToken()
		p.nextToken() // Move to first token inside the block
		for !p.curTokenIs(lexer.TOKEN_RBRACE) && !p.curTokenIs(lexer.TOKEN_EOF) {
			if p.curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "var" {
				varDecl := p.parseVariable()
				if varDecl != nil {
					frame.Variables = append(frame.Variables, *varDecl)
				}
			} else if p.curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "block" {
				block := p.parseBlock()
				if block != nil {
					frame.Blocks = append(frame.Blocks, *block)
				}
			}
			p.nextToken()
		}
	}

	return frame
}

func (p *Parser) parseVariable() *model.VariableDSLModel {
	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	key := p.curToken.Literal
	if !p.expectPeek(lexer.TOKEN_COLON) {
		return nil
	}
	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	typ := p.curToken.Literal
	if !p.expectPeek(lexer.TOKEN_ASSIGN) {
		return nil
	}
	p.nextToken()
	var value string
	if p.curTokenIs(lexer.TOKEN_BOOLEAN) {
		value = p.curToken.Literal
	} else if p.curTokenIs(lexer.TOKEN_STRING) {
		value = p.curToken.Literal
	} else if p.curTokenIs(lexer.TOKEN_INT) || p.curTokenIs(lexer.TOKEN_LONG) || p.curTokenIs(lexer.TOKEN_FLOAT) || p.curTokenIs(lexer.TOKEN_DOUBLE) {
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

func (p *Parser) parseBlock() *model.BlockDSLModel {
	block := &model.BlockDSLModel{
		Data:       []model.BlockDataDSLModel{},
		Properties: []model.BlockPropertyDSLModel{},
		Slots:      []model.BlockSlotDSLModel{},
		Blocks:     []model.BlockDSLModel{},
		Actions:    []model.ActionDSLModel{},
	}

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}
	for !p.curTokenIs(lexer.TOKEN_RPAREN) && !p.curTokenIs(lexer.TOKEN_EOF) {
		p.nextToken()
		key := p.curToken.Literal
		if !p.expectPeek(lexer.TOKEN_ASSIGN) {
			return nil
		}
		p.nextToken()
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
		if p.peekTokenIs(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.peekTokenIs(lexer.TOKEN_RPAREN) {
			p.nextToken()
			break
		}
	}

	for p.peekTokenIs(lexer.TOKEN_DOT) {
		p.nextToken()
		if p.expectPeek(lexer.TOKEN_KEYWORD) {
			switch p.curToken.Literal {
			case "data":
				p.parseBlockDataMethod(block)

			case "prop":
				p.parseBlockPropMethod(block)

			case "props":
				p.parseBlockPropsMethod(block)

			case "slot":
				p.parseSlot(block)

			case "action":
				block.Actions = append(block.Actions, p.parseAction())
			}
		}
	}
	return block
}

func (p *Parser) parseSlot(block *model.BlockDSLModel) {
	p.expectPeek(lexer.TOKEN_LPAREN)
	p.expectPeek(lexer.TOKEN_STRING)
	slotName := p.curToken.Literal
	p.expectPeek(lexer.TOKEN_RPAREN)
	p.expectPeek(lexer.TOKEN_LBRACE)
	p.nextToken() // Move to first token inside the slot

	// Create a slot
	slot := model.BlockSlotDSLModel{Slot: slotName}
	block.Slots = append(block.Slots, slot)

	// Parse child blocks in the slot
	for !p.curTokenIs(lexer.TOKEN_RBRACE) && !p.curTokenIs(lexer.TOKEN_EOF) {
		if p.curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "block" {
			child := p.parseBlock()
			if child != nil {
				child.Slot = slotName
				block.Blocks = append(block.Blocks, *child)
			}
		}
		p.nextToken()
	}
}

func (p *Parser) parseBlockDataMethod(block *model.BlockDSLModel) {
	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return
	}

	for !p.curTokenIs(lexer.TOKEN_RPAREN) && !p.curTokenIs(lexer.TOKEN_EOF) {
		p.nextToken()
		key := p.curToken.Literal
		if !p.expectPeek(lexer.TOKEN_ASSIGN) {
			return
		}
		p.nextToken()
		value := p.curToken.Literal

		block.Data = append(block.Data, model.BlockDataDSLModel{
			Key:   key,
			Value: value,
			Type:  "STRING",
		})

		if p.peekTokenIs(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.peekTokenIs(lexer.TOKEN_RPAREN) {
			p.nextToken()
			break
		}
	}
}

func (p *Parser) parseBlockPropsMethod(block *model.BlockDSLModel) {
	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return
	}

	for !p.curTokenIs(lexer.TOKEN_RPAREN) && !p.curTokenIs(lexer.TOKEN_EOF) {
		p.nextToken()
		key := p.curToken.Literal
		if !p.expectPeek(lexer.TOKEN_ASSIGN) {
			return
		}
		p.nextToken()
		value := p.curToken.Literal

		block.Properties = append(block.Properties, model.BlockPropertyDSLModel{
			Key:          key,
			ValueMobile:  value,
			ValueTablet:  value,
			ValueDesktop: value,
			Type:         "STRING",
		})

		if p.peekTokenIs(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.peekTokenIs(lexer.TOKEN_RPAREN) {
			p.nextToken()
			break
		}
	}
}

func (p *Parser) parseBlockPropMethod(block *model.BlockDSLModel) {
	p.expectPeek(lexer.TOKEN_LPAREN)
	p.expectPeek(lexer.TOKEN_IDENT)
	key := p.curToken.Literal
	p.expectPeek(lexer.TOKEN_ASSIGN)
	p.expectPeek(lexer.TOKEN_LPAREN)

	// Parse property values for different device types
	var valueMobile, valueTablet, valueDesktop string
	for !p.curTokenIs(lexer.TOKEN_RPAREN) && !p.curTokenIs(lexer.TOKEN_EOF) {
		p.nextToken()
		propKey := p.curToken.Literal
		p.expectPeek(lexer.TOKEN_ASSIGN)
		p.nextToken()
		propValue := p.curToken.Literal

		switch propKey {
		case "valueMobile":
			valueMobile = propValue
		case "mobile":
			valueMobile = propValue
		case "valueTablet":
			valueTablet = propValue
		case "tablet":
			valueTablet = propValue
		case "valueDesktop":
			valueDesktop = propValue
		case "desktop":
			valueDesktop = propValue
		case "value":
			valueMobile = propValue
			valueTablet = propValue
			valueDesktop = propValue
		}

		if p.peekTokenIs(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else {
			break
		}
	}

	p.expectPeek(lexer.TOKEN_RPAREN)
	p.expectPeek(lexer.TOKEN_RPAREN)

	block.Properties = append(block.Properties, model.BlockPropertyDSLModel{
		Key:          key,
		ValueMobile:  valueMobile,
		ValueTablet:  valueTablet,
		ValueDesktop: valueDesktop,
		Type:         "STRING",
	})
}

func (p *Parser) parseAction() model.ActionDSLModel {
	action := model.ActionDSLModel{
		Triggers: []model.ActionTriggerDSLModel{},
	}

	p.expectPeek(lexer.TOKEN_LPAREN)
	p.expectPeek(lexer.TOKEN_IDENT)
	p.expectPeek(lexer.TOKEN_ASSIGN)
	p.nextToken()
	action.Event = p.curToken.Literal
	p.expectPeek(lexer.TOKEN_RPAREN)
	p.expectPeek(lexer.TOKEN_LBRACE)
	p.nextToken() // Move to first token inside the action block

	for !p.curTokenIs(lexer.TOKEN_RBRACE) && !p.curTokenIs(lexer.TOKEN_EOF) {
		if p.curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "trigger" {
			trigger := p.parseTriggerWithContext("NEXT")
			if trigger != nil {
				action.Triggers = append(action.Triggers, *trigger)
			}
		}
		p.nextToken()
	}
	return action
}

func (p *Parser) parseTrigger() *model.ActionTriggerDSLModel {
	return p.parseTriggerWithContext("")
}

func (p *Parser) parseTriggerWithContext(defaultThen string) *model.ActionTriggerDSLModel {
	trigger := &model.ActionTriggerDSLModel{
		Properties: []model.TriggerPropertyDSLModel{},
		Data:       []model.TriggerDataDSLModel{},
		Triggers:   []model.ActionTriggerDSLModel{},
		Then:       defaultThen,
	}

	p.expectPeek(lexer.TOKEN_LPAREN)
	for !p.curTokenIs(lexer.TOKEN_RPAREN) && !p.curTokenIs(lexer.TOKEN_EOF) {
		p.nextToken()
		key := p.curToken.Literal
		if !p.expectPeek(lexer.TOKEN_ASSIGN) {
			return nil
		}
		p.nextToken()
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
		if p.peekTokenIs(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.peekTokenIs(lexer.TOKEN_RPAREN) {
			p.nextToken()
			break
		}
	}

	for p.peekTokenIs(lexer.TOKEN_DOT) {
		p.nextToken()
		if p.expectPeek(lexer.TOKEN_KEYWORD) {
			switch p.curToken.Literal {
			case "data":
				p.parseTriggerDataMethod(trigger)

			case "prop":

				p.parseTriggerPropMethod(trigger)

			case "props":
				p.parseTriggerPropsMethod(trigger)

			case "then":
				p.expectPeek(lexer.TOKEN_LPAREN)
				p.expectPeek(lexer.TOKEN_STRING)
				thenValue := p.curToken.Literal
				p.expectPeek(lexer.TOKEN_RPAREN)
				p.expectPeek(lexer.TOKEN_LBRACE)
				p.nextToken() // Move to first token inside the then block

				// Parse nested triggers inside the then block
				for !p.curTokenIs(lexer.TOKEN_RBRACE) && !p.curTokenIs(lexer.TOKEN_EOF) {
					if p.curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "trigger" {
						nestedTrigger := p.parseTriggerWithContext(thenValue)
						if nestedTrigger != nil {
							trigger.Triggers = append(trigger.Triggers, *nestedTrigger)
						}
					}
					p.nextToken()
				}
			}
		}
	}
	return trigger
}

func (p *Parser) parseTriggerDataMethod(trigger *model.ActionTriggerDSLModel) {
	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return
	}

	for !p.curTokenIs(lexer.TOKEN_RPAREN) && !p.curTokenIs(lexer.TOKEN_EOF) {
		p.nextToken()
		key := p.curToken.Literal
		if !p.expectPeek(lexer.TOKEN_ASSIGN) {
			return
		}
		p.nextToken()
		value := p.curToken.Literal

		trigger.Data = append(trigger.Data, model.TriggerDataDSLModel{
			Key:   key,
			Value: value,
			Type:  "STRING",
		})

		if p.peekTokenIs(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.peekTokenIs(lexer.TOKEN_RPAREN) {
			p.nextToken()
			break
		}
	}
}

func (p *Parser) parseTriggerPropsMethod(trigger *model.ActionTriggerDSLModel) {
	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return
	}

	for !p.curTokenIs(lexer.TOKEN_RPAREN) && !p.curTokenIs(lexer.TOKEN_EOF) {
		p.nextToken()
		key := p.curToken.Literal
		if !p.expectPeek(lexer.TOKEN_ASSIGN) {
			return
		}
		p.nextToken()
		value := p.curToken.Literal

		trigger.Properties = append(trigger.Properties, model.TriggerPropertyDSLModel{
			Key:   key,
			Value: value,
			Type:  "STRING",
		})

		if p.peekTokenIs(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.peekTokenIs(lexer.TOKEN_RPAREN) {
			p.nextToken()
			break
		}
	}
}

func (p *Parser) parseTriggerPropMethod(trigger *model.ActionTriggerDSLModel) {
	p.expectPeek(lexer.TOKEN_LPAREN)
	p.expectPeek(lexer.TOKEN_IDENT)
	key := p.curToken.Literal
	p.expectPeek(lexer.TOKEN_ASSIGN)
	p.expectPeek(lexer.TOKEN_LPAREN)
	p.expectPeek(lexer.TOKEN_IDENT)
	p.expectPeek(lexer.TOKEN_ASSIGN)
	p.nextToken()
	val := p.curToken.Literal
	p.expectPeek(lexer.TOKEN_RPAREN)
	p.expectPeek(lexer.TOKEN_RPAREN)
	trigger.Properties = append(trigger.Properties, model.TriggerPropertyDSLModel{Key: key, Value: val, Type: "STRING"})
}
