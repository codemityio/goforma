package doc

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/input.txt
var inputText string

//go:embed testdata/output.txt
var outputText string

func TestDefaultParser_Parse(t *testing.T) {
	parser := New()
	require.NotNil(t, parser)

	assert.Equal(t, outputText, parser.Parse(strings.Split(inputText, "\n")))
}
