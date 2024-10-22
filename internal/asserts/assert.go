package asserts

import (
	"errors"
)

func Assert(predicate bool, reason string) {
    if !predicate {
        panic(errors.New(reason))
    }
}
