package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	textTypeUndefined = iota
	textTypeText
	textTypeBold
	textTypeItalic
	textTypeCode
	textTypeLink
	textTypeImage
)

type TextNode struct {
	TextType  int
	Text, URL string
}
type TextNodeSlice []TextNode
type TextNodeSplitFunc func(*TextNode) ([]TextNode, error)

func (node TextNode) ToHTMLNode() (HtmlNode, error) {
	var result HtmlNode
	var err error

	switch node.TextType {
	case textTypeText:
		result = NewHtmlNode("", node.Text, nil, nil)
		break
	case textTypeBold:
		result = NewHtmlNode("strong", node.Text, nil, nil)
		break
	case textTypeItalic:
		result = NewHtmlNode("em", node.Text, nil, nil)
		break
	case textTypeCode:
		result = NewHtmlNode("code", node.Text, nil, nil)
		break
	case textTypeLink:
		result = NewHtmlNode("a", node.Text, nil,
			map[string]string{
				"href": node.URL,
			},
		)
		break
	case textTypeImage:
		result = NewHtmlNode("img", "", nil,
			map[string]string{
				"src": node.URL,
				"alt": node.Text,
			},
		)
		break
	default:
		err = errors.New("TextNode.ToHtmlNode(): Invalid TextType")
	}

	return result, err
}

func (node TextNode) Split(delim string, splitType int) ([]TextNode, error) {
	if node.TextType != textTypeText {
		return []TextNode{node}, nil
	}

	var result []TextNode
	var err error

	var indices []int
	for i := 0; i < len(node.Text); i++ {
		if strings.HasPrefix(node.Text[i:], delim) {
			if i > 0 && node.Text[i-1] == '\\' {
				if !(i > 1 && node.Text[i-2] == '\\') {
					continue
				}
			}

			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		return []TextNode{node}, nil
	}

	if indices[0] != 0 {
		result = append(result, TextNode{
			TextType: textTypeText,
			Text:     node.Text[:indices[0]],
		})
	}

	for i := range indices {
		var chunkType int
		if i%2 == 0 {
			chunkType = splitType
		} else {
			chunkType = textTypeText
		}

		var splitText string
		if i == len(indices)-1 {
			splitText = node.Text[indices[i]+len(delim):]
		} else {
			splitText = node.Text[indices[i]+len(delim) : indices[i+1]]
		}

		if len(splitText) > 0 {
			result = append(result, TextNode{
				TextType: chunkType,
				Text:     splitText,
			})
		}
	}

	return result, err
}

func (node TextNode) SplitExp(pattern string, marshal func([]string) TextNode) ([]TextNode, error) {
	if node.TextType != textTypeText {
		return []TextNode{node}, nil
	}

	expr, err := regexp.Compile(pattern)
	if err != nil {
		return []TextNode{node}, err
	}

	var result []TextNode
	var textStrs []string = expr.Split(node.Text, -1)
	var exprStrs [][]string = expr.FindAllStringSubmatch(node.Text, -1)
	var i, t int

	// If line starts with the pattern, append that first
	if exprStrs != nil && expr.FindStringIndex(node.Text)[0] == 0 {
		exprNode := marshal(exprStrs[i])
		result = append(result, exprNode)
		i++
	}

	for t < len(textStrs) || i < len(exprStrs) {
		if t < len(textStrs) {
			if textStrs[t] != "" {
				textNode := TextNode{
					TextType: textTypeText,
					Text:     textStrs[t],
				}
				result = append(result, textNode)
			}
			t++
		}

		if i < len(exprStrs) {
			exprNode := marshal(exprStrs[i])
			result = append(result, exprNode)
			i++
		}
	}

	return result, err
}

func (node TextNode) SplitImageNodes() ([]TextNode, error) {
	const pattern = `!\[(.*?)\]\((.*?)\)`
	marshal := func(match []string) TextNode {
		var result TextNode
		if len(match) == 3 {
			result = TextNode{
				TextType: textTypeImage,
				Text:     match[1],
				URL:      match[2],
			}
		}
		return result
	}

	return node.SplitExp(pattern, marshal)
}

// NOTE: This will also match images, and so should be run after SplitImageNodes
func (node TextNode) SplitLinkNodes() ([]TextNode, error) {
	const link = `\[(.*?)\]\((.*?)\)`

	marshal := func(match []string) TextNode {
		var result TextNode
		if len(match) == 3 {
			result = TextNode{
				TextType: textTypeLink,
				Text:     match[1],
				URL:      match[2],
			}
		}
		return result
	}

	return node.SplitExp(link, marshal)
}

func (nodeList TextNodeSlice) ForEach(f func(TextNode)) {
	for _, node := range nodeList {
		f(node)
	}
}

func (nodeList TextNodeSlice) ToString() string {
	var result string
	nodeList.ForEach(func(n TextNode) {
		result += fmt.Sprintf("%#v\n", n)
	})
	return result
}

func (nodeList TextNodeSlice) SplitFunc(split TextNodeSplitFunc) (TextNodeSlice, error) {
	var result []TextNode
	var err error

	var nResult []TextNode
	for _, n := range nodeList {
		nResult, err = split(&n)
		/*
			if err != nil {
				break
			}
		*/
		result = append(result, nResult...)
	}

	return result, err
}

func (nodeList TextNodeSlice) Split(delim string, splitType int) ([]TextNode, error) {
	return nodeList.SplitFunc(
		func(n *TextNode) ([]TextNode, error) {
			return n.Split(delim, splitType)
		},
	)
}

func (nodeList TextNodeSlice) SplitLinkNodes() ([]TextNode, error) {
	return nodeList.SplitFunc(
		func(n *TextNode) ([]TextNode, error) {
			return n.SplitLinkNodes()
		},
	)
}

func (nodeList TextNodeSlice) SplitImageNodes() ([]TextNode, error) {
	return nodeList.SplitFunc(
		func(n *TextNode) ([]TextNode, error) {
			return n.SplitImageNodes()
		},
	)
}

func (nodeList TextNodeSlice) SplitAll() ([]TextNode, error) {
	type textDelim struct {
		d string
		t int
	}

	delims := []textDelim{
		{"`", textTypeCode},
		{"**", textTypeBold},
		{"__", textTypeBold},
		{"*", textTypeItalic},
		{"_", textTypeItalic},
	}

	var result TextNodeSlice
	var err error

	result, err = nodeList.SplitImageNodes()
	result, err = result.SplitLinkNodes()
	for _, delim := range delims {
		result, err = result.Split(delim.d, delim.t)
	}

	return result, err
}
