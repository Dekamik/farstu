package deviations

import (
	"net/http"

	"github.com/Dekamik/farstu/internal/components/shared"
)

type Deviations struct {}

func Handler(w http.ResponseWriter, _ *http.Request) {
	data := Deviations{}
	shared.ExecuteLayout(w, "internal/components/deviations/deviations.html", "deviations", data)
}
