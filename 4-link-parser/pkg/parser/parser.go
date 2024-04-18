package parser

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type Cursor struct {
	Pointer int
	Depth   int
}

func ParseLinks(r io.Reader) ([]Link, error) {
	var (
		z           = html.NewTokenizer(r)
		links       = make([]Link, 0)
		linkCursors = make([]Cursor, 0)
		depth       = 0
	)

	for {
		tt := z.Next()
		switch tt {

		case html.ErrorToken:
			if z.Err() == io.EOF {
				links = processLinkText(links)
				return links, nil
			}
			return links, z.Err()

		case html.TextToken:
			if len(linkCursors) > 0 {
				str := string(z.Text())
				for _, cur := range linkCursors {
					links[cur.Pointer].Text += str
				}
			}

		case html.StartTagToken, html.EndTagToken:
			name, hasAttr := z.TagName()
			if !bytes.Equal(name, []byte("a")) {
				continue
			}
			if !hasAttr {
				if tt == html.StartTagToken {
					depth++
				} else {
					depth--
					if len(linkCursors) > 0 && linkCursors[len(linkCursors)-1].Depth > depth {
						linkCursors = linkCursors[:len(linkCursors)-1]
					}
				}
				continue
			}

			var (
				key, val []byte
				moreAttr = true
			)
			for moreAttr {
				key, val, moreAttr = z.TagAttr()
				if !bytes.Equal(key, []byte("href")) {
					continue
				}

				if tt == html.StartTagToken {
					links = append(links, Link{
						Href: string(val),
						Text: "",
					})
					depth++
					linkCursors = append(linkCursors, Cursor{
						Pointer: len(links) - 1,
						Depth:   depth,
					})
				}
			}

		default:
			// Any other tag type is irrelevant to us, ignore it.

		}
	}
}

func processLinkText(links []Link) []Link {
	exp := regexp.MustCompile(`\s+`)
	for i := range links {
		s := html.UnescapeString(links[i].Text)
		s = strings.Join(strings.Fields(s), " ")
		s = strings.TrimSpace(s)
		links[i].Text = exp.ReplaceAllString(s, " ")
	}
	return links
}
