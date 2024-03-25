package main

import (
	"errors"
	"fmt"
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
		err = errors.New("TextNode.ToHtmlNode: Invalid TextType")
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
