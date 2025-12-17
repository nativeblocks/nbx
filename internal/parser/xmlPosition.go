package parser

import "strings"

type Position struct {
	Line   int
	Column int
}

type PositionTracker struct {
	source    string
	positions map[int]Position // byte offset -> line:column
}

func NewPositionTracker(source string) *PositionTracker {
	return &PositionTracker{
		source:    source,
		positions: make(map[int]Position),
	}
}

func (pt *PositionTracker) _getPosition(offset int) Position {
	if pos, ok := pt.positions[offset]; ok {
		return pos
	}

	line := 1
	column := 1

	for i := 0; i < offset && i < len(pt.source); i++ {
		if pt.source[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}

	pos := Position{Line: line, Column: column}
	pt.positions[offset] = pos
	return pos
}

func (pt *PositionTracker) FindElementPosition(elementName, attrValue string) Position {
	searchStr := "<" + elementName

	offset := 0
	for {
		idx := strings.Index(pt.source[offset:], searchStr)
		if idx == -1 {
			return Position{Line: 0, Column: 0}
		}

		// Calculate absolute position
		absIdx := offset + idx

		// Find the end of this element's opening tag
		endIdx := strings.Index(pt.source[absIdx:], ">")
		if endIdx == -1 {
			return Position{Line: 0, Column: 0}
		}

		// Check if this element contains the attribute value
		segment := pt.source[absIdx : absIdx+endIdx]
		if attrValue == "" || strings.Contains(segment, attrValue) {
			return pt._getPosition(absIdx)
		}

		// Move past this element and continue searching
		offset = absIdx + 1
		if offset >= len(pt.source) {
			break
		}
	}

	return Position{Line: 0, Column: 0}
}
