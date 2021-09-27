package golib

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

// RenderTemplate takes bytes and executes it and returns the rendered bytes
func RenderTemplate(in []byte, funcMap template.FuncMap, data interface{}) (out []byte, err error) {
	tmpl, err := template.New("tmpl").Funcs(funcMap).Parse(string(in))
	if err != nil {
		return nil, errors.Wrap(err, "TmplToWriter[New] error")
	}

	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, errors.Wrap(err, "TmplToWriter[Execute] error")
	}
	return buf.Bytes(), nil
}
