package golib

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func RenderMarkdown(md []byte, gm goldmark.Markdown) (bs []byte, err error) {
	buf := bytes.Buffer{}
	err = gm.Convert(md, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MarkdownDefault() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
}
