package main

import (
	"fmt"
)

type TextNode struct {
	Text, TextType, URL string
}

type HtmlNode struct {
	Tag, Value string
	Children   []HtmlNode
	Props      map[string]string
}

func (node *HtmlNode) ToHTML() string {
	var result string

	var children string
	for _, child := range node.Children {
		children += child.ToHTML() + "\n"
	}

	result = fmt.Sprintf("<%s %s>%s%s</%s>", node.Tag, node.PropsToHTML(), node.Value, children, node.Tag)

	return result
}

func (node *HtmlNode) PropsToHTML() string {
	var result string
	for key, val := range node.Props {
		result += fmt.Sprintf(" %s=\"%s\"", key, val)
	}

	return result
}
