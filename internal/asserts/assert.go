package asserts

import (
	"errors"
	"log/slog"
	"os"
)

func Assert(predicate bool, reason string) {
    if !predicate {
        panic(errors.New(reason))
    }
}

func ErrAssert(predicate bool, message string) {
    if !predicate {
        slog.Error(message)
        os.Exit(1)
    }
}
