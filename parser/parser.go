package parser

import (
	"strings"
)

func MDtoHTML(src, template string) (string, error) {
	var result string

	blocks := ParseMDBlocks(src)
	blockNodes, err := BlocksToHTMLNodes(blocks)
	if err != nil {
		return result, err
	}

	var pageTitle string

	for i := range blockNodes {
		blockNodes[i].ProcessInnerText()
		blockNodes[i].UnescapeMD()
	}

	var body string
	for _, node := range blockNodes {
		if pageTitle == "" && node.Tag == "h1" {
			pageTitle = node.Value
		}
		body += node.ToHTML()
	}

	result = strings.Join(
		strings.Split(template, "{{ Title }}"),
		pageTitle,
	)
	result = strings.Join(
		strings.Split(result, "{{ Content }}"),
		body,
	)

	return result, nil
}
