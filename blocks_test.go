package main

import (
	"testing"
)

func TestMarkdownBlocks(t *testing.T) {
	blocks := []string{
		"# This is a heading",
		"This is a paragraph of text. It has some **bold** and *italic* words inside of it.",
		"* This is a list item\n* This is another list item",
	}
	var md string = blocks[0]
	for i := 1; i < len(blocks); i++ {
		md += "\n\n" + blocks[i]
	}

	result := ParseMDBlocks(md)
	if len(result) != len(blocks) {
		t.Fatalf("Incorrect block count in result")
	}
}

func TestGetBlockType(t *testing.T) {
	type BlockTest struct {
		Name   string
		Result int
		Block  string
	}

	blocks := []BlockTest{
		{Name: "Empty Paragraph", Result: blockTypeParagraph, Block: ""},
		{Name: "Heading 3", Result: blockTypeHeading, Block: "### Heading 3"},
		{Name: "Broken heading 3", Result: blockTypeParagraph, Block: "###Heading"},
		{Name: "Heading 7", Result: blockTypeParagraph, Block: "####### Heading"},
		{Name: "Code block", Result: blockTypeCode, Block: "```\nCode block line 1\nCode block line 2\n```"},
		{Name: "Broken code block", Result: blockTypeParagraph, Block: "```\nCode block line 1\nCode block line 2\n``"},
		{Name: "Blockquote", Result: blockTypeQuote, Block: ">>>quote line 1\n>>>quote line 2\n>>>quote line 3"},
		{Name: "Broken Blockquote", Result: blockTypeParagraph, Block: ">>>quote line 1\nquote line 2\n>>>quote line 3"},
		{Name: "Unordered List", Result: blockTypeUnorderedList, Block: "* UL item 1\n- UL item 2\n* UL item 3"},
		{Name: "Ordered list", Result: blockTypeOrderedList, Block: "1. OL item 1\n2. OL item 2\n3. OL item 3"},
	}

	for _, b := range blocks {
		t.Run(b.Name, func(t *testing.T) {
			result := GetBlockType(b.Block)
			if result != b.Result {
				t.Fatalf("GetBlockType() failed for input:\n%s\nExpected: %v Result: %v", b.Block, b.Result, result)
			}
		})
	}
}
