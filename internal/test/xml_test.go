package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/nativeblocks/nbx"
)

func TestXMLIntegration(t *testing.T) {
	fmt.Println("=== Test 1: Parse Simple XML ===")
	simpleXML, err := os.ReadFile("../example/simple.xml")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	frame, errs := nbx.Parse(string(simpleXML))
	if len(errs) > 0 {
		fmt.Printf("Parse errors: %v\n", errs.FormatAll())
		return
	}

	fmt.Printf("✓ Parsed frame: %s (route: %s)\n", frame.Name, frame.Route)
	fmt.Printf("✓ Variables: %d\n", len(frame.Variables))
	fmt.Printf("✓ Blocks: %d\n", len(frame.Blocks))

	fmt.Println("\n=== Test 2: Parse Login XML ===")
	loginXML, err := os.ReadFile("../example/login.xml")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	loginFrame, errs := nbx.Parse(string(loginXML))
	if len(errs) > 0 {
		fmt.Printf("Parse errors: %v\n", errs.FormatAll())
		return
	}

	fmt.Printf("✓ Parsed frame: %s (route: %s)\n", loginFrame.Name, loginFrame.Route)
	fmt.Printf("✓ Variables: %d\n", len(loginFrame.Variables))
	fmt.Printf("✓ Blocks: %d\n", len(loginFrame.Blocks))

	if len(loginFrame.Blocks) > 0 && len(loginFrame.Blocks[0].Actions) > 0 {
		fmt.Printf("✓ Actions: %d\n", len(loginFrame.Blocks[0].Actions))
	}

	fmt.Println("\n=== Test 3: Convert XML to DSL ===")
	dslString := nbx.ToString(frame)
	fmt.Println("✓ Converted to DSL:")
	fmt.Println(dslString[:200] + "...")

	fmt.Println("\n=== Test 4: Round-trip XML -> Model -> XML ===")
	xmlOutput := nbx.ToXML(frame)
	fmt.Println("✓ Converted back to XML:")
	fmt.Println(xmlOutput[:200] + "...")

	_, errs = nbx.ParseXML(xmlOutput)
	if len(errs) > 0 {
		fmt.Printf("✗ Round-trip failed: %v\n", errs.FormatAll())
		return
	}
	fmt.Println("✓ Round-trip successful!")

	fmt.Println("\n=== Test 5: Format Detection ===")
	fmt.Printf("Simple XML format: %s\n", nbx.DetectFormat(string(simpleXML)))

	dslExample := `frame(name = "test", route = "/test") {}`
	fmt.Printf("DSL format: %s\n", nbx.DetectFormat(dslExample))

	fmt.Println("\n✓ All tests passed!")
}
