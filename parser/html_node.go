package parser

import (
	"fmt"
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

func (node *HtmlNode) ToHTML() string {
	var result string

	if node.Tag == "" {
		result = node.Value
	} else {
		result = fmt.Sprintf("<%s%s>\n", node.Tag, node.PropsToHTML())
		if len(node.Value) > 0 {
			result += fmt.Sprintf("%s\n", node.Value)
		}
		for _, child := range node.Children {
			result += child.ToHTML()
		}
		result += fmt.Sprintf("</%s>\n", node.Tag)
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
