package note

import (
	"bytes"
	"html/template"
	"os"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Metadata map[string]interface{}

type Note struct {
	Metadata
	Content template.HTML
}

// Setup goldmark
var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Footnote,
		extension.Typographer, // -- into –, ... into … etc.
		meta.Meta,             // YAML metadata header parsing
		mathjax.MathJax,
		emoji.Emoji, // Github Emoji set
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

// New reads the file at the given path, extracts the metadata and converts the
// content to HTML.
func New(path string) (*Note, error) {
	note := new(Note)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var html bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(content, &html, parser.WithContext(context)); err != nil {
		return nil, err
	}
	metadata := meta.Get(context)

	note.Metadata = metadata
	note.Content = template.HTML(html.String())

	return note, nil
}
