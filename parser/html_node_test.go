package parser

import (
	"testing"
)

func TestEscapeMD(t *testing.T) {
	const input = "\\* \\_ \\*\\* \\_\\_ \\\\"
	const expected = "* _ ** __ \\"
	node := HtmlNode{
		Tag:   "p",
		Value: input,
	}

	node.UnescapeMD()
	if node.Value != expected {
		t.Fatalf("Escape character processing failed.\nExpected: %s\nResult: %s", expected, node.Value)
	}
}
