package test

import (
	"os"
	"strings"
	"testing"

	"github.com/nativeblocks/nbx"
)

func TestToString(t *testing.T) {
	content, err := os.ReadFile("../example/welcome_android.nbx")
	if err != nil {
		t.Fatal(err)
	}

	originalDSL := string(content)

	dslModel, err := nbx.Parse(originalDSL)
	if err != nil {
		t.Fatal("Parse error:", err)
	}

	reconstructedDSL := nbx.ToString(dslModel)

	//t.Logf("=== ORIGINAL DSL ===\n%s", originalDSL)
	//t.Logf("\n=== RECONSTRUCTED DSL ===\n%s", reconstructedDSL)

	if strings.TrimSpace(originalDSL) == strings.TrimSpace(reconstructedDSL) {
		t.Log("originalDsl and stringify dsl are identical")
	}

	_, err = nbx.Parse(reconstructedDSL)
	if err != nil {
		t.Fatal("Reconstructed DSL cannot be parsed:", err)
	}

	t.Log("Reconstructed DSL is valid and parseable!")
}
