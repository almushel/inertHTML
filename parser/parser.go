package parser

type InertParserResult struct {
	Title string
	Body  string
}

func MDtoHTML(src string) (InertParserResult, error) {
	var result InertParserResult

	blocks := ParseMDBlocks(src)
	blockNodes, err := BlocksToHTMLNodes(blocks)
	if err != nil {
		return result, err
	}

	for i := range blockNodes {
		blockNodes[i].ProcessInnerText()
		blockNodes[i].UnescapeMD()
	}

	for _, node := range blockNodes {
		if result.Title == "" && node.Tag == "h1" {
			result.Title = node.Value
		}
		result.Body += node.ToHTML()
	}

	return result, nil
}
