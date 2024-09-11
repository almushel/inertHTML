package parser

import (
	"fmt"
	"slices"
	"strings"
)

type HtmlNode struct {
	Tag, Value string
	Children   []HtmlNode
	Props      map[string]string
}

func (node *HtmlNode) ProcessInnerText() {
	if node.Tag == "pre" || node.Tag == "code" {
		return
	}

	if node.Value != "" {
		var innerTextNodes TextNodeSlice
		innerTextNodes = append(innerTextNodes, TextNode{
			TextType: textTypeText,
			Text:     node.Value,
		})

		innerTextNodes, _ = innerTextNodes.SplitAll()
		if len(innerTextNodes) > 1 || innerTextNodes[0].TextType != textTypeText {
			for _, itNode := range innerTextNodes {
				child, _ := itNode.ToHTMLNode()
				node.Children = append(node.Children, child)
			}

			node.Value = ""
		} else {
			node.Value = innerTextNodes[0].Text
		}
	}

	if len(node.Children) > 0 {
		for i := range node.Children {
			node.Children[i].ProcessInnerText()
		}
	}
}

func (node *HtmlNode) UnescapeMD() {
	replacer := strings.NewReplacer(
		"\\*", "*",
		"\\_", "_",
		"\\`", "`",
		"\\\\", "\\",
	)

	node.Value = replacer.Replace(node.Value)

	for i := range node.Children {
		node.Children[i].UnescapeMD()
	}
}

func (node *HtmlNode) ToHTML() string {
	var result string

	if node.Tag == "" {
		result = node.Value
	} else {
		var children string = node.Value
		for _, child := range node.Children {
			children += child.ToHTML()
		}

		result = fmt.Sprintf("<%s%s>%s</%s>", node.Tag, node.PropsToHTML(), children, node.Tag)
	}

	return result
}

func (node *HtmlNode) PropsToHTML() string {
	keys := make([]string, len(node.Props))

	var i int
	for k := range node.Props {
		keys[i] = k
		i++
	}
	slices.Sort(keys)

	var result string
	for _, key := range keys {
		val := node.Props[key]
		result += fmt.Sprintf(" %s=\"%s\"", key, val)
	}

	return result
}

func NewHtmlNode(tag, value string, children []HtmlNode, props map[string]string) HtmlNode {
	result := HtmlNode{
		Tag:      tag,
		Value:    value,
		Children: children,
		Props:    props,
	}
	if props == nil {
		result.Props = make(map[string]string)
	}

	return result
}
