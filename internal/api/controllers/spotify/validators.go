package spotify

import (
	"errors"
	"net/url"
)

const state = "auth"

func validateStateAndCode(values url.Values) error {
	code := values.Get("code")
	if code == "" {
		return errors.New("didn't get access code")
	}

	actualState := values.Get("state")
	if actualState != state {
		return errors.New("redirect state parameter doesn't match")
	}

	return nil
}
