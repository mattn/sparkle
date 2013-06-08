package template

import (
	ht "html/template"
	"io"
	"sparkle/views"
)

type templateView struct {
	template *ht.Template
}

func (v *templateView) Execute(w io.Writer, model interface{}) error {
	return v.template.Execute(w, model)
}

func New(templateFiles ...string) views.ViewWriter {	
	return &templateView{
		ht.Must(ht.ParseFiles(templateFiles...)),
	}
}