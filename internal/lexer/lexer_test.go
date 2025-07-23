package lexer

import (
	"testing"
)

func TestLexer_ComplexInput(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login",
) {
    var visible: BOOLEAN = true
    var username: STRING = ""
    var password: STRING = ""

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "column", key = "main", visibility = visible)
        .slot("content") {
            block(keyType = "input", key = "username", visibility = visible)
            .assignData(text = username)
            block(keyType = "input", key = "password", visibility = visible)
            .prop(
                color = (mobile = "NONE", valueTablet = "NONE", valueDesktop = "NONE")
            )
            .assignData(text = password)
            .action(event = "onTextChange") {
                trigger(keyType = "validate", name = "validate password")
                .then("FAILURE") {
                    trigger(keyType = "show_error", name = "show error 1")
                    trigger(keyType = "change_color", name = "change color to red")
                    .prop(
                        color = "RED"
                    )
                }
                .then("SUCCESS") {
                    trigger(keyType = "show_ok", name = "show ok")
                    .prop(
                        color = "GREEN"
                    )
                }
            }
        }    
    }   
}
`

	l := NewLexer(input)
	seen := 0

	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_EOF {
			break
		}

		if tok.Type == TOKEN_ILLEGAL {
			t.Fatalf("Unexpected illegal token: %v at line %d, col %d", tok.Literal, tok.Line, tok.Column)
		}

		seen++
	}

	if seen < 50 {
		t.Fatalf("Expected many tokens, but only got %d", seen)
	}
}

func TestLexer_ComplexInput2(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login",
) {
    var visible: BOOLEAN = true
    var username: STRING = ""
    var password: STRING = ""

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "nativeblocks/column", key = "main", visibility = visible)
        .slot("content") {
            block(keyType = "nativeblocks/text_field", key = "username", visibility = visible)
            .assignData(text = username)
            block(keyType = "nativeblocks/text_field", key = "password", visibility = visible)
            .prop(
                textColor = (mobile = "NONE", valueTablet = "NONE", valueDesktop = "NONE")
            )
            .assignData(text = password)
            .action(event = "onTextChange") {
                trigger(keyType = "nativeblocks/change_block_property", name = "show error")
                .prop(
                    propertyKey = "textColor"
                )
                .prop(
                    propertyValueDesktop = "RED"
                )
				trigger(keyType = "nativeblocks/change_block_property", name = "show success")
				.prop(
                    propertyKey = "textColor"
                )
				.prop(
                    propertyValueDesktop = "GREEN"
                )
            }
        }
    }
}
`

	l := NewLexer(input)
	seen := 0

	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_EOF {
			break
		}

		if tok.Type == TOKEN_ILLEGAL {
			t.Fatalf("Unexpected illegal token: %v at line %d, col %d", tok.Literal, tok.Line, tok.Column)
		}

		seen++
	}

	if seen < 50 {
		t.Fatalf("Expected many tokens, but only got %d", seen)
	}
}

func TestLexerKeyword_ComplexInput(t *testing.T) {

	input := `
frame(
    name = "login",
    route = "/login",
) {
    var visible: BOOLEAN = true
    var username: STRING = ""
    var password: STRING = ""

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "column", key = "main", visibility = visible)
        .slot("content") {
            block(keyType = "input", key = "username", visibility = visible)
            .data(text = username)
            block(keyType = "input", key = "password", visibility = visible)
            .prop(
                color = (mobile = "NONE", tablet = "NONE", desktop = "NONE")
            )
            .data(text = password)
            .action(event = "onTextChange") {
                trigger(keyType = "validate", name = "validate password")
                .then("FAILURE") {
                    trigger(keyType = "show_error", name = "show error 1")
                    trigger(keyType = "change_color", name = "change color to red")
                    .prop(
                        color = "RED"
                    )
                }
                .then("SUCCESS") {
                    trigger(keyType = "show_ok", name = "show ok")
                    .prop(
                        color = "GREEN"
                    )
                }
            }
        }    
    }   
}
`

	l := NewLexer(input)
	seen := 0

	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_EOF {
			break
		}

		if tok.Type == TOKEN_ILLEGAL {
			t.Fatalf("Unexpected illegal token: %v at line %d, col %d", tok.Literal, tok.Line, tok.Column)
		}

		seen++
	}

	if seen < 50 {
		t.Fatalf("Expected many tokens, but only got %d", seen)
	}
}

func TestLexerKeyword_ComplexInput2(t *testing.T) {
	input := `
frame(
    name = "login",
    route = "/login",
) {
    var visible: BOOLEAN = true
    var username: STRING = ""
    var password: STRING = ""

    block(keyType = "ROOT", key = "root", visibility = visible)
    .slot("content") {
        block(keyType = "nativeblocks/column", key = "main", visibility = visible)
        .slot("content") {
            block(keyType = "nativeblocks/text_field", key = "username", visibility = visible)
            .data(text = username)
            block(keyType = "nativeblocks/text_field", key = "password", visibility = visible)
            .prop(
                textColor = (mobile = "NONE", tablet = "NONE", desktop = "NONE")
            )
            .data(text = password)
            .action(event = "onTextChange") {
                trigger(keyType = "nativeblocks/change_block_property", name = "show error")
                .prop(
                    propertyKey = "textColor"
                )
                .prop(
                    propertydesktop = "RED"
                )
				trigger(keyType = "nativeblocks/change_block_property", name = "show success")
				.prop(
                    propertyKey = "textColor"
                )
				.prop(
                    propertydesktop = "GREEN"
                )
            }
        }
    }
}
`

	l := NewLexer(input)
	seen := 0

	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_EOF {
			break
		}

		if tok.Type == TOKEN_ILLEGAL {
			t.Fatalf("Unexpected illegal token: %v at line %d, col %d", tok.Literal, tok.Line, tok.Column)
		}

		seen++
	}

	if seen < 50 {
		t.Fatalf("Expected many tokens, but only got %d", seen)
	}
}

func TestLexer_WithComments(t *testing.T) {
	input := `
frame(
    name = "dashboard",
    route = "/dashboard",
) {
	// Counter variable
    var counter: INT = 42
}
`

	l := NewLexer(input)
	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_ILLEGAL {
			t.Fatalf("Unexpected illegal token: %v at line %d, col %d", tok.Literal, tok.Line, tok.Column)
		}
		if tok.Type == TOKEN_EOF {
			break
		}
	}
}

func TestLexer_SyntaxErrors(t *testing.T) {
	input := `
frame(
    name = "broken"
    missing_comma route = "/broken"
) {
    var invalid@name: STRING = "error"
    block(missingKey = ) // Missing value
    .slot( // Missing slot name
}
`
	l := NewLexer(input)
	illegalCount := 0

	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_ILLEGAL {
			illegalCount++
		}
		if tok.Type == TOKEN_EOF {
			break
		}
	}

	if illegalCount == 0 {
		t.Fatalf("Expected illegal tokens but found none")
	}
}

func TestLexer_NumberTypes(t *testing.T) {
	input := `
frame(name = "numbers") {
    var int1: INT = 42
    var int2: INT = -42
    var long1: LONG = 9223372036854775807
    var float1: FLOAT = 3.14
    var float2: FLOAT = -3.14
    var double1: DOUBLE = 3.141592653589793
}
`

	l := NewLexer(input)
	numbers := make(map[TokenType]int)

	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_EOF {
			break
		}

		if tok.Type == TOKEN_INT || tok.Type == TOKEN_LONG || tok.Type == TOKEN_FLOAT || tok.Type == TOKEN_DOUBLE {
			numbers[tok.Type]++
		}
	}

	if numbers[TOKEN_INT] < 2 || numbers[TOKEN_LONG] < 1 ||
		numbers[TOKEN_FLOAT] < 2 || numbers[TOKEN_DOUBLE] < 1 {
		t.Fatalf("Not all number types were correctly identified")
	}
}

func TestLexer_NestedStructures(t *testing.T) {
	input := `
frame(name = "nested") {
    block(keyType = "container") {
        block(keyType = "row") {
            block(keyType = "column") {
                block(keyType = "text") 
                .action(event = "onClick") {
                    trigger(keyType = "navigation") {
                        trigger(keyType = "push") {
                            trigger(keyType = "animate")
                        }
                    }
                }
            }
        }
    }
}
`

	l := NewLexer(input)
	depth := 0
	maxDepth := 0

	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_EOF {
			break
		}

		if tok.Type == TOKEN_LBRACE {
			depth++
			if depth > maxDepth {
				maxDepth = depth
			}
		} else if tok.Type == TOKEN_RBRACE {
			depth--
		}
	}

	if maxDepth < 5 {
		t.Fatalf("Expected deeply nested structure with depth >= 5, got %d", maxDepth)
	}
}
