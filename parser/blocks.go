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
	blockTypeTable
)

func IsNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

func ParseMDBlocks(md string) []string {
	var result []string

	// NOTE: Dealing with code blocks separately, because they are allowed to break whitespace rules
	var start, end int
	for start = 0; start < len(md); start++ {
		if strings.HasPrefix(md[start:], "```") {
			for inside := start + 3; inside < len(md); inside++ {
				if strings.HasPrefix(md[inside:], "\n```") {
					result = append(result, md[start:inside+4])

					end = inside + 4
					start = end

					break
				}
			}
		} else if strings.HasPrefix(md[start:], "\r\n\r\n") {
			block := strings.TrimSpace(md[end:start])
			if len(block) > 0 {
				result = append(result, block)
			}

			end = start
			start += 3
		} else if strings.HasPrefix(md[start:], "\n\n") {
			block := strings.TrimSpace(md[end:start])

			if len(block) > 0 {
				result = append(result, block)
			}

			end = start
			start += 1
		}
	}

	if end < len(md) {
		block := strings.TrimSpace(md[end:start])
		if len(block) > 0 {
			result = append(result, block)
		}
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
			if !(strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "- ")) {
				valid = false
				break
			}
		}
		if valid {
			return blockTypeUnorderedList
		}
	} else if block[0] == '|' {
		var bars int
		for i, line := range strings.Split(strings.TrimSpace(block), "\n") {
			if i == 0 { // | title1 | title2 |
				if line[0] != '|' || line[len(line)-1] != '|' {
					return blockTypeParagraph
				}

				bars = strings.Count(line, "|")
			} else if i == 1 { // | :--- | ---: |
				if line[0] != '|' ||
					line[len(line)-1] != '|' ||
					strings.Count(line, "|") != bars {
					return blockTypeParagraph
				}

				for _, cell := range strings.Split(strings.Trim(line, "|"), "|") {
					cell = strings.TrimSpace(cell)
					if strings.HasPrefix(cell, ":-") {
						cell = cell[1:]
					}
					if strings.HasSuffix(cell, "-:") {
						cell = cell[:len(cell)-1]
					}

					if strings.Count(cell, "-") < 3 ||
						len(strings.Trim(cell, "-")) > 0 {
						return blockTypeParagraph
					}
				}
			} else { // | content | content |
				if strings.Count(line, "|") != bars {
					return blockTypeParagraph
				}
			}
		}

		return blockTypeTable
	} else if block == "***" || block == "---" || block == "___" {
		return blockTypeHorizontalRule
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
		if strings.HasSuffix(block, "\n```") {
			return blockTypeCode
		}
	}

	return blockTypeParagraph
}

func generateValidId(value string) string {
	remove := []string{
		"`", "[", "]", "(", ")", ":", ";", ".", "?",
		"=", "+", "%", "^", "$", "#", "@", "!",
		"*", "~", "{", "}", "^", "<", ">",
	}
	replace := []string{
		" ", "-",
	}
	for _, r := range remove {
		replace = append(replace, r, "")
	}

	replacer := strings.NewReplacer(replace...)

	return strings.ToLower(
		html.EscapeString(
			replacer.Replace(value),
		),
	)
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

			newNode.Props = map[string]string{"id": generateValidId(newNode.Value)}
			break

		case blockTypeCode:
			opening, body, _ := strings.Cut(block, "\n")
			lang, name, _ := strings.Cut(opening[len("```"):], " ")
			newNode = HtmlNode{
				Tag: "pre",
				Value: html.EscapeString(
					strings.TrimSpace(
						body[:len(body)-len("```")],
					),
				),
			}
			newNode.Props = make(map[string]string)
			if len(lang) > 0 {
				newNode.Props["class"] = "language-" + lang
			}

			if len(name) > 0 {
				newNode.Props["title"] = name
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
			newNode = NewHtmlNode("ol", "", nil, nil)
			for lineNumber, line := range strings.Split(block, "\n") {
				var i int
				var c rune
				for i, c = range line {
					if !IsNumeric(c) {
						break
					}
				}
				if lineNumber == 0 {
					newNode.Props["start"] = line[:i]
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

		case blockTypeTable:
			newNode = HtmlNode{
				Tag:   "div",
				Props: map[string]string{"style": "overflow-x:auto;"},
			}
			table := HtmlNode{
				Tag: "table",
			}

			lines := strings.Split(block, "\n")
			divider := strings.Split(strings.Trim(lines[1], "|"), "|")
			var alignments []string

			for _, cell := range divider {
				c := strings.TrimSpace(cell)
				left := strings.HasPrefix(c, ":")
				right := strings.HasSuffix(c, ":")

				if left && right {
					alignments = append(alignments, "center")
				} else if right && !left {
					alignments = append(alignments, "right")
				} else {
					alignments = append(alignments, "left")
				}
			}

			head := HtmlNode{Tag: "thead"}
			body := HtmlNode{Tag: "tbody"}

			var tag string
			for i, line := range lines {
				if i == 1 {
					continue
				} else if i == 0 {
					tag = "th"
				} else {
					tag = "td"
				}

				row := HtmlNode{
					Tag: "tr",
				}
				cells := strings.Split(strings.Trim(line, "|"), "|")
				for i, cell := range cells {
					row.Children = append(row.Children, HtmlNode{
						Tag:   tag,
						Value: strings.TrimSpace(cell),
						Props: map[string]string{
							"style": "text-align: " + alignments[i],
						},
					})
				}

				if i == 0 {
					head.Children = append(head.Children, row)
				} else {
					body.Children = append(body.Children, row)
				}
			}

			table.Children = append(table.Children, head, body)
			newNode.Children = append(newNode.Children, table)

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
