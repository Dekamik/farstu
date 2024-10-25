package asserts

import (
	"errors"
	"os"
)

var isDevelopment = os.Getenv("FARSTU_ENVIRONMENT") == "Development"

func Assert(predicate bool, reason string) {
	if isDevelopment && !predicate {
		panic(errors.New(reason))
	}
}
