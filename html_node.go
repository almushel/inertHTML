package main

import (
	"fmt"
)

type HtmlNode struct {
	Tag, Value string
	Children   []HtmlNode
	Props      map[string]string
}

func (node *HtmlNode) ToHTML() string {
	var result string

	if node.Tag == "" {
		result = node.Value
	} else {
		var children string
		for _, child := range node.Children {
			children += child.ToHTML() + "\n"
		}

		result = fmt.Sprintf("<%s%s>%s%s</%s>", node.Tag, node.PropsToHTML(), node.Value, children, node.Tag)
	}

	return result
}

func (node *HtmlNode) PropsToHTML() string {
	var result string
	for key, val := range node.Props {
		result += fmt.Sprintf(" %s=\"%s\"", key, val)
	}

	return result
}

func NewHTMLNode(tag, value string, children []HtmlNode, props map[string]string) HtmlNode {
	result := HtmlNode{
		Tag:      tag,
		Value:    value,
		Children: children,
		Props:    props,
	}

	return result
}

func NewLeafNode(tag, value string) HtmlNode {
	result := HtmlNode{
		Tag:   tag,
		Value: value,
	}

	return result
}

func NewParentNode(tag string, children []HtmlNode) HtmlNode {
	result := HtmlNode{
		Tag:      tag,
		Children: children,
	}

	return result
}
