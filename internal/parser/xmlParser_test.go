package parser

import (
	"testing"
)

func TestParseXML_Simple(t *testing.T) {
	xmlInput := `
<frame name="HomePage" route="/home">
  <var key="title" type="STRING" value="Hello World" />
  <block keyType="ROOT" key="root">
  </block>
</frame>`

	frame, errs := ParseXML(xmlInput)

	if len(errs) > 0 {
		t.Fatalf("Unexpected errors: %v", errs)
	}

	if frame.Name != "HomePage" {
		t.Errorf("Expected frame name 'HomePage', got '%s'", frame.Name)
	}

	if frame.Route != "/home" {
		t.Errorf("Expected route '/home', got '%s'", frame.Route)
	}

	if frame.Type != "FRAME" {
		t.Errorf("Expected type 'FRAME', got '%s'", frame.Type)
	}

	if len(frame.Variables) != 1 {
		t.Fatalf("Expected 1 variable, got %d", len(frame.Variables))
	}

	v := frame.Variables[0]
	if v.Key != "title" {
		t.Errorf("Expected variable key 'title', got '%s'", v.Key)
	}
	if v.Type != "STRING" {
		t.Errorf("Expected variable type 'STRING', got '%s'", v.Type)
	}
	if v.Value != "Hello World" {
		t.Errorf("Expected variable value 'Hello World', got '%s'", v.Value)
	}

	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(frame.Blocks))
	}

	b := frame.Blocks[0]
	if b.KeyType != "ROOT" {
		t.Errorf("Expected block keyType 'ROOT', got '%s'", b.KeyType)
	}
	if b.Key != "root" {
		t.Errorf("Expected block key 'root', got '%s'", b.Key)
	}
}

func TestParseXML_WithProperties(t *testing.T) {
	xmlInput := `<frame name="test" route="/test">
  <block keyType="ROOT" key="root">
    <prop key="fontSize" value="16" />
    <prop key="layout" mobile="vertical" tablet="horizontal" desktop="horizontal" />
  </block>
</frame>`

	frame, errs := ParseXML(xmlInput)

	if len(errs) > 0 {
		t.Fatalf("Unexpected errors: %v", errs)
	}

	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(frame.Blocks))
	}

	b := frame.Blocks[0]
	if len(b.Properties) != 2 {
		t.Fatalf("Expected 2 properties, got %d", len(b.Properties))
	}

	// Check single-value property
	p1 := b.Properties[0]
	if p1.Key != "fontSize" {
		t.Errorf("Expected property key 'fontSize', got '%s'", p1.Key)
	}
	if p1.ValueMobile != "16" || p1.ValueTablet != "16" || p1.ValueDesktop != "16" {
		t.Errorf("Expected all device values to be '16', got mobile='%s', tablet='%s', desktop='%s'",
			p1.ValueMobile, p1.ValueTablet, p1.ValueDesktop)
	}

	// Check multi-device property
	p2 := b.Properties[1]
	if p2.Key != "layout" {
		t.Errorf("Expected property key 'layout', got '%s'", p2.Key)
	}
	if p2.ValueMobile != "vertical" {
		t.Errorf("Expected mobile value 'vertical', got '%s'", p2.ValueMobile)
	}
	if p2.ValueTablet != "horizontal" {
		t.Errorf("Expected tablet value 'horizontal', got '%s'", p2.ValueTablet)
	}
	if p2.ValueDesktop != "horizontal" {
		t.Errorf("Expected desktop value 'horizontal', got '%s'", p2.ValueDesktop)
	}
}

func TestParseXML_WithSlots(t *testing.T) {
	xmlInput := `<frame name="test" route="/test">
  <block keyType="ROOT" key="root">
    <slot name="content">
      <block keyType="nativeblocks/text" key="title">
        <data key="text" value="title_text" />
      </block>
    </slot>
  </block>
</frame>`

	frame, errs := ParseXML(xmlInput)

	if len(errs) > 0 {
		t.Fatalf("Unexpected errors: %v", errs)
	}

	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(frame.Blocks))
	}

	b := frame.Blocks[0]
	if len(b.Slots) != 1 {
		t.Fatalf("Expected 1 slot, got %d", len(b.Slots))
	}

	s := b.Slots[0]
	if s.Slot != "content" {
		t.Errorf("Expected slot name 'content', got '%s'", s.Slot)
	}

	if len(b.Blocks) != 1 {
		t.Fatalf("Expected 1 child block, got %d", len(b.Blocks))
	}

	child := b.Blocks[0]
	if child.KeyType != "nativeblocks/text" {
		t.Errorf("Expected child keyType 'nativeblocks/text', got '%s'", child.KeyType)
	}
	if child.Key != "title" {
		t.Errorf("Expected child key 'title', got '%s'", child.Key)
	}
	if child.Slot != "content" {
		t.Errorf("Expected child slot 'content', got '%s'", child.Slot)
	}

	if len(child.Data) != 1 {
		t.Fatalf("Expected 1 data binding, got %d", len(child.Data))
	}

	d := child.Data[0]
	if d.Key != "text" {
		t.Errorf("Expected data key 'text', got '%s'", d.Key)
	}
	if d.Value != "title_text" {
		t.Errorf("Expected data value 'title_text', got '%s'", d.Value)
	}
}

func TestParseXML_WithActions(t *testing.T) {
	xmlInput := `<frame name="test" route="/test">
  <block keyType="ROOT" key="root">
    <action event="onClick">
      <trigger keyType="nativeblocks/navigate" name="goHome" version="1">
        <prop key="route" value="/home" />
        <data key="userId" value="userId" />
      </trigger>
    </action>
  </block>
</frame>`

	frame, errs := ParseXML(xmlInput)

	if len(errs) > 0 {
		t.Fatalf("Unexpected errors: %v", errs)
	}

	if len(frame.Blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(frame.Blocks))
	}

	b := frame.Blocks[0]
	if len(b.Actions) != 1 {
		t.Fatalf("Expected 1 action, got %d", len(b.Actions))
	}

	a := b.Actions[0]
	if a.Event != "onClick" {
		t.Errorf("Expected event 'onClick', got '%s'", a.Event)
	}
	if a.Key != "root" {
		t.Errorf("Expected action key 'root', got '%s'", a.Key)
	}

	if len(a.Triggers) != 1 {
		t.Fatalf("Expected 1 trigger, got %d", len(a.Triggers))
	}

	tr := a.Triggers[0]
	if tr.KeyType != "nativeblocks/navigate" {
		t.Errorf("Expected trigger keyType 'nativeblocks/navigate', got '%s'", tr.KeyType)
	}
	if tr.Name != "goHome" {
		t.Errorf("Expected trigger name 'goHome', got '%s'", tr.Name)
	}
	if tr.IntegrationVersion != 1 {
		t.Errorf("Expected trigger version 1, got %d", tr.IntegrationVersion)
	}

	if len(tr.Properties) != 1 {
		t.Fatalf("Expected 1 trigger property, got %d", len(tr.Properties))
	}

	p := tr.Properties[0]
	if p.Key != "route" {
		t.Errorf("Expected property key 'route', got '%s'", p.Key)
	}
	if p.Value != "/home" {
		t.Errorf("Expected property value '/home', got '%s'", p.Value)
	}

	if len(tr.Data) != 1 {
		t.Fatalf("Expected 1 trigger data, got %d", len(tr.Data))
	}

	d := tr.Data[0]
	if d.Key != "userId" {
		t.Errorf("Expected data key 'userId', got '%s'", d.Key)
	}
	if d.Value != "userId" {
		t.Errorf("Expected data value 'userId', got '%s'", d.Value)
	}
}

func TestParseXML_WithConditionalFlow(t *testing.T) {
	xmlInput := `<frame name="test" route="/test">
  <block keyType="ROOT" key="root">
    <action event="onClick">
      <trigger keyType="VALIDATE" name="check">
        <then value="SUCCESS">
          <trigger keyType="NAVIGATE" name="goHome">
            <prop key="route" value="/home" />
          </trigger>
        </then>
        <then value="FAILURE">
          <trigger keyType="SHOW_ERROR" name="showErr">
            <prop key="message" value="Failed" />
          </trigger>
        </then>
      </trigger>
    </action>
  </block>
</frame>`

	frame, errs := ParseXML(xmlInput)

	if len(errs) > 0 {
		t.Fatalf("Unexpected errors: %v", errs)
	}

	b := frame.Blocks[0]
	a := b.Actions[0]
	tr := a.Triggers[0]

	if len(tr.Triggers) != 2 {
		t.Fatalf("Expected 2 nested triggers, got %d", len(tr.Triggers))
	}

	// Check SUCCESS trigger
	successTr := tr.Triggers[0]
	if successTr.Then != "SUCCESS" {
		t.Errorf("Expected then value 'SUCCESS', got '%s'", successTr.Then)
	}
	if successTr.KeyType != "NAVIGATE" {
		t.Errorf("Expected keyType 'NAVIGATE', got '%s'", successTr.KeyType)
	}

	// Check FAILURE trigger
	failureTr := tr.Triggers[1]
	if failureTr.Then != "FAILURE" {
		t.Errorf("Expected then value 'FAILURE', got '%s'", failureTr.Then)
	}
	if failureTr.KeyType != "SHOW_ERROR" {
		t.Errorf("Expected keyType 'SHOW_ERROR', got '%s'", failureTr.KeyType)
	}
}

func TestParseXML_MissingRequiredFields(t *testing.T) {
	// Missing name
	xmlInput1 := `<frame route="/test"></frame>`
	_, errs1 := ParseXML(xmlInput1)
	if len(errs1) == 0 {
		t.Error("Expected error for missing name")
	}

	// Missing route
	xmlInput2 := `<frame name="test"></frame>`
	_, errs2 := ParseXML(xmlInput2)
	if len(errs2) == 0 {
		t.Error("Expected error for missing route")
	}

	// Empty content
	_, errs3 := ParseXML("")
	if len(errs3) == 0 {
		t.Error("Expected error for empty content")
	}
}
