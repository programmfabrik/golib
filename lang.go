package lib

import "golang.org/x/text/language"

type Language string
type Languages []Language

func (lg Language) Tag() language.Tag {
	return language.Make(string(lg))
}

func (lgs Languages) Tags() []language.Tag {
	t := make([]language.Tag, len(lgs))
	for idx, lg := range lgs {
		t[idx] = lg.Tag()
	}
	return t
}

func (lgs Languages) Strings() []string {
	s := make([]string, len(lgs))
	for idx, lg := range lgs {
		s[idx] = string(lg)
	}
	return s
}

func (lgs Languages) Contains(lang string) bool {
	for _, lg := range lgs {
		if string(lg) == lang {
			return true
		}
	}
	return false
}

func NewLanguage(tag language.Tag) Language {
	return Language(tag.String())
}

func NewLanguages(tags []language.Tag) Languages {
	lgs := make(Languages, len(tags))
	for idx, tag := range tags {
		lgs[idx] = NewLanguage(tag)
	}
	return lgs
}
