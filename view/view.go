package view

import (
	"html/template"
)

func New(dirs []string) *Renderer {
	return &Renderer{
		Cached: false,
		// Engine: jet.NewHTMLSet(dirs...),
		Engine: NewSetLoader(template.HTMLEscape, dirs...),
	}
}
