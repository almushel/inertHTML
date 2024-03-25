package main

import (
	"fmt"
)

func main() {
	txt := TextNode{
		Text:     "Lorem Ipsum",
		TextType: textTypeBold,
		URL:      "https://www.inerthtml.com",
	}

	html := HtmlNode{
		Tag:   "p",
		Value: "Lorem Ipsum",
		Props: map[string]string{
			"href":   "https://www.inerthtml.com",
			"target": "_blank",
		},
	}

	fmt.Printf("%#v\n", txt)
	fmt.Printf("%#v\n", html)
	fmt.Println(html.ToHTML())
}
