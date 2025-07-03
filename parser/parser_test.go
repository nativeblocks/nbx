package parser

import (
	"github.com/nativeblocks/nbx/lexer"
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
	frame := p.ParseSDUI()

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
	frame := p.ParseSDUI()

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
	frame := p.ParseSDUI()

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
            .assignData(text = username)
            block(keyType = "INPUT", key = "password", visibility = visible)
            .assignProperty(color = (valueMobile = "NONE", valueTablet = "NONE", valueDesktop = "NONE"))
            .assignData(text = password)
            .action(event = "onTextChange") {
                trigger(keyType = "VALIDATE", name = "validate password")
                .then("FAILURE") {
                    trigger(keyType = "SHOW_ERROR", name = "show error 1")
                    trigger(keyType = "CHANGE_COLOR", name = "change color to red")
                    .assignProperty(color = "RED")
                }
                .then("SUCCESS") {
                    trigger(keyType = "SHOW_OK", name = "show ok")
                    .assignProperty(color = "GREEN")
                }
            }
        }    
    }   
}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseSDUI()

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
            .assignData(text = username)
            block(keyType = "nativeblocks/TEXT_FIELD", key = "password", visibility = visible)
            .assignProperty(textColor = (valueMobile = "NONE", valueTablet = "NONE", valueDesktop = "NONE"))
            .assignData(text = password)
            .action(event = "onTextChange") {
                trigger(keyType = "nativeblocks/CHANGE_BLOCK_PROPERTY", name = "show error")
                .assignProperty(propertyKey = "textColor")
                .assignProperty(propertyValueDesktop = "RED")
				trigger(keyType = "nativeblocks/CHANGE_BLOCK_PROPERTY", name = "show success")
				.assignProperty(propertyKey = "textColor")
				.assignProperty(propertyValueDesktop = "GREEN")
            }
        }
    }
}
`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseSDUI()

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
            .assignData(text = username)
        }
    }
}`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseSDUI()

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
    .assignData(text = password)
    .action(event = "onTextChange") {
        trigger(keyType = "VALIDATE", name = "validate password")
        .then("SUCCESS") {
            trigger(keyType = "SHOW_OK", name = "show ok")
            .assignProperty(color = "GREEN")
        }
    }
}`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	frame := p.ParseSDUI()

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
	frame := p.ParseSDUI()

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
