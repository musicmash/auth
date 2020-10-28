package spotify

import (
	"errors"
	"net/url"
)

const stateAuth = "auth"

var (
	errCodeIsEmpty  = errors.New("query arg 'code' is empty")
	errUnknownState = errors.New("unknown state")
)

func validateStateAndCode(values url.Values) error {
	code := values.Get("code")
	if code == "" {
		return errCodeIsEmpty
	}

	actualState := values.Get("state")
	if actualState != stateAuth {
		return errUnknownState
	}

	return nil
}
