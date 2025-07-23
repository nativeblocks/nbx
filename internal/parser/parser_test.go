package parser

import (
	"github.com/nativeblocks/nbx/internal/lexer"
	"testing"
)

func TestParser_FrameOnly(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login"
)`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatalf("Expected frame to be parsed, got nil: %v", p.Errors())
	}

	if frame.Name != "login" {
		t.Errorf("Expected name to be 'login', got %s", frame.Name)
	}

	if frame.Route != "/login" {
		t.Errorf("Expected route to be '/login', got %s", frame.Route)
	}

}

func TestParser_FrameWithVariables(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login"
) {
	var visible: BOOLEAN = true
	var username: STRING = ""
	var password: STRING = ""
}
`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseNBX()

	if len(frame.Variables) != 3 {
		t.Fatalf("Expected 3 variables, got %d", len(frame.Variables))
	}

	if frame.Variables[0].Key != "visible" || frame.Variables[0].Type != "BOOLEAN" || frame.Variables[0].Value != "true" {
		t.Errorf("Unexpected variable 1: %+v", frame.Variables[0])
	}

	if frame.Variables[1].Key != "username" || frame.Variables[1].Type != "STRING" || frame.Variables[1].Value != "" {
		t.Errorf("Unexpected variable 2: %+v", frame.Variables[1])
	}

	if frame.Variables[2].Key != "password" || frame.Variables[2].Type != "STRING" || frame.Variables[2].Value != "" {
		t.Errorf("Unexpected variable 3: %+v", frame.Variables[2])
	}
}

func TestParser_FrameWithBlock(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login"
) {
	var visible: BOOLEAN = true
	var username: STRING = ""
	var password: STRING = ""
	block(keyType = "ROOT", key = "root", visibility = visible)
}
`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseNBX()

	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(frame.Blocks))
	}

	block := frame.Blocks[0]
	if block.KeyType != "ROOT" {
		t.Errorf("Expected block type ROOT, got %s", block.KeyType)
	}

	if block.Key != "root" {
		t.Errorf("Expected block key 'root', got %s", block.Key)
	}

	if block.VisibilityKey != "visible" {
		t.Errorf("Expected visibility 'visible', got %s", block.VisibilityKey)
	}
}

func TestParser_ComplexFrameWithBlocksAndActions(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login"
) {
    var visible: BOOLEAN = true
    var username: STRING = ""
    var password: STRING = ""

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "COLUMN", key = "main", visibility = visible)
        .slot("content") {
            block(keyType = "INPUT", key = "username", visibility = visible)
            .data(text = username)
            block(keyType = "INPUT", key = "password", visibility = visible)
            .prop(color = (valueMobile = "NONE", valueTablet = "NONE", valueDesktop = "NONE"))
            .data(text = password)
            .action(event = "onTextChange") {
                trigger(keyType = "VALIDATE", name = "validate password")
                .then("FAILURE") {
                    trigger(keyType = "SHOW_ERROR", name = "show error 1")
                    trigger(keyType = "CHANGE_COLOR", name = "change color to red")
                    .prop(color = "RED")
                }
                .then("SUCCESS") {
                    trigger(keyType = "SHOW_OK", name = "show ok")
                    .prop(color = "GREEN")
                }
            }
        }    
    }   
}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatalf("Expected frame to be parsed, got nil: %v", p.Errors())
	}

	if frame.Name != "login" {
		t.Errorf("Expected name to be 'login', got %s", frame.Name)
	}
	if frame.Route != "/login" {
		t.Errorf("Expected route to be '/login', got %s", frame.Route)
	}
	if len(frame.Variables) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(frame.Variables))
	}
	if len(frame.Blocks) == 0 {
		t.Fatalf("Expected at least 1 block, got 0")
	}
}

func TestParser_ComplexFrameWithBlocksAndActions2(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login"
) {
    var visible: BOOLEAN = true
    var username: STRING = ""
    var password: STRING = ""

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "nativeblocks/COLUMN", key = "main", visibility = visible)
        .slot("content") {
            block(keyType = "nativeblocks/TEXT_FIELD", key = "username", visibility = visible)
            .data(text = username)
            block(keyType = "nativeblocks/TEXT_FIELD", key = "password", visibility = visible)
            .prop(textColor = (valueMobile = "NONE", valueTablet = "NONE", valueDesktop = "NONE"))
            .data(text = password)
            .action(event = "onTextChange") {
                trigger(keyType = "nativeblocks/CHANGE_BLOCK_PROPERTY", name = "show error")
                .prop(propertyKey = "textColor")
                .prop(propertyValueDesktop = "RED")
				trigger(keyType = "nativeblocks/CHANGE_BLOCK_PROPERTY", name = "show success")
				.prop(propertyKey = "textColor")
				.prop(propertyValueDesktop = "GREEN")
            }
        }
    }
}
`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatalf("Expected frame to be parsed, got nil: %v", p.Errors())
	}

	if frame.Name != "login" {
		t.Errorf("Expected name to be 'login', got %s", frame.Name)
	}
	if frame.Route != "/login" {
		t.Errorf("Expected route to be '/login', got %s", frame.Route)
	}
	if len(frame.Variables) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(frame.Variables))
	}
	if len(frame.Blocks) == 0 {
		t.Fatalf("Expected at least 1 block, got 0")
	}
}

func TestParser_BlockWithPropertiesAndSlots(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login"
) {
    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "COLUMN", key = "main", visibility = visible)
        .slot("content") {
            block(keyType = "INPUT", key = "username", visibility = visible)
            .data(text = username)
        }
    }
}`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatalf("Expected frame to be parsed, got nil: %v", p.Errors())
	}
	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(frame.Blocks))
	}
	rootBlock := frame.Blocks[0]
	if rootBlock.KeyType != "ROOT" {
		t.Errorf("Expected root block type ROOT, got %s", rootBlock.KeyType)
	}
	if len(rootBlock.Blocks) != 1 {
		t.Fatalf("Expected 1 nested block in slot, got %d", len(rootBlock.Blocks))
	}
	columnBlock := rootBlock.Blocks[0]
	if columnBlock.KeyType != "COLUMN" {
		t.Errorf("Expected column block type COLUMN, got %s", columnBlock.KeyType)
	}
}

func TestParser_BlockWithActionAndData(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login"
) {
    block(keyType = "INPUT", key = "password", visibility = visible)
    .data(text = password)
    .action(event = "onTextChange") {
        trigger(keyType = "VALIDATE", name = "validate password")
        .then("SUCCESS") {
            trigger(keyType = "SHOW_OK", name = "show ok")
            .prop(color = "GREEN")
        }
    }
}`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatalf("Expected frame to be parsed, got nil: %v", p.Errors())
	}
	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(frame.Blocks))
	}
	block := frame.Blocks[0]
	if len(block.Actions) != 1 {
		t.Fatalf("Expected 1 action, got %d", len(block.Actions))
	}
	action := block.Actions[0]
	if action.Event != "onTextChange" {
		t.Errorf("Expected action event 'onTextChange', got %s", action.Event)
	}
	if len(action.Triggers) != 1 {
		t.Fatalf("Expected 1 trigger, got %d", len(action.Triggers))
	}
	trigger := action.Triggers[0]
	if trigger.Name != "validate password" {
		t.Errorf("Expected trigger name 'validate password', got %s", trigger.Name)
	}
	if len(trigger.Triggers) == 0 {
		t.Errorf("Expected trigger to have then blocks")
	}
}

func TestParser_VariableTypes(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login"
) {
    var visible: BOOLEAN = true
    var username: STRING = ""
    var password: STRING = ""
}`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatalf("Expected frame to be parsed, got nil: %v", p.Errors())
	}
	if len(frame.Variables) != 3 {
		t.Fatalf("Expected 3 variables, got %d", len(frame.Variables))
	}
	if frame.Variables[0].Key != "visible" || frame.Variables[0].Type != "BOOLEAN" || frame.Variables[0].Value != "true" {
		t.Errorf("Unexpected variable: %+v", frame.Variables[0])
	}
	if frame.Variables[1].Key != "username" || frame.Variables[1].Type != "STRING" || frame.Variables[1].Value != "" {
		t.Errorf("Unexpected variable: %+v", frame.Variables[1])
	}
	if frame.Variables[2].Key != "password" || frame.Variables[2].Type != "STRING" || frame.Variables[2].Value != "" {
		t.Errorf("Unexpected variable: %+v", frame.Variables[2])
	}
}

func TestParser_ComplexFrame(t *testing.T) {
	input := `
frame(
    name = "welcome",
    route = "/welcome"
) {
    var visible: BOOLEAN = true
    var enable: BOOLEAN = true
    var count: INT = 0
    var increaseButton: STRING = "+"
    var decreaseButton: STRING = "-"
    var welcome: STRING = "Welcome to Nativeblocks"
    var logo: STRING = "https://nativeblocks.io/nativeblocks_logo.png"

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "nativeblocks/column", key = "mainColumn", visibility = visible, version = 1)
        .prop(horizontalAlignment = "centerHorizontally", width = "match", height = "match")
        .slot("content") {

            block(keyType = "nativeblocks/column", key = "nativeblocksColumn", VISIBILITY = visible, version = 1)
            .prop(horizontalAlignment = "centerHorizontally", paddingTop = "64", weight = "0.4f", verticalArrangement = "spaceAround")
            .slot("content") {

                block(keyType = "nativeblocks/image", key = "logo", visibility = visible, version = 1)
                .prop(scaleType = "inside", width = "128", height = "128")
                .data(imageUrl = logo)

                block(keyType = "nativeblocks/text", key = "welcome", visibility = visible, version = 1)
                .prop(fontSize = "24", textAlign = "center", width = "wrap")
                .data(text = welcome)
            }

            block(keyType = "nativeblocks/row", key = "buttonsRow", visibility = visible, version = 1)
            .prop(horizontalArrangement = "spaceAround", verticalAlignment = "centerVertically", paddingTop = "12", weight = "0.6f")
            .slot("content") {

                block(keyType = "nativeblocks/button", key = "decreaseButton", visibility = visible, version = 1)
                .prop(backgroundColor = "#2563EB", borderColor = "#2563EB", radiusTopStart = "32", radiusTopEnd = "32", radiusBottomStart = "32", radiusBottomEnd = "32", fontSize = "20")
                .data(text = decreaseButton, enable = enable)
                .action(event = "onClick") {
                    trigger(keyType = "nativeblocks/change_variable", name = "decrease", version = 1)
                    .prop(variableValue = "#SCRIPT
                    const count = {var:count}
                    let result = count
                    if (count >= 1) {
                        result = count - 1
                    } else {
                        result = count
                    }
                    result
                    #ENDSCRIPT")
                    .data(variableKey = count)
                }

                block(keyType = "nativeblocks/text", key = "countText", visibility = visible, version = 1)
                .prop(fontSize = "18", textAlign = "center", width = "128")
                .data(text = count)

                block(keyType = "nativeblocks/button", key = "increaseButton", visibility = visible, version = 1)
                .prop(backgroundColor = "#2563EB", borderColor = "#2563EB", radiusTopStart = "32", radiusTopEnd = "32", radiusBottomStart = "32", radiusBottomEnd = "32", fontSize = "20")
                .data(text = increaseButton, enable = enable)
                .action(event = "onClick") {
                    trigger(keyType = "nativeblocks/change_variable", name = "increase", version = 1)
                    .prop(variableValue = "#SCRIPT
                    const count = {var:count}
                    const result = count + 1
                    result
                    #ENDSCRIPT")
                    .data(variableKey = count)
                }
            }
        }
    }
}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseNBX()

	if frame == nil {
		t.Fatalf("Expected frame to be parsed, got nil: %v", p.Errors())
	}
	if frame.Name != "welcome" {
		t.Errorf("Expected name to be 'welcome', got %s", frame.Name)
	}
	if frame.Route != "/welcome" {
		t.Errorf("Expected route to be '/login', got %s", frame.Route)
	}
	if len(frame.Variables) != 7 {
		t.Errorf("Expected 3 variables, got %d", len(frame.Variables))
	}
	if len(frame.Blocks) == 0 {
		t.Fatalf("Expected at least 1 block, got 0")
	}
}
