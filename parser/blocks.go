package parser

import (
	"fmt"
	"html"
	"strings"
)

const (
	blockTypeParagraph = iota
	blockTypeCode
	blockTypeHeading
	blockTypeHorizontalRule
	blockTypeOrderedList
	blockTypeUnorderedList
	blockTypeQuote
)

func IsNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

func ParseMDBlocks(md string) []string {
	var result []string

	// NOTE: Dealing with code blocks separately, because they are allowed to break whitespace rules
	var start, end int
	for start = 0; start < len(md); start++ {
		if strings.HasPrefix(md[start:], "```\n") {
			for inside := start + 3; inside < len(md); inside++ {
				if strings.HasPrefix(md[inside:], "\n```") {
					result = append(result, md[start:inside+4])

					end = inside + 4
					start = end

					break
				}
			}
		} else if strings.HasPrefix(md[start:], "\r\n\r\n") {
			block := strings.TrimSpace(strings.TrimSpace(md[end:start]))
			if len(block) > 0 {
				result = append(result, block)
			}

			end = start
			start += 3
		} else if strings.HasPrefix(md[start:], "\n\n") {
			block := strings.TrimSpace(strings.TrimSpace(md[end:start]))

			if len(block) > 0 {
				result = append(result, block)
			}

			end = start
			start += 1
		}
	}

	block := strings.TrimSpace(md[end:start])
	if len(block) > 0 {
		result = append(result, block)
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
		/* TODO: Nested lists */
	} else if (block[0] == '-' || block[0] == '*') && block[1] == ' ' {
		// Unordered List
		var valid bool = true
		for _, line := range strings.Split(block, "\n") {
			if !(strings.HasPrefix(line, "* ") || !strings.HasPrefix(line, "- ")) {
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
	} else if strings.HasPrefix(block, "```\n") {
		// Code Block
		if strings.HasSuffix(block, "\n```") {
			return blockTypeCode
		}
	} else if block == "***" || block == "---" || block == "___" {
		return blockTypeHorizontalRule
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
				Tag:   "pre",
				Value: html.EscapeString(block[3+1 : len(block)-4]),
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
				var i int
				var c rune
				for i, c = range line {
					if !IsNumeric(c) {
						break
					}
				}
				newNode.Children = append(newNode.Children, HtmlNode{
					Tag:   "li",
					Value: line[i+2:],
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
					Value: line[2:],
				})
			}
			break
		case blockTypeHorizontalRule:
			newNode = HtmlNode{
				Tag: "hr",
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
