package main

import (
	"fmt"
)

func main() {
	txtNode := TextNode{
		Text:     "Lorem Ipsum",
		TextType: textTypeBold,
		URL:      "https://www.inerthtml.com",
	}

	hNode := HtmlNode{
		Tag:   "p",
		Value: "Lorem Ipsum",
		Props: map[string]string{
			"href":   "https://www.inerthtml.com",
			"target": "_blank",
		},
	}

	fmt.Printf("%#v\n", txtNode)
	fmt.Printf("%#v\n", hNode)
	fmt.Println(hNode.ToHTML())
}
