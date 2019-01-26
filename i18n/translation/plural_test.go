package translation

import (
	"reflect"
	"testing"

	"github.com/go-flow/flow/i18n/language"
)

func mustTemplate(t *testing.T, src string) *template {
	tmpl, err := newTemplate(src)
	if err != nil {
		t.Fatal(err)
	}
	return tmpl
}

func pluralTranslationFixture(t *testing.T, id string, pluralCategories ...language.Plural) *pluralTranslation {
	templates := make(map[language.Plural]*template, len(pluralCategories))
	for _, pc := range pluralCategories {
		templates[pc] = mustTemplate(t, string(pc))
	}
	return &pluralTranslation{id, templates}
}

func verifyDeepEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\n%#v\nnot equal to expected value\n%#v", actual, expected)
	}
}

func TestPluralTranslationMerge(t *testing.T) {
	pt := pluralTranslationFixture(t, "id", language.One, language.Other)
	oneTemplate, otherTemplate := pt.templates[language.One], pt.templates[language.Other]

	pt.Merge(pluralTranslationFixture(t, "id"))
	verifyDeepEqual(t, pt.templates, map[language.Plural]*template{
		language.One:   oneTemplate,
		language.Other: otherTemplate,
	})

	pt2 := pluralTranslationFixture(t, "id", language.One, language.Two)
	pt.Merge(pt2)
	verifyDeepEqual(t, pt.templates, map[language.Plural]*template{
		language.One:   pt2.templates[language.One],
		language.Two:   pt2.templates[language.Two],
		language.Other: otherTemplate,
	})
}
