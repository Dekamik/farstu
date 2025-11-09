package helloworld

import (
	"net/http"

	"github.com/Dekamik/farstu/internal/components/shared"
)

type HelloWorld struct {
	Text string
}

func Handler(w http.ResponseWriter, _ *http.Request) {
	data := HelloWorld{
		Text: "World",
	}
	shared.ExecuteLayout(w, "internal/components/helloworld/hello_world.html", data)
}
