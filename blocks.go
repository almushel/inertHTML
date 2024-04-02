package main

import (
	"fmt"
	"strings"
)

const (
	blockTypeParagraph = iota
	blockTypeHeading
	blockTypeCode
	blockTypeQuote
	blockTypeOrderedList
	blockTypeUnorderedList
)

func IsNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

func ParseMDBlocks(md string) []string {
	var result []string

	var sep string
	// Attempting to account for Windows newline format
	if strings.Count(md, "\n\n") > strings.Count(md, "\r\n\r\n") {
		sep = "\n\n"
	} else {
		sep = "\r\n\r\n"
	}

	blocks := strings.Split(md, sep)
	for _, b := range blocks {
		result = append(result, strings.TrimSpace(b))
	}

	return result
}

func GetBlockType(block string) int {
	if block == "" {
		return blockTypeParagraph
	}
	if block[0] == '#' {
		// Heading
		for level, char := range block {
			if char != '#' {
				if level < 6 && block[level] == ' ' {
					return blockTypeHeading
				} else {
					break
				}
			}
		}
	} else if block[0] == '>' {
		// Blockquote
		var valid bool = true
		for _, line := range strings.Split(block, "\n") {
			if line[0] != '>' {
				valid = false
				break
			}
		}
		if valid {
			return blockTypeQuote
		}
	} else if block[0] == '-' || block[0] == '*' {
		// Unordered List
		var valid bool = true
		for _, line := range strings.Split(block, "\n") {
			if !(strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "- ")) {
				valid = false
				break
			}
		}
		if valid {
			return blockTypeUnorderedList
		}
	} else if IsNumeric([]rune(block)[0]) {
		// Ordered List
		var valid bool = true
		for _, line := range strings.Split(block, "\n") {
			for i, c := range line {
				if !IsNumeric(c) {
					valid = (c == '.' && line[i+1] == ' ')
					break
				}
			}
			if !valid {
				break
			}
		}
		if valid {
			return blockTypeOrderedList
		}
	} else if strings.HasPrefix(block, "```") {
		// Code Block
		if strings.HasSuffix(block, "```") {
			return blockTypeCode
		}
	}

	return blockTypeParagraph
}

func BlocksToHTMLNodes(blocks []string) ([]HtmlNode, error) {
	var result []HtmlNode
	var err error
	var newNode HtmlNode
	for _, block := range blocks {
		switch GetBlockType(block) {

		case blockTypeHeading:
			var i int = 0
			for range block {
				if block[i] != '#' {
					break
				}
				i++
			}

			newNode = HtmlNode{
				Tag:   "h" + fmt.Sprintf("%v", i),
				Value: block[i+1:],
			}
			break

		case blockTypeCode:
			newNode = HtmlNode{
				Tag:   "code",
				Value: block[3+1 : len(block)-3],
			}
			break

		case blockTypeQuote:
			var quoteText string
			for _, line := range strings.Split(block, "\n") {
				quoteText += line[3:] + "\n"
			}
			newNode = HtmlNode{
				Tag:   "blockquote",
				Value: quoteText[:len(quoteText)-1],
			}
			break
		case blockTypeOrderedList:
			newNode = HtmlNode{
				Tag: "ol",
			}
			for _, line := range strings.Split(block, "\n") {
				newNode.Children = append(newNode.Children, HtmlNode{
					Tag:   "li",
					Value: line[3:],
				})
			}
			break
		case blockTypeUnorderedList:
			newNode = HtmlNode{
				Tag: "ul",
			}
			for _, line := range strings.Split(block, "\n") {
				newNode.Children = append(newNode.Children, HtmlNode{
					Tag:   "li",
					Value: line[3:],
				})
			}
			break
		default:
			newNode = HtmlNode{
				Tag:   "p",
				Value: block,
			}
			break
		}

		result = append(result, newNode)
	}

	return result, err
}
