package main

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

func (node *TextNode) ToHTMLNode() (HtmlNode, error) {
	var result HtmlNode
	var err error

	switch node.TextType {
	case textTypeText:
		result = NewLeafNode("", node.Text)
		break
	case textTypeBold:
		result = NewLeafNode("strong", node.Text)
		break
	case textTypeItalic:
		result = NewLeafNode("em", node.Text)
		break
	case textTypeCode:
		result = NewLeafNode("code", node.Text)
		break
	case textTypeLink:
		result = NewLeafNode("a", node.Text)
		result.Props["href"] = node.URL
		break
	case textTypeImage:
		result = NewLeafNode("img", "")
		result.Props["src"] = node.URL
		result.Props["alt"] = node.Text
		break
	default:
		err = errors.New("TextNode.ToHtmlNode(): Invalid TextType")
	}

	return result, err
}

func (node *TextNode) Split(delim string, splitType int) ([]TextNode, error) {
	if node.TextType != textTypeText {
		return []TextNode{*node}, nil
	}

	var result []TextNode
	var err error

	// If node text starts with delim,
	// substrings with even indexes will be enclosed in the delim
	var evenType, oddType int
	if strings.HasPrefix(node.Text, delim) {
		evenType = splitType
		oddType = textTypeText
	} else {
		evenType = textTypeText
		oddType = splitType
	}

	chunks := strings.Split(node.Text, delim)
	if len(chunks)%2 != 0 {
		err = errors.New(fmt.Sprintf("TextNode.Split(): missing closing delimiter %s", delim))
	}

	for i, str := range chunks {
		if str == "" {
			continue
		}

		var chunkType int
		if i%2 == 0 {
			chunkType = evenType
		} else {
			chunkType = oddType
		}

		result = append(result, TextNode{
			TextType: chunkType,
			Text:     str,
		})
	}

	return result, err
}

func (nodeList TextNodeSlice) Split(delim string, splitType int) (TextNodeSlice, error) {
	var result []TextNode
	var err error

	var nResult []TextNode
	for _, n := range nodeList {
		nResult, err = n.Split(delim, splitType)
		result = append(result, nResult...)
	}

	return result, err
}

func (node *TextNode) SplitExp(pattern string, marshal func([]string) TextNode) ([]TextNode, error) {
	var result []TextNode
	expr, err := regexp.Compile(pattern)
	if err != nil {
		return result, err
	}

	textStrs := expr.Split(node.Text, -1)
	exprStrs := expr.FindAllStringSubmatch(node.Text, -1)
	var i, t int

	// If line starts with an image pattern, append that first
	if expr.FindStringIndex(node.Text)[0] == 0 {
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

func (node *TextNode) SplitImageNodes() ([]TextNode, error) {
	const pattern = "!\\[(.*?)\\]\\((.*?)\\)"
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

func (node *TextNode) SplitLinkNodes() ([]TextNode, error) {
	const pattern = "\\[(.*?)\\]\\((.*?)\\)"
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

	return node.SplitExp(pattern, marshal)
}
