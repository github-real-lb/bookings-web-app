package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewText(t *testing.T) {
	text := NewText()
	require.IsType(t, &Text{}, text)
	assert.Len(t, text.Lines, 0)
}

func TestText_Add(t *testing.T) {
	text := NewText().
		Add("part a").
		Add("part b").
		Add("part c")

	require.IsType(t, &Text{}, text)
	assert.Len(t, text.Lines, 1)
	assert.Equal(t, "part apart bpart c", text.Lines[0])

	text = NewText()
	text.Add("part a")
	text.Add("part b")
	text.Add("part c")

	require.IsType(t, &Text{}, text)
	assert.Len(t, text.Lines, 1)
	assert.Equal(t, "part apart bpart c", text.Lines[0])
}

func TestText_AddLine(t *testing.T) {
	text := NewText().
		AddLine("line a").
		AddLine("line b").
		AddLine("line c")

	require.IsType(t, &Text{}, text)
	assert.Len(t, text.Lines, 3)

	assert.Equal(t, "line a", text.Lines[0])
	assert.Equal(t, "line b", text.Lines[1])
	assert.Equal(t, "line c", text.Lines[2])

	text = NewText()
	text.AddLine("line a")
	text.AddLine("line b")
	text.AddLine("line c")

	require.IsType(t, &Text{}, text)
	assert.Len(t, text.Lines, 3)

	assert.Equal(t, "line a", text.Lines[0])
	assert.Equal(t, "line b", text.Lines[1])
	assert.Equal(t, "line c", text.Lines[2])
}

func TestText_AddLineIndent(t *testing.T) {
	text := NewText().
		AddLineIndent("line a", "\t").
		AddLineIndent("line b", "\t").
		AddLineIndent("line c", "\t")

	require.IsType(t, &Text{}, text)
	assert.Len(t, text.Lines, 3)

	assert.Equal(t, "\tline a", text.Lines[0])
	assert.Equal(t, "\tline b", text.Lines[1])
	assert.Equal(t, "\tline c", text.Lines[2])

	text = NewText()
	text.AddLineIndent("line a", "\t")
	text.AddLineIndent("line b", "\t")
	text.AddLineIndent("line c", "\t")

	require.IsType(t, &Text{}, text)
	assert.Len(t, text.Lines, 3)

	assert.Equal(t, "\tline a", text.Lines[0])
	assert.Equal(t, "\tline b", text.Lines[1])
	assert.Equal(t, "\tline c", text.Lines[2])
}

func TestText_String(t *testing.T) {
	text := NewText().
		AddLineIndent("line a", "\t").
		AddLineIndent("line b", "\t").
		AddLineIndent("line c", "\t")

	require.IsType(t, &Text{}, text)
	assert.Len(t, text.Lines, 3)

	actual := text.String()
	expected := "\tline a\n\tline b\n\tline c\n"
	assert.Equal(t, expected, actual)
}

func TestText_Error(t *testing.T) {
	text := NewText().
		AddLineIndent("line a", "\t").
		AddLineIndent("line b", "\t").
		AddLineIndent("line c", "\t")

	require.IsType(t, &Text{}, text)
	assert.Len(t, text.Lines, 3)

	actual := text.Error()
	expected := "\tline a\n\tline b\n\tline c\n"
	assert.Equal(t, expected, actual)
}
