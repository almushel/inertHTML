package main

import (
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

	blocks := strings.Split(md, "\n\n")
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
