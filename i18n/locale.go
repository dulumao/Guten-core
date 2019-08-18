package i18n

import (
	"github.com/dulumao/Guten-utils/conv"
	"github.com/gookit/rux"
	"golang.org/x/text/language"
)

type Locale struct {
	c *rux.Context
}

func Get(c *rux.Context) *Locale {
	return &Locale{
		c: c,
	}
}

func (self *Locale) Tr(i ...interface{}) string {
	var key = conv.String(i[0])
	var args []interface{}

	if len(i) > 1 {
		args = i[1:]
	}

	t, _, err := language.ParseAcceptLanguage(self.c.Req.Header.Get("Accept-Language"))

	if err == nil {
		if len(t) > 0 {
			return I18n.Tr(t[0].String(), key, args...)
		}
	}

	return I18n.DefTr(key, args...)
}
