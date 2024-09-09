package parser

import (
	"strings"
	"testing"
)

var testBlocks map[string]string = map[string]string{
	"Empty Paragraph":   "",
	"Heading 3":         "### Heading 3",
	"Broken Heading 3":  "###Heading",
	"Heading 7":         "####### Heading",
	"Code Block":        "```\nCode block line 1\nCode block line 2\n```",
	"Broken Code Block": "```\nCode block line 1\nCode block line 2\n``",
	"Blockquote":        ">>>quote line 1\n>>>quote line 2\n>>>quote line 3",
	"Broken Blockquote": ">>>quote line 1\nquote line 2\n>>>quote line 3",
	"Unordered List":    "* UL item 1\n- UL item 2\n* UL item 3",
	"Ordered List":      "1. OL item 1\n2. OL item 2\n3. OL item 3",
	"Leading Escape":    "\\* This line starts with an escaped asterisk",
}

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
		t.Fatalf("Incorrect block count in result:\n\n%v", result)
	}
}

func TestGetBlockType(t *testing.T) {
	type BlockTest struct {
		Name   string
		Result int
	}

	blocks := map[string]int{
		"Empty Paragraph":   blockTypeParagraph,
		"Heading 3":         blockTypeHeading,
		"Broken Heading 3":  blockTypeParagraph,
		"Heading 7":         blockTypeParagraph,
		"Code Block":        blockTypeCode,
		"Broken Code Block": blockTypeParagraph,
		"Blockquote":        blockTypeQuote,
		"Broken Blockquote": blockTypeParagraph,
		"Unordered List":    blockTypeUnorderedList,
		"Ordered List":      blockTypeOrderedList,
		"Leading Escape":    blockTypeParagraph,
	}

	for key, b := range blocks {
		t.Run(key, func(t *testing.T) {
			result := GetBlockType(testBlocks[key])
			if result != b {
				t.Fatalf("Input:\n%s\nExpected: %v Result: %v", testBlocks[key], b, result)
			}
		})
	}
}

// NOTE: Tag result of "p" represents failure for all non-paragraph blocks
func TestBlockstoHtmlNodes(t *testing.T) {
	blocks := map[string]HtmlNode{
		"Empty Paragraph":  {Tag: "p"},
		"Heading 3":        {Tag: "h3", Value: "Heading 3"},
		"Broken Heading 3": {Tag: "p", Value: testBlocks["Broken Heading 3"]},
		"Heading 7":        {Tag: "p", Value: testBlocks["Heading 7"]},
		"Code Block": {
			Tag:   "pre",
			Value: strings.TrimSpace(testBlocks["Code Block"][4 : len(testBlocks["Code Block"])-3]),
		},
		"Broken Code Block": {Tag: "p", Value: testBlocks["Broken Code Block"]},
		"Blockquote":        {Tag: "blockquote", Value: "quote line 1\nquote line 2\nquote line 3"},
		"Broken Blockquote": {Tag: "p", Value: testBlocks["Broken Blockquote"]},
		"Unordered List":    {Tag: "ul"},
		"Ordered List":      {Tag: "ol"},
		"Leading Escape":    {Tag: "p"},
	}

	for key, b := range blocks {
		t.Run(key, func(t *testing.T) {
			result, _ := BlocksToHTMLNodes([]string{testBlocks[key]})
			if result[0].Tag != b.Tag {
				t.Fatalf("Input:\n%s\nExpected Tag: %s\nResult: %s", testBlocks[key], b.Tag, result[0].Tag)
			}
			if b.Value != "" && result[0].Value != b.Value {
				t.Fatalf("Input:\n%s\nExpected Value: %s\nResult: %s", testBlocks[key], b.Value, result[0].Value)
			}
		})
	}
}
