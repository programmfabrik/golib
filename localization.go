package golib

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"

	"golang.org/x/text/language"
)

var re *regexp.Regexp

func init() {
	re = regexp.MustCompile(`\%\((.*?)\)([si])`)
}

type Localization struct {
	Keys        map[string]LocaValue
	defaultLang language.Tag
}

func NewLocalization() *Localization {
	loca := &Localization{}
	loca.Keys = map[string]LocaValue{}
	return loca
}

// Merge l2 into l1, existing keys are overwritten
func (l1 Localization) Merge(l2 *Localization) {
	if l2 == nil {
		return
	}
	for key, value := range l2.Keys {
		l1.Keys[key] = value
	}
}

// ToApi returns an api compatible loca map for the given language
func (l Localization) ToApi(lang string) map[string]map[string]string {
	apiLoca := make(map[string]string, len(l.Keys))
	for k, lV := range l.Keys {
		v, ok := lV[lang]
		if !ok || v == "" {
			continue
		}
		apiLoca[k] = v
	}

	return map[string]map[string]string{lang: apiLoca}
}

func (l *Localization) SetDefaultLang(lang language.Tag) error {
	l.defaultLang = lang
	return nil
}

func (loca Localization) Load(reader io.Reader) error {

	// read CSV file
	r := csv.NewReader(reader)

	var (
		keyCol int = -1
		row    int = -1
	)

	langByColIdx := make(map[int]string, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Localization.Load: failed: %w", err)
		}

		row = row + 1
		// fmt.Println("row:", row, record)

		// learn the order of the columns
		if row == 0 {
			for i, col := range record {
				if col == "key" {
					keyCol = i
				} else if col != "" {
					langByColIdx[i] = col
				}
			}
			if keyCol == -1 {
				return fmt.Errorf(`Localization.Load: Column "key" found in CSV data.`)
			}
			continue
		}

		locaValue := make(LocaValue, 0)

		for i, value := range record {
			lang, ok := langByColIdx[i]
			if ok {
				locaValue[lang] = value
			}
		}

		loca.Keys[record[keyCol]] = locaValue
	}

	// fmt.Printf("%+v", loca.Keys)
	return nil
}

func (l Localization) Get(locaKey string, langs []language.Tag) string {
	return l.GetWithParams(locaKey, langs, nil)
}

func (l Localization) GetLocaValue(key string) LocaValue {
	locaValue := l.Keys[key]
	if len(locaValue) == 0 {
		return LocaValue{l.defaultLang.String(): fmt.Sprintf("[!%s]", key)}
	}
	return locaValue
}

type LocaParams map[string]interface{}

func (l Localization) GetWithParams(key string, langs []language.Tag, info LocaParams) string {
	if len(langs) == 0 {
		langs = []language.Tag{l.defaultLang}
	}
	var text, value string
	locaValue := l.Keys[key]
	if len(locaValue) == 0 {
		return fmt.Sprintf("[!%s]", key)
	}
	for _, lang := range langs {
		value = locaValue[lang.String()]
		if value != "" {
			break
		}
	}
	if value == "" {
		return fmt.Sprintf("[!%s[%s]", key, l.defaultLang)
	}
	text = replaceText(value, info)
	return text
}

// GetLocaValue returns the loca value with all values replaced
// using the passed info.
func (l Localization) GetLocaValueWithParams(key string, info LocaParams) (lv LocaValue) {
	lv = LocaValue{}
	for lang, value := range l.GetLocaValue(key) {
		lv[lang] = replaceText(value, info)
	}
	return lv
}

// replaces the given string with info from the given
// map[string]interface{} (thats what pflib.H is)
func replaceText(str string, info map[string]interface{}) string {

	var params []interface{}

	str = re.ReplaceAllStringFunc(str, func(m string) string {
		parts := re.FindStringSubmatch(m)

		value, ok := info[parts[1]]
		if ok {
			params = append(params, value)
		} else {
			params = append(params, "!["+parts[1]+"]")
		}
		return "%v"
	})

	return fmt.Sprintf(str, params...)
}
