package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	HREF string
	TEXT string
}

func parseLinks(node *html.Node) []Link {
	var links []Link

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				text := extractText(node)
				links = append(links, Link{HREF: attr.Val, TEXT: text})
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		exLinks := parseLinks(child)
		links = append(links, exLinks...)
	}
	return links

}

func extractText(node *html.Node) string {
	var text string

	if node.Type != html.ElementNode && node.Type != html.CommentNode && node.Data != "a" {
		text = node.Data
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		text += extractText(child)
	}

	return strings.Trim(text, "\n")
}

func main() {
	HTMLbytes, err := os.Open("doc.html")
	checkError(err)

	doc, err := html.Parse(HTMLbytes)
	checkError(err)

	listLinks := parseLinks(doc)

	fmt.Println(listLinks)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
