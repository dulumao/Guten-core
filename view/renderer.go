package view

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/dulumao/Guten-utils/conv"
	"github.com/dulumao/Guten-utils/dump"
	"github.com/gookit/rux"
	"github.com/gookit/validate"
	"github.com/dulumao/Guten-core/i18n"
	"io"
	"strings"
)

// 模板注册
type Renderer struct {
	Cached bool
	Engine *jet.Set
}

func NewSetLoader(escapee jet.SafeWriter, dirs ...string) *jet.Set {
	return jet.NewSetLoader(escapee, &OSFileSystemLoader{dirs: dirs})
}

func (self *Renderer) Render(w io.Writer, name string, data interface{}, ctx *rux.Context) error {
	self.Engine.SetDevelopmentMode(!self.Cached)

	self.Engine.AddGlobal("dump", func(i ...interface{}) {
		dump.DD(i...)
	})
	self.Engine.AddGlobal("dump2", func(i ...interface{}) {
		dump.DD2(i...)
	})
	self.Engine.AddGlobal("printf", func(format string, a ...interface{}) string {
		return fmt.Sprintf(format, a ...)
	})
	self.Engine.AddGlobal("isEqual", func(v1 interface{}, v2 interface{}) bool {
		if v1 == v2 {
			return true
		}

		return false
	})
	self.Engine.AddGlobal("tr", func(i ...interface{}) string {
		locale := i18n.Get(ctx)

		return locale.Tr(i...)
	})
	self.Engine.AddGlobal("isCurrentRoute", func(name string) bool {
		var currentRouteName = ctx.MustGet(rux.CTXCurrentRouteName)

		if conv.String(currentRouteName) == name {
			return true
		}

		return false
	})
	self.Engine.AddGlobal("currentRouteName", func() string {
		var currentRouteName = ctx.MustGet(rux.CTXCurrentRouteName)

		return conv.String(currentRouteName)
	})
	self.Engine.AddGlobal("route", func(name string, args ...interface{}) string {
		if name == "" {
			return "#"
		}

		var router = ctx.Router()

		if len(args) > 0 {
			// var buildRequestURLs []interface{}
			//
			// for _, arg := range args {
			// 	buildRequestURLs = append(buildRequestURLs, conv.String(arg))
			// }
			//
			// return router.BuildRequestURL(name, buildRequestURLs...).String()
			return router.BuildRequestURL(name, args...).String()
		}

		return router.BuildRequestURL(name).String()
	})
	self.Engine.AddGlobal("valid", func(key string) bool {
		var err = ctx.MustGet("error")

		if validErrors, can := err.(validate.Errors); can {
			validResult := validErrors.All()

			if _, ok := validResult[key]; ok {
				return true
			}
		}

		return false

		return false
	})
	self.Engine.AddGlobal("validText", func(key string) string {
		var err = ctx.MustGet("error")

		if validErrors, can := err.(validate.Errors); can {
			return validErrors.Get(key)
		}

		return ""
	})
	self.Engine.AddGlobal("validField", func(key string) []string {
		var err = ctx.MustGet("error")

		if validErrors, can := err.(validate.Errors); can {
			return validErrors.Field(key)
		}

		return []string{}
	})

	for k, v := range Funcs.Items() {
		self.Engine.AddGlobal(conv.String(k), v)
	}

	t, err := self.Engine.GetTemplate(name)

	if err != nil {
		panic(err)
	}

	vars := make(jet.VarMap)
	vars.Set("ctx", ctx)

	for k, v := range Vars.Items() {
		vars.Set(conv.String(k), v)
	}

	if err = t.Execute(w, vars, data); err != nil {
		panic(err)
	}

	return nil
}

func GetViewEventName(name string) string {
	var names = strings.Split(name, "/")

	return strings.Join(names, ".")
}
