# NBX

A Go library to parse and convert the Nativeblocks eXchange (NBX) Domain-Specific Language (DSL) for describing UI frames, variables, blocks, slots, actions, and triggers to structured JSON, and vice versa.

## What is NBX?

NBX lets you define UI screens, layouts, and their logic in a human-readable DSL text format, which can be converted to JSON for further processing (codegen, cross-platform, rendering). Conversion is bidirectional: NBX DSL ⇄ JSON.

---

## Concepts

- **Frame**: The root structure representing a UI screen or component.
- **Variable**: Named data available in the frame (flags, form fields, intermediate data).
- **Block**: UI element such as containers, buttons, inputs, images, etc. Blocks can be nested.
- **Slot**: Named regions in a block to inject other blocks (“children” into layouts).
- **Action**: Response logic to UI events (e.g., onClick, onChange). Belongs to a Block.
- **Trigger**: The invocation of an effect (function) in response to an event inside an Action, Triggers can conditionally run more triggers via `.then("NEXT") { ... }` blocks, for handling success, failure, or custom logic.
---

## Example

Example NBX DSL describing a login UI with logic:

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

## Syntax Overview

- **Frame Declaration**  
  ```
  frame(name = "screenName", route = "/route") { ... }
  ```

- **Variable Declaration**  
  ```
  var variableName: TYPE = value
  ```

- **Block Declaration**  
  ```
  block(keyType = "TYPE", key = "name", visibility = someVariable, version = 1)
  ```

- **Slot Injection**  
  ```
  .slot("slotName") { ... }
  ```

- **Data Assignment**  
  ```
  .data(key = value)
  // or
  .data(
      key1 = value1, 
      key2 = value2
  )
  ```
  
- **Property Assignment (single or multi-device)**  
  ```
  .prop(
      property1 = value1,
      property2 = value2
  )
  // or
  .prop(
      property1 = (mobile = "NONE", tablet = "NONE", desktop = "NONE"),
      property2 = (mobile = "NONE", tablet = "NONE", desktop = "NONE")
  )
  // or
  .prop(
      property1 = (value = "NONE"),
      property2 = (value = "NONE")
  )
  ```

- **Event Action**  
  ```
  .action(event = "eventName") { ... }
  ```
  (Multiple triggers can be handled inside.)

- **Trigger**  
  ```
  trigger(keyType = "TYPE", name = "description")
  .then("NEXT") { ... }
  ```

---

## Usage

Use the root `nbx` package. All implementation details are in `internal/` (not for import).

### Convert JSON to NBX DSL

```go
import "github.com/nativeblocks/nbx"

frameDSL := nbx.ToDSL(frameJson)
```

### Convert NBX DSL to JSON

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