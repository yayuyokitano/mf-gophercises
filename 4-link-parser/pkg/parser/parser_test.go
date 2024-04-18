package parser_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"reflect"
	"testing"

	"golang.org/x/net/html"

	"github.com/yayuyokitano/mf-gophercises/4-link-parser/pkg/parser"
)

var testDir = path.Join("..", "..", "test")

func TestParseLinks(t *testing.T) {
	dir, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range dir {
		testFile(t, file)
	}
}

func testFile(t *testing.T, file os.DirEntry) {
	r, err := os.Open(path.Join(testDir, file.Name()))
	if err != nil {
		t.Errorf("open %s: %v", file.Name(), err)
		return
	}

	z := html.NewTokenizer(r)
	var expected []parser.Link
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			t.Errorf("%s: error token: %v", file.Name(), z.Err())
			return
		}
		if tt != html.StartTagToken {
			continue
		}
		name, hasAttr := z.TagName()
		if !bytes.Equal(name, []byte("body")) || !hasAttr {
			continue
		}

		key, val, _ := z.TagAttr()
		if !bytes.Equal(key, []byte("data-expected")) {
			t.Errorf("%s: expected %s, got %s", file.Name(), string(key), string(val))
			return
		}
		err := json.NewDecoder(bytes.NewReader(val)).Decode(&expected)
		if err != nil {
			t.Errorf("read test case %s: %v", file.Name(), err)
			return
		}
		break
	}

	r, err = os.Open(path.Join(testDir, file.Name()))
	if err != nil {
		t.Errorf("open %s: %v", file.Name(), err)
		return
	}
	actual, err := parser.ParseLinks(r)
	if err != nil {
		t.Errorf("parse %s: %v", file.Name(), err)
		return
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%s: expected %#v, got %#v", file.Name(), expected, actual)
	}
}
