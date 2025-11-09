package shared

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Dekamik/farstu/internal/asserts"
)

var layouts = template.Must(template.ParseGlob("internal/components/shared/layout/*.html"))

type Nav struct {
	Highlight string
}

type Site struct {
	Lang  string
	Theme string
}

type Layout[T any] struct {
	Nav  Nav
	Site Site
	Data T
}

func ExecuteLayout[T any](w http.ResponseWriter, templatePath string, highlightNav string, data T) {
	asserts.Assert(templatePath != "", "template path cannot be empty")

	_, err := os.Stat(templatePath)
	asserts.Assert(!errors.Is(err, os.ErrNotExist), fmt.Sprintf("template %s must exist", templatePath))
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.Must(layouts.Clone()).ParseFiles(templatePath))

	layoutData := Layout[T]{
		Nav: Nav{
			Highlight: highlightNav,
		},
		Site: Site{
			Lang:  "sv",
			Theme: "synthwave",
		},
		Data: data,
	}

	tmpl.ExecuteTemplate(w, "layout.html", layoutData)
}
