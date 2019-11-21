package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func check(r []Link, s []Link, t *testing.T) {

	if len(r) != len(s) {
		t.Errorf("Error yaw, expected %v, got %v", s, r)
	}

	for i, v := range r {
		if v != s[i] {
			t.Errorf("Error yaw, expected %v, got %v", s, r)
			break
		}
	}
}

func TestHtmlParsing(t *testing.T) {

	t.Run("parse an html without anchor tag", func(t *testing.T) {
		str := `<html>
			<body>

			<p>This is a paragraph.</p>
			<p>This is a paragraph.</p>
			<p>This is a paragraph.</p>

			</body>
			</html>
			`

		doc := strings.NewReader(str)
		node, _ := html.Parse(doc)

		result := parseLinks(node)
		slice := []Link{}

		check(result, slice, t)
	})

	t.Run("parse an html with anchor tag", func(t *testing.T) {
		str := `<html>
		<body>

		<a href="/webpage">Some Text!</a>
		<a href="/otherpage">Other Text!</a>

		</body>
		</html>
		`

		doc := strings.NewReader(str)
		node, _ := html.Parse(doc)

		result := parseLinks(node)
		slice := []Link{
			Link{
				HREF: "/webpage",
				TEXT: "Some Text!",
			},
			Link{
				HREF: "/otherpage",
				TEXT: "Other Text!",
			},
		}

		check(result, slice, t)
	})

	t.Run("parse an html with anchor tag and span for text", func(t *testing.T) {
		str := `<html>
		<body>

		<a href="/webpage"><span>some span</span></a>
		<a href="/otherpage">Other Text!<span>some span</span></a>

		</body>
		</html>
		`

		doc := strings.NewReader(str)
		node, _ := html.Parse(doc)

		result := parseLinks(node)
		slice := []Link{
			Link{
				HREF: "/webpage",
				TEXT: "some span",
			},
			Link{
				HREF: "/otherpage",
				TEXT: "Other Text!some span",
			},
		}

		check(result, slice, t)
	})

	t.Run("parse an html with anchor tag and comments", func(t *testing.T) {
		str := `<html>
		<body>

		<a href="/webpage">Text!--comment</a>
		//some comment

		</body>
		</html>
		`
	})

}
