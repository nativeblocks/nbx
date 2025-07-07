# NBX

A Go library for parsing and converting the NBX DSL (Domain-Specific Language) for describing UI frames, variables, blocks, and actions into structured JSON, and vice versa.

## What is this?

This project lets you describe UI screens and their logic in a readable text format (DSL), and then turn that into JSON for use in other systems. You can also go the other way: take JSON and turn it back into the DSL.

---

## Concepts

- **Frame**: The top-level container. Think of it as a screen or page.
- **Variable**: A named value you can use in your frame (like a flag or a field value).
- **Block**: A UI element (like a button, input, or container). Blocks can be nested.
- **Slot**: A named area inside a block where you can put other blocks.
- **Action**: Something that happens in response to an event (like a button click).
- **Trigger**: A function execution on a specific events (NEXT, END, FAILURE, SUCCESS).

---

## Example

Here's an example of the DSL:

```
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
}
```

---

## Syntax

- **Frame**:  
  `frame(name = "screenName", route = "/route") { ... }`

- **Variable**:  
  `var variableName: TYPE = value`

- **Block**:  
  `block(keyType = "TYPE", key = "name", visibility = variable)`

- **Slot**:  
  `.slot("slotName") { ... }`

- **Assign Data/Property**:  
  `.assignData(key = value)`  
  You can assign for different devices:  
  `.assignProperty(color = (valueMobile = "NONE", valueTablet = "NONE", valueDesktop = "NONE"))`

- **Action**:  
  `.action(event = "eventName") { ... }`

- **Trigger**:  
  `trigger(keyType = "TYPE", name = "a function name or description of what this function does")`  
  Triggers can have `.then("NEXT") { ... }` blocks for branching.

- **Assign Data/Property**:  
  `.assignData(key = value)`
  `.assignProperty(key = value)`

---

## How to Use

### Public API

You should only use the root `nbx` package. All implementation details are hidden in `internal/` and cannot be imported directly.


#### Convert JSON to DSL struct

```go
import "github.com/nativeblocks/nbx"

frameDSL := nbx.ToDSL(frameJson)
```

#### Convert DSL struct to JSON

```go
import "github.com/nativeblocks/nbx"

jsonFrame, err := nbx.ToJSON(frameDSL, schemaString, "")
if err != nil {
    // handle error
}
```

---

## License

MIT

---
