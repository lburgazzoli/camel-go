package wasm

type MatchError struct {
	Status bool
}

func (e MatchError) Error() string {
	if !e.Status {
		return "does not match"
	}

	return "matches"
}

var ErrPredicateDoesNotMatch = MatchError{Status: false}
var ErrPredicateMatches = MatchError{Status: true}
