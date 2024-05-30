package golib

import (
	"bytes"
	"fmt"
	"text/template"
)

// RenderTemplate takes bytes and executes it and returns the rendered bytes
func RenderTemplate(in []byte, funcMap template.FuncMap, data interface{}) (out []byte, err error) {
	tmpl, err := template.New("tmpl").Funcs(funcMap).Parse(string(in))
	if err != nil {
		return nil, fmt.Errorf("TmplToWriter[New] error: %w", err)
	}

	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("TmplToWriter[Execute] error: %w", err)
	}
	return buf.Bytes(), nil
}
