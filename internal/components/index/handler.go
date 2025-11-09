package index

import (
	"net/http"

	"github.com/Dekamik/farstu/internal/components/shared"
)

type Index struct {}

func Handler(w http.ResponseWriter, _ *http.Request) {
	data := Index{}
	shared.ExecuteLayout(w, "internal/components/index/index.html", "/", data)
}
