package settings

import (
	"net/http"

	"github.com/Dekamik/farstu/internal/routes/shared"
)

type Settings struct {}

func Handler(w http.ResponseWriter, _ *http.Request) {
	data := Settings{}
	shared.ExecuteLayout(w, "settings", data, "internal/components/settings/settings.html")
}
