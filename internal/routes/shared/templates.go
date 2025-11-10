package shared

import (
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/Dekamik/farstu/internal/asserts"
)

var layouts = template.Must(template.ParseGlob("internal/routes/shared/layout/*.html"))

type Nav struct {
	CurrentTime TimeData
	Highlight   string
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

func ExecuteLayout[T any](w http.ResponseWriter, highlightNav string, data T, templatePath ...string) {
	for _, t := range templatePath {
		_, err := os.Stat(t)
		asserts.Assert(!errors.Is(err, os.ErrNotExist), fmt.Sprintf("template %s must exist", templatePath))
		if err != nil {
			panic(err)
		}
	}

	tmpl := template.Must(template.Must(layouts.Clone()).ParseFiles(templatePath...))

	layoutData := Layout[T]{
		Nav: Nav{
			CurrentTime: GetTime(),
			Highlight: highlightNav,
		},
		Site: Site{
			Lang:  "sv",
			Theme: "synthwave",
		},
		Data: data,
	}

	err := tmpl.ExecuteTemplate(w, "layout.html", layoutData)
	if err != nil {
		slog.Error(err.Error())
	}
}
