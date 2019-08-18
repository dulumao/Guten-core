package i18n

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gookit/ini/v2"
	"io/ioutil"
	"os"
	"strings"
)

// language file load mode
const (
	// language name is file name. "en" -> "lang/en.ini"
	FileMode uint8 = 0
	// language name is dir name, will load all file in the dir. "en" -> "lang/en/*.ini"
	DirMode uint8 = 1
)

// I18n language manager
type i18n struct {
	// languages data
	data map[string]*ini.Ini

	// language files directory
	i18nDir string
	// language list {en:English, zh-CN:简体中文}
	languages map[string]string
	// loaded lang files
	// loadedFiles []string

	// mode for the load language files. mode: 0 single file, 1 multi dir
	LoadMode uint8
	// default language name. eg. "en"
	DefaultLang string
	// spare(fallback) language name. eg. "en"
	FallbackLang string
}

var I18n *i18n

func New(loadMode uint8, i18nDir string, defaultLang string, languages map[string]string) {
	I18n = &i18n{
		data:      make(map[string]*ini.Ini, 0),
		LoadMode:  loadMode,
		i18nDir:   i18nDir,
		languages: languages,

		DefaultLang: defaultLang,
	}

	if I18n.LoadMode == FileMode {
		I18n.loadSingleFiles()
	} else if I18n.LoadMode == DirMode {
		I18n.loadDirFiles()
	} else {
		panic("invalid load mode setting. only allow 0, 1")
	}
}

// load language files when LoadMode is 0
func (l *i18n) loadSingleFiles() {
	pathSep := string(os.PathSeparator)

	for lang := range l.languages {
		lData := ini.New()
		err := lData.LoadFiles(l.i18nDir + pathSep + lang + ".ini")
		if err != nil {
			panic("fail to load language: " + lang + ", error " + err.Error())
		}

		l.data[lang] = lData
	}
}

// load language files when LoadMode is 1
func (l *i18n) loadDirFiles() {
	pathSep := string(os.PathSeparator)

	for lang := range l.languages {
		dirPath := l.i18nDir + pathSep + lang
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			panic("read dir fail: " + dirPath + ", error " + err.Error())
		}

		sl := l.data[lang]

		for _, fi := range files {
			// skip dir and filter the specified format
			if fi.IsDir() || !strings.HasSuffix(fi.Name(), ".ini") {
				continue
			}

			var err error
			if sl != nil {
				err = sl.LoadFiles(dirPath + pathSep + fi.Name())
			} else { // create new language instance
				sl = ini.New()
				err = sl.LoadFiles(dirPath + pathSep + fi.Name())
				l.data[lang] = sl
			}

			if err != nil {
				panic("fail to load language file: " + lang + ", error " + err.Error())
			}
		}
	}
}

func (l *i18n) DefTr(key string, args ...interface{}) string {
	return l.Tr(l.DefaultLang, key, args...)
}

// Tr translate from a lang by key
// site.name => [site]
//  			name = my blog
func (l *i18n) Tr(lang, key string, args ...interface{}) string {
	if !l.HasLang(lang) {
		// find from fallback lang
		val := l.transFromFallback(key)
		if val == "" {
			return key
		}

		if len(args) > 0 { // if has args
			val = fmt.Sprintf(val, args...)
		}

		return val
	}

	val, ok := l.data[lang].GetValue(key)
	if !ok {
		// find from fallback lang
		val = l.transFromFallback(key)
		if val == "" {
			return key
		}
	}

	if len(args) > 0 { // if has args
		val = fmt.Sprintf(val, args...)
	}

	return val
}

// translate from fallback language
func (l *i18n) transFromFallback(key string) (val string) {
	fb := l.FallbackLang
	if !l.HasLang(fb) {
		return
	}

	return l.data[fb].String(key)
}

// HasLang in the manager
func (l *i18n) HasLang(lang string) bool {
	_, ok := l.languages[lang]
	return ok
}

// NewLang create/add a new language
// Usage:
// 	i18n.NewLang("zh-CN", "简体中文")
func (l *i18n) NewLang(lang string, name string) {
	// lang exist
	if _, ok := l.languages[lang]; ok {
		return
	}

	l.data[lang] = ini.New()
	l.languages[lang] = name
}

// LoadFile append data to a exist language
// Usage:
// 	i18n.LoadFile("zh-CN", "path/to/zh-CN.ini")
func (l *i18n) LoadFile(lang string, file string) (err error) {
	// append data
	if ld, ok := l.data[lang]; ok {
		err = ld.LoadFiles(file)
		if err != nil {
			return
		}
	} else {
		err = errors.New("language" + lang + " not exist, please create it before load data")
	}

	return
}

// LoadString load language data form a string
// Usage:
// i18n.Set("zh-CN", "name = blog")
func (l *i18n) LoadString(lang string, data string) (err error) {
	// append data
	if ld, ok := l.data[lang]; ok {
		err = ld.LoadStrings(data)
		if err != nil {
			return
		}
	} else {
		err = errors.New("language" + lang + " not exist, please create it before load data")
	}

	return
}

// Export a language data as INI string
func (l *i18n) Export(lang string) string {
	if _, ok := l.languages[lang]; !ok {
		return ""
	}

	var buf bytes.Buffer

	_, err := l.data[lang].WriteTo(&buf)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

// Languages get all languages
func (l *i18n) Languages() map[string]string {
	return l.languages
}
