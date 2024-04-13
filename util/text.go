package util

import (
	"fmt"
	"strings"
)

// Text holds lines of string
type Text struct {
	Lines []string
}

// NewErrorInfo returns an empty ErrorInfo
func NewText() *Text {
	return &Text{}
}

// Add adds v to last line in t.
func (t *Text) Add(v any) *Text {
	if len(t.Lines) == 0 {
		return t.AddLine(v)
	}

	i := len(t.Lines) - 1
	t.Lines[i] = fmt.Sprint(t.Lines[i], v)

	return t
}

// AddLine adds v to t as a new line.
func (t *Text) AddLine(v any) *Text {
	t.Lines = append(t.Lines, fmt.Sprint(v))
	return t
}

// AddLine adds v to t as a new line with indent.
func (t *Text) AddLineIndent(v any, indent string) *Text {
	t.Lines = append(t.Lines, fmt.Sprint(indent, v))
	return t
}

// String creates a string from t
func (t *Text) String() string {
	var sb strings.Builder

	for _, v := range t.Lines {
		sb.WriteString(v)
		sb.WriteString("\n")
	}

	return sb.String()
}

// Error creates an error from e.
func (t Text) Error() string {
	return t.String()
}
