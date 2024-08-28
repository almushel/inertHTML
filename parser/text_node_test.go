package parser

import (
	"fmt"
	"testing"
)

func TestDelimSplit(t *testing.T) {
	node := TextNode{
		TextType: textTypeText,
		Text:     "Before code tag `inside of code tag` after code tag **Bold text**",
	}

	var result TextNodeSlice
	result = append(result, node)

	result, _ = result.Split("`", textTypeCode)
	if len(result) != 3 || result[1].TextType != textTypeCode {
		t.Fatalf("Incorrect Split() output for delimiter %s: %s", "`", fmt.Sprintf("%#v", result))
	}

	result, _ = result.Split("**", textTypeBold)
	if len(result) != 4 || result[3].TextType != textTypeBold {
		t.Fatalf("Incorrect Split() output for delimiter %s: %s", "**", fmt.Sprintf("%#v", result))
	}
}

func TestImgSplit(t *testing.T) {
	node := TextNode{
		TextType: textTypeText,
		Text: `Before image1 ![img1 alt text](http://img1.url) after img1,
		before img2 ![img2 alt text](http://img2.url)`,
	}
	var result TextNodeSlice
	result, _ = node.SplitImageNodes()
	if len(result) != 4 || result[1].TextType != textTypeImage || result[3].TextType != textTypeImage {
		t.Fatalf("Incorrect output for SplitImageNodes():\nInput:\n%s\nResult:\n%s", node.Text, result.ToString())
	}
}

func TestLinkSplit(t *testing.T) {
	node := TextNode{
		TextType: textTypeText,
		Text: `Before link1 [link1 anchor text](http://link1.url) after link1,
		before link2 [link2 anchor text](http://link2.url)`,
	}
	var result TextNodeSlice
	result, _ = node.SplitLinkNodes()
	if len(result) != 4 || result[1].TextType != textTypeLink || result[3].TextType != textTypeLink {
		t.Fatalf("Incorrect output for SplitLinkNodes():\nInput:\n%s\nResult:\n%s", node.Text, result.ToString())
	}
}

// TODO: Figure out why this test appears to be non-deterministic
func TestSplitAll(t *testing.T) {
	const nText string = "**bold text** *italic text* `code text` " +
		"![image alt text](http://image.url) [link text](http://link.url)"

	var nodes TextNodeSlice
	nodes = append(nodes, TextNode{TextType: textTypeText, Text: nText})
	nodes, _ = nodes.SplitAll()

	// NOTE: Shouldn't length be 9 and nodes[7] be the single space between the image and the link?
	if len(nodes) != 8 ||
		nodes[0].TextType != textTypeBold ||
		nodes[2].TextType != textTypeItalic ||
		nodes[4].TextType != textTypeCode ||
		nodes[6].TextType != textTypeImage ||
		nodes[7].TextType != textTypeLink {
		t.Fatalf("Incorrect output for SplitAll():\nInput:\n%s\n%s", nText, nodes.ToString())
	}
}
