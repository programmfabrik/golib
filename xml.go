package golib

import (
	"encoding/xml"
	"strconv"

	"github.com/antchfx/xmlquery"
)

func AttrByName(attrs []xmlquery.Attr, name string) string {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

func XmlStart(name string) xml.StartElement {
	return xml.StartElement{Name: xml.Name{Local: name}}
}

func XmlEnd(name string) xml.EndElement {
	return xml.EndElement{Name: xml.Name{Local: name}}
}

func XmlSetAttr(start *xml.StartElement, attr xml.Attr) {
	start.Attr = append(start.Attr, attr)
}

func XmlAttr(name, value string) xml.Attr {
	return xml.Attr{Name: xml.Name{Local: name}, Value: value}
}

func XmlAttrAppend(start *xml.StartElement, name, value string) {
	start.Attr = append(start.Attr, XmlAttr(name, value))
}

func XmlAttrInt64(name string, value int64) xml.Attr {
	return XmlAttr(name, strconv.FormatInt(value, 10))
}

func XmlAttrInt64Append(start *xml.StartElement, name string, value int64) {
	start.Attr = append(start.Attr, XmlAttrInt64(name, value))
}

func XmlAttrInt(name string, value int) xml.Attr {
	return XmlAttr(name, strconv.Itoa(value))
}

func XmlAttrIntAppend(start *xml.StartElement, name string, value int) {
	start.Attr = append(start.Attr, XmlAttrInt(name, value))
}
