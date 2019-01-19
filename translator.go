package flow

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-flow/flow/i18n"
	"github.com/go-flow/flow/i18n/language"
	"github.com/go-flow/flow/i18n/translation"
)

// LanguageExtractor can be implemented for custom finding of search
// languages. This can be useful if you want to load a user's language
// from something like a database.
type LanguageExtractor func(LanguageExtractorOptions, *Context) []string

// LanguageExtractorOptions is a map of options for a LanguageExtractor.
type LanguageExtractorOptions map[string]interface{}

// Translator for handling all your i18n needs.
type Translator struct {
	// Path - where are the files?
	Path string
	// DefaultLanguage - default is passed as a parameter on New.
	DefaultLanguage string
	// HelperName - name of the view helper. default is "t"
	HelperName string
	// LanguageExtractors - a sorted list of user language extractors.
	LanguageExtractors []LanguageExtractor
	// LanguageExtractorOptions - a map with options to give to LanguageExtractors.
	LanguageExtractorOptions LanguageExtractorOptions
}

// Load translations.
func (t *Translator) Load() error {
	return filepath.Walk(t.Path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			base := filepath.Base(path)
			dir := filepath.Dir(path)

			// Add a prefix to the loaded string, to avoid colilision with ISO lang code
			err = i18n.ParseTranslationFileBytes(fmt.Sprintf("%sbuff%s", dir, base), data)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// AddTranslation directly, without using a file. This is useful if you wish to load translations
// from a database, instead of disk.
func (t *Translator) AddTranslation(lang *language.Language, translations ...translation.Translation) {
	i18n.AddTranslation(lang, translations...)
}

// NewTranslator -
//
// This willalso call t.Load() and load the translations from disk.
func NewTranslator(filePath string, language string) (*Translator, error) {
	t := &Translator{
		Path:            filePath,
		DefaultLanguage: language,
		HelperName:      "t",
		LanguageExtractorOptions: LanguageExtractorOptions{
			"CookieName":       "lang",
			"SessionName":      "lang",
			"RequestParamName": "lang",
			"QueryStringName":  "lang",
		},
		LanguageExtractors: []LanguageExtractor{
			CookieLanguageExtractor,
			SessionLanguageExtractor,
			RequestParamLanguageExtractor,
			QueryStringLanguageExtractor,
			HeaderLanguageExtractor,
		},
	}
	return t, t.Load()
}

// AvailableLanguages gets the list of languages provided by the app.
func (t *Translator) AvailableLanguages() []string {
	lt := i18n.LanguageTags()
	sort.Strings(lt)
	return lt
}

// ExtractLanguage gets language from defined LanguageExtractors
func (t *Translator) ExtractLanguage(c *Context) []string {
	langs := []string{}
	for _, extractor := range t.LanguageExtractors {
		langs = append(langs, extractor(t.LanguageExtractorOptions, c)...)
	}
	// Add default language, even if no language extractor is defined
	langs = append(langs, t.DefaultLanguage)
	return langs
}

// CookieLanguageExtractor is a LanguageExtractor implementation, using a cookie.
func CookieLanguageExtractor(o LanguageExtractorOptions, c *Context) []string {
	langs := make([]string, 0)
	// try to get the language from a cookie:
	if cookieName := o["CookieName"].(string); cookieName != "" {
		if cookie, err := c.Cookie(cookieName); err == nil {
			if cookie != "" {
				langs = append(langs, cookie)
			}
		}
	} else {
		c.Logger().Error("i18n Translator: \"CookieName\" is not defined in LanguageExtractorOptions")
	}
	return langs
}

// SessionLanguageExtractor is a LanguageExtractor implementation, using a session.
func SessionLanguageExtractor(o LanguageExtractorOptions, c *Context) []string {
	langs := make([]string, 0)
	// try to get the language from the session
	if sessionName := o["SessionName"].(string); sessionName != "" {
		if s := c.Session().Get(sessionName); s != nil {
			langs = append(langs, s.(string))
		}
	} else {
		c.Logger().Error("i18n Translator: \"SessionName\" is not defined in LanguageExtractorOptions")
	}
	return langs
}

// HeaderLanguageExtractor is a LanguageExtractor implementation, using a HTTP Accept-Language
// header.
func HeaderLanguageExtractor(o LanguageExtractorOptions, c *Context) []string {
	langs := make([]string, 0)
	// try to get the language from a header:
	acceptLang := c.GetHeader("Accept-Language")
	if acceptLang != "" {
		langs = append(langs, parseAcceptLanguage(acceptLang)...)
	}
	return langs
}

// QueryStringLanguageExtractor is a LanguageExtractor implementation, using a query param from request.
func QueryStringLanguageExtractor(o LanguageExtractorOptions, c *Context) []string {
	langs := make([]string, 0)
	// try to get the language from an URL prefix:
	if urlPrefixName := o["QueryStringName"].(string); urlPrefixName != "" {
		paramLang := c.Query(urlPrefixName)
		if paramLang != "" {
			langs = append(langs, paramLang)
		}
	} else {
		c.Logger().Error("i18n Translator: \"QueryStringName\" is not defined in LanguageExtractorOptions")
	}
	return langs
}

// RequestParamLanguageExtractor is a LanguageExtractor implementation, using a param in the URL.
func RequestParamLanguageExtractor(o LanguageExtractorOptions, c *Context) []string {
	langs := make([]string, 0)
	// try to get the language from an URL prefix:
	if urlPrefixName := o["RequestParamName"].(string); urlPrefixName != "" {
		paramLang := c.Param(urlPrefixName)
		if paramLang != "" {
			langs = append(langs, paramLang)
		}
	} else {
		c.Logger().Error("i18n Translator: \"RequestParamName\" is not defined in LanguageExtractorOptions")
	}
	return langs
}

// Inspired from https://siongui.github.io/2015/02/22/go-parse-accept-language/
// Parse an Accept-Language string to get usable lang values for i18n system
func parseAcceptLanguage(acptLang string) []string {
	var lqs []string

	langQStrs := strings.Split(acptLang, ",")
	for _, langQStr := range langQStrs {
		trimedLangQStr := strings.Trim(langQStr, " ")

		langQ := strings.Split(trimedLangQStr, ";")
		lq := langQ[0]
		lqs = append(lqs, lq)
	}
	return lqs
}
