package parser

import (
	"fmt"
	"strconv"

	"github.com/nativeblocks/nbx/internal/errors"
	"github.com/nativeblocks/nbx/internal/lexer"
	"github.com/nativeblocks/nbx/internal/model"
	"github.com/nativeblocks/nbx/internal/types"
)

type Parser struct {
	l              *lexer.Lexer
	curToken       lexer.Token // current token being processed
	peekToken      lexer.Token // next token (for lookahead)
	errorCollector *errors.ErrorCollector
}

func NewParser(l *lexer.Lexer, source string) *Parser {
	p := &Parser{
		l:              l,
		errorCollector: errors.NewErrorCollector(source),
	}
	// read two tokens to initialize curToken and peekToken
	p._nextToken()
	p._nextToken()
	return p
}

func (p *Parser) ErrorCollector() *errors.ErrorCollector {
	return p.errorCollector
}

func (p *Parser) ParseNBX() *model.FrameDSLModel {
	if !p._curTokenIs(lexer.TOKEN_KEYWORD) || p.curToken.Literal != "frame" {
		p.errorCollector.AddTokenError(
			"Program must start with a frame declaration",
			p.curToken,
			"Add 'frame(name=\"...\", route=\"...\") { ... }' at the beginning",
		)
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
	expected := lexer.Token{Type: t}
	p.errorCollector.AddError(errors.UnexpectedTokenError(expected, p.peekToken))
}

// _curTokenIs checks if the current token is of the given type
func (p *Parser) _curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

// _peekTokenIs checks if the next token is of the given type
func (p *Parser) _peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func _enforceSliceCap[T any](slice []T) []T {
	if cap(slice) != len(slice) {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}
	return slice
}

func (p *Parser) _parseKeyValuePairs() map[string]string {
	pairs := make(map[string]string)

	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errorCollector.AddTokenError(
				"Expected identifier in attribute list",
				p.curToken,
				"Use format: key=\"value\"",
			)
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

func (p *Parser) _parseBlockData() []model.BlockDataDSLModel {
	dataList := make([]model.BlockDataDSLModel, 0)

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return dataList
	}

	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return dataList
	}

	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		keyLine, keyColumn := p.curToken.Line, p.curToken.Column

		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errorCollector.AddTokenError(
				"Expected identifier in data declaration",
				p.curToken,
				"Use format: data(key=value, ...)",
			)
			return dataList
		}
		key := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return dataList
		}
		p._nextToken()
		value := p.curToken.Literal

		inferredType := p._inferTypeFromToken(p.curToken)

		dataList = append(dataList, model.BlockDataDSLModel{
			Key:    key,
			Value:  value,
			Type:   inferredType.Name(),
			Line:   keyLine,
			Column: keyColumn,
		})

		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
		} else if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	return _enforceSliceCap(dataList)
}

func (p *Parser) _parseTriggerData() []model.TriggerDataDSLModel {
	dataList := make([]model.TriggerDataDSLModel, 0)

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return dataList
	}

	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return dataList
	}

	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()
		keyLine, keyColumn := p.curToken.Line, p.curToken.Column

		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errorCollector.AddTokenError(
				"Expected identifier in data declaration",
				p.curToken,
				"Use format: data(key=value, ...)",
			)
			return dataList
		}
		key := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return dataList
		}
		p._nextToken()
		value := p.curToken.Literal

		inferredType := p._inferTypeFromToken(p.curToken)
		dataList = append(dataList, model.TriggerDataDSLModel{
			Key:    key,
			Value:  value,
			Type:   inferredType.Name(),
			Line:   keyLine,
			Column: keyColumn,
		})

		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
		} else if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	return _enforceSliceCap(dataList)
}

func (p *Parser) _parseBlockProperty() []model.BlockPropertyDSLModel {
	propList := make([]model.BlockPropertyDSLModel, 0)

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return propList
	}

	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return propList
	}

	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()

		// handle trailing comma, if we hit closing paren after a comma, break
		if p._curTokenIs(lexer.TOKEN_RPAREN) {
			break
		}

		propLine, propColumn := p.curToken.Line, p.curToken.Column

		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errorCollector.AddTokenError(
				"Expected identifier in prop declaration",
				p.curToken,
				"Use format: prop(key=value, ...)",
			)
			return propList
		}
		propKey := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return propList
		}

		p._nextToken()
		if p.curToken.Type == lexer.TOKEN_LPAREN {
			deviceValues := make(map[string]string)
			for {
				p._nextToken()

				if p.curToken.Type == lexer.TOKEN_RPAREN {
					break
				}

				if p.curToken.Type != lexer.TOKEN_IDENT {
					p.errorCollector.AddTokenError(
						"Expected identifier inside property value parenthesis",
						p.curToken,
						"Use format: prop(key=(device=value, ...))",
					)
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
					continue
				}
			}

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

			inferValue := mobile
			if inferValue == "" {
				inferValue = tablet
			}
			if inferValue == "" {
				inferValue = desktop
			}
			inferredType := types.InferType(inferValue)

			propList = append(propList, model.BlockPropertyDSLModel{
				Key:          propKey,
				ValueMobile:  mobile,
				ValueTablet:  tablet,
				ValueDesktop: desktop,
				Type:         inferredType.Name(),
				Line:         propLine,
				Column:       propColumn,
			})
		} else {
			// single value, possibly multiline, parse everything up to next ',' or ')'
			value := p.curToken.Literal
			for {
				// if at end or comma (i.e. end of prop), break.
				if p._peekTokenIs(lexer.TOKEN_COMMA) || p._peekTokenIs(lexer.TOKEN_RPAREN) || p._peekTokenIs(lexer.TOKEN_EOF) {
					break
				}
				// otherwise, accumulate new tokens as part of the value (including newlines, etc.)
				p._nextToken()
				// Separate by space or newline character for clarity. If token type was NEWLINE (optional), add \n, else add a space.
				// But for now, we just add Literal with a newline.
				// Go lexers rarely make a newline token so we just use the literal.
				value += "\n" + p.curToken.Literal
			}
			inferredType := types.InferType(value)
			propList = append(propList, model.BlockPropertyDSLModel{
				Key:          propKey,
				ValueMobile:  value,
				ValueTablet:  value,
				ValueDesktop: value,
				Type:         inferredType.Name(),
				Line:         propLine,
				Column:       propColumn,
			})
		}

		// prop separator
		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
			continue
		}
		if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	return _enforceSliceCap(propList)
}

func (p *Parser) _parseTriggerProperty() []model.TriggerPropertyDSLModel {
	propList := make([]model.TriggerPropertyDSLModel, 0)

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return propList
	}

	if p._peekTokenIs(lexer.TOKEN_RPAREN) {
		p._nextToken()
		return propList
	}

	for !p._curTokenIs(lexer.TOKEN_RPAREN) && !p._curTokenIs(lexer.TOKEN_EOF) {
		p._nextToken()

		// handle trailing comma, if we hit closing paren after a comma, break
		if p._curTokenIs(lexer.TOKEN_RPAREN) {
			break
		}

		propLine, propColumn := p.curToken.Line, p.curToken.Column

		// Property key
		if !p._curTokenIs(lexer.TOKEN_IDENT) {
			p.errorCollector.AddTokenError(
				"Expected identifier in prop declaration",
				p.curToken,
				"Use format: prop(key=value, ...)",
			)
			return propList
		}
		propKey := p.curToken.Literal

		if !p._expectPeek(lexer.TOKEN_ASSIGN) {
			return propList
		}

		// only single value is supported
		p._nextToken()
		value := p.curToken.Literal
		for {
			// if at end or comma (i.e. end of prop), break.
			if p._peekTokenIs(lexer.TOKEN_COMMA) || p._peekTokenIs(lexer.TOKEN_RPAREN) || p._peekTokenIs(lexer.TOKEN_EOF) {
				break
			}
			// otherwise, accumulate new tokens as part of the value (including newlines, etc.)
			p._nextToken()
			value += "\n" + p.curToken.Literal
		}

		inferredType := types.InferType(value)
		propList = append(propList, model.TriggerPropertyDSLModel{
			Key:    propKey,
			Value:  value,
			Type:   inferredType.Name(),
			Line:   propLine,
			Column: propColumn,
		})

		// prop separator
		if p._peekTokenIs(lexer.TOKEN_COMMA) {
			p._nextToken()
			continue
		}
		if p._peekTokenIs(lexer.TOKEN_RPAREN) {
			p._nextToken()
			break
		}
	}

	return _enforceSliceCap(propList)
}

func (p *Parser) _parseSlot(block *model.BlockDSLModel) {
	slotLine, slotColumn := p.curToken.Line, p.curToken.Column

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

	slot := model.BlockSlotDSLModel{
		Slot:   slotName,
		Line:   slotLine,
		Column: slotColumn,
	}
	block.Slots = append(block.Slots, slot)
	block.Slots = _enforceSliceCap(block.Slots)

	for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
		if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "block" {
			child := p._parseBlock()
			if child != nil {
				child.Slot = slotName
				block.Blocks = append(block.Blocks, *child)
				block.Blocks = _enforceSliceCap(block.Blocks)
			}
		}
		p._nextToken()
	}
}

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
	p._nextToken() // move to first token inside the then block

	for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
		if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "trigger" {
			nestedTrigger := p._parseTrigger(thenValue)
			if nestedTrigger != nil {
				trigger.Triggers = append(trigger.Triggers, *nestedTrigger)
				trigger.Triggers = _enforceSliceCap(trigger.Triggers)
			}
		}
		p._nextToken()
	}
}

func (p *Parser) _parseFrame() *model.FrameDSLModel {
	frameLine, frameColumn := p.curToken.Line, p.curToken.Column

	frame := &model.FrameDSLModel{
		Type:      "FRAME",
		Variables: make([]model.VariableDSLModel, 0),
		Blocks:    make([]model.BlockDSLModel, 0),
		Line:      frameLine,
		Column:    frameColumn,
	}

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	frameAttrs := p._parseKeyValuePairs()
	frame.Name = frameAttrs["name"]
	frame.Route = frameAttrs["route"]

	for key := range frameAttrs {
		if key != "name" && key != "route" {
			validAttrs := []string{"name", "route"}
			p.errorCollector.AddError(errors.UnknownAttributeError(
				key, "frame", frame.Line, frame.Column, validAttrs,
			))
		}
	}

	if !p._curTokenIs(lexer.TOKEN_RPAREN) {
		p.errorCollector.AddTokenError(
			"Expected ')' to close frame header",
			p.curToken,
			"Add ')' after frame attributes",
		)
		return nil
	}

	if p._peekTokenIs(lexer.TOKEN_LBRACE) {
		p._nextToken()
		p._nextToken() // move to first token inside the block
		for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
			if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "var" {
				varDecl := p._parseVariable()
				if varDecl != nil {
					frame.Variables = append(frame.Variables, *varDecl)
					frame.Variables = _enforceSliceCap(frame.Variables)
				}
			} else if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "block" {
				block := p._parseBlock()
				if block != nil {
					frame.Blocks = append(frame.Blocks, *block)
					frame.Blocks = _enforceSliceCap(frame.Blocks)
				}
			} else {
				p.errorCollector.AddTokenError(
					fmt.Sprintf("Unexpected token '%s' in frame body", p.curToken.Literal),
					p.curToken,
					"Expected 'var' or 'block' declaration",
				)
			}
			p._nextToken()
		}
	}

	if p.errorCollector.HasErrors() {
		return nil
	}

	return frame
}

func (p *Parser) _parseVariable() *model.VariableDSLModel {
	varLine, varColumn := p.curToken.Line, p.curToken.Column

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
		Key:    key,
		Type:   typ,
		Value:  value,
		Line:   varLine,
		Column: varColumn,
	}
}

func (p *Parser) _parseBlock() *model.BlockDSLModel {
	blockLine, blockColumn := p.curToken.Line, p.curToken.Column

	block := &model.BlockDSLModel{
		Data:       make([]model.BlockDataDSLModel, 0),
		Properties: make([]model.BlockPropertyDSLModel, 0),
		Slots:      make([]model.BlockSlotDSLModel, 0),
		Blocks:     make([]model.BlockDSLModel, 0),
		Actions:    make([]model.ActionDSLModel, 0),
		Line:       blockLine,
		Column:     blockColumn,
	}

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	blockAttrs := p._parseKeyValuePairs()
	block.KeyType = blockAttrs["keyType"]
	block.Key = blockAttrs["key"]
	block.VisibilityKey = blockAttrs["visibility"]
	if version, ok := blockAttrs["version"]; ok {
		block.IntegrationVersion, _ = strconv.Atoi(version)
	}

	for p._peekTokenIs(lexer.TOKEN_DOT) {
		p._nextToken()
		if p._expectPeek(lexer.TOKEN_KEYWORD) {
			switch p.curToken.Literal {
			case "data":
				dataItems := p._parseBlockData()
				block.Data = append(block.Data, dataItems...)
				block.Data = _enforceSliceCap(block.Data)
			case "prop":
				propItems := p._parseBlockProperty()
				block.Properties = append(block.Properties, propItems...)
				block.Properties = _enforceSliceCap(block.Properties)
			case "slot":
				p._parseSlot(block)
			case "action":
				action := p._parseAction()
				block.Actions = append(block.Actions, action)
				block.Actions = _enforceSliceCap(block.Actions)
			}
		}
	}
	return block
}

func (p *Parser) _parseAction() model.ActionDSLModel {
	actionLine, actionColumn := p.curToken.Line, p.curToken.Column

	action := model.ActionDSLModel{
		Triggers: make([]model.ActionTriggerDSLModel, 0),
		Line:     actionLine,
		Column:   actionColumn,
	}

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return action
	}

	actionAttrs := p._parseKeyValuePairs()
	action.Event = actionAttrs["event"]

	if !p._expectPeek(lexer.TOKEN_LBRACE) {
		return action
	}
	p._nextToken() // move to first token inside the action block

	for !p._curTokenIs(lexer.TOKEN_RBRACE) && !p._curTokenIs(lexer.TOKEN_EOF) {
		if p._curTokenIs(lexer.TOKEN_KEYWORD) && p.curToken.Literal == "trigger" {
			trigger := p._parseTrigger("NEXT")
			if trigger != nil {
				action.Triggers = append(action.Triggers, *trigger)
				action.Triggers = _enforceSliceCap(action.Triggers)
			}
		}
		p._nextToken()
	}
	return action
}

func (p *Parser) _parseTrigger(defaultThen string) *model.ActionTriggerDSLModel {
	triggerLine, triggerColumn := p.curToken.Line, p.curToken.Column

	trigger := &model.ActionTriggerDSLModel{
		Properties: make([]model.TriggerPropertyDSLModel, 0),
		Data:       make([]model.TriggerDataDSLModel, 0),
		Triggers:   make([]model.ActionTriggerDSLModel, 0),
		Then:       defaultThen,
		Line:       triggerLine,
		Column:     triggerColumn,
	}

	if !p._expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	triggerAttrs := p._parseKeyValuePairs()
	trigger.KeyType = triggerAttrs["keyType"]
	trigger.Name = triggerAttrs["name"]
	if then, ok := triggerAttrs["then"]; ok {
		trigger.Then = then
	}
	if version, ok := triggerAttrs["version"]; ok {
		trigger.IntegrationVersion, _ = strconv.Atoi(version)
	}

	for key := range triggerAttrs {
		if key != "keyType" && key != "name" && key != "then" && key != "version" {
			validAttrs := []string{"keyType", "name", "then", "version"}
			p.errorCollector.AddError(errors.UnknownAttributeError(
				key, "trigger", trigger.Line, trigger.Column, validAttrs,
			))
		}
	}

	for p._peekTokenIs(lexer.TOKEN_DOT) {
		p._nextToken()
		if p._expectPeek(lexer.TOKEN_KEYWORD) {
			switch p.curToken.Literal {
			case "data":
				dataItems := p._parseTriggerData()
				trigger.Data = append(trigger.Data, dataItems...)
				trigger.Data = _enforceSliceCap(trigger.Data)
			case "prop":
				propItems := p._parseTriggerProperty()
				trigger.Properties = append(trigger.Properties, propItems...)
				trigger.Properties = _enforceSliceCap(trigger.Properties)
			case "then":
				p._parseThen(trigger)
			}
		}
	}
	return trigger
}

func (p *Parser) _inferTypeFromToken(token lexer.Token) types.Type {
	switch token.Type {
	case lexer.TOKEN_STRING:
		return types.TypeString
	case lexer.TOKEN_BOOLEAN:
		return types.TypeBoolean
	case lexer.TOKEN_INT:
		return types.TypeInt
	case lexer.TOKEN_LONG:
		return types.TypeLong
	case lexer.TOKEN_FLOAT:
		return types.TypeFloat
	case lexer.TOKEN_DOUBLE:
		return types.TypeDouble
	case lexer.TOKEN_IDENT:
		return types.InferType(token.Literal)
	default:
		return types.InferType(token.Literal)
	}
}
