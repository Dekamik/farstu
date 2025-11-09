package deviations

import (
	"net/http"

	"github.com/Dekamik/farstu/internal/routes/shared"
)

type Deviations struct {}

func Handler(w http.ResponseWriter, _ *http.Request) {
	data := Deviations{}
	shared.ExecuteLayout(w, "deviations", data, "internal/components/deviations/deviations.html")
}
