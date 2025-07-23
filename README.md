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

    block(keyType = "ROOT", key = "root", visibility = visible, version = 1)
    .slot("content") {
        block(keyType = "COLUMN", key = "main", visibility = visible, version = 1)
        .slot("content") {
            block(keyType = "INPUT", key = "username", visibility = visible, version = 1)
            .data(text = username)
            block(keyType = "INPUT", key = "password", visibility = visible, version = 1)
            .prop(fontSize = (mobile = "14", tablet = "14", desktop = "14"))
            .data(text = password)
            .action(event = "onTextChange") {
                trigger(keyType = "VALIDATE", name = "validate password", version = 1)
                .then("FAILURE") {
                    trigger(keyType = "SHOW_ERROR", name = "show error 1", version = 1)
                    trigger(keyType = "CHANGE_COLOR", name = "change color to red", version = 1)
                    .prop(color = "RED")
                }
                .then("SUCCESS") {
                    trigger(keyType = "SHOW_OK", name = "show ok", version = 1)
                    .prop(color = "GREEN")
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

- **Data/Property**:  
  `.data(key = value)`  
  `.prop(key = value, key2 = value2)`  
  You can assign for different devices:  
  `.prop(color = (mobile = "NONE", tablet = "NONE", desktop = "NONE"))`

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
