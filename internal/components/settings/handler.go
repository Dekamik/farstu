package settings

import (
	"net/http"

	"github.com/Dekamik/farstu/internal/components/shared"
)

type Settings struct {}

func Handler(w http.ResponseWriter, _ *http.Request) {
	data := Settings{}
	shared.ExecuteLayout(w, "internal/components/settings/settings.html", "settings", data)
}
