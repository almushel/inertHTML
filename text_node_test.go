package main

import (
	"fmt"
	"testing"
)

func TestTNSplit(t *testing.T) {
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
