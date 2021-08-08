package lib

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

type LocaValue map[string]string

// Trim removes all whitespace from set values
// and remove empty values entirely
func (lv LocaValue) TrimSpace() (lv2 LocaValue) {
	lv2 = LocaValue{}
	for lang, value := range lv {
		value = strings.TrimSpace(value)
		if value != "" {
			lv2[lang] = value
		}
	}
	return lv2
}

// LangShort cuts the lang returns the first part
//
func LangShort(lang string) string {
	if lang == "" {
		return lang
	}
	parts := strings.Split(lang, "-")
	return parts[0]
}

// First returns first set lang, set in langs
func (lv LocaValue) Best(prefLangs ...language.Tag) string {
	if lv == nil {
		return ""
	}
	for _, lang := range prefLangs {
		v := lv[lang.String()]
		if v != "" {
			return v
		}
	}
	// Return first lang which is set
	for _, v := range lv {
		if v != "" {
			return v
		}
	}
	return ""
}

func (lv LocaValue) Value() (driver.Value, error) {
	// if len(lv) == 0 {
	// 	return nil, nil
	// }
	return json.Marshal(lv)
}

func (lv LocaValue) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	enc.EncodeToken(start)
	for key, value := range lv {
		enc.EncodeElement(value, XmlStart(key))
	}
	enc.EncodeToken(start.End())
	return nil
}

func (lv *LocaValue) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	switch v2 := v.(type) {
	case []uint8:
		return json.Unmarshal(v2, lv)
	case string:
		return json.Unmarshal([]uint8(v2), lv)
	default:
		panic(fmt.Sprintf("LocaValuenc.Scan: Unsupported type for scan: %T", v))
	}
}

func (lv LocaValue) Copy() LocaValue {
	if lv == nil {
		return nil
	}
	lv2 := LocaValue{}
	for k, v := range lv {
		lv2[k] = v
	}
	return lv2
}

// Fill returns the LocaValue with all languages filled out
// If a value is missing, "Best" is used to fill it
func (lv LocaValue) Fill(prefLangs ...language.Tag) (lv2 LocaValue) {
	lv2 = LocaValue{}
	for _, prefLang := range prefLangs {
		pl := prefLang.String()
		if lv[pl] != "" {
			lv2[pl] = lv[pl]
		} else {
			lv2[pl] = lv2.Best(prefLangs...)
		}
	}
	return lv2
}
