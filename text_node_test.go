package main

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
	result, _ := node.SplitImageNodes()
	if len(result) != 4 && result[1].TextType != textTypeImage && result[2].TextType != textTypeImage {
		t.Fatalf("Incorrect output for SplitImageNodes():\nInput:%s\nResult:%#v", node.Text, result)
	}
}

func TestLinkSplit(t *testing.T) {
	node := TextNode{
		TextType: textTypeText,
		Text: `Before link1 ![link1 anchor text](http://link1.url) after link1,
		before link2 ![link2 anchor text](http://link2.url)`,
	}
	result, _ := node.SplitLinkNodes()
	if len(result) != 4 && result[1].TextType != textTypeLink && result[2].TextType != textTypeLink {
		t.Fatalf("Incorrect output for SplitImageNodes():\nInput:%s\nResult:%#v", node.Text, result)
	}
}
