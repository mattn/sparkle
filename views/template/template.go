package template

import (
	ht "html/template"
	"io"
	"github.com/sekhat/sparkle/views"
)

type templateView struct {
	template *ht.Template
}

func (v *templateView) Execute(w io.Writer, model interface{}) error {
	return v.template.Execute(w, model)
}

// Creates a new ViewWriter using html/template as it's implementation,
// specifing the files to use for the template
func New(templateFiles ...string) views.ViewWriter {
	return &templateView{
		ht.Must(ht.ParseFiles(templateFiles...)),
	}
}
