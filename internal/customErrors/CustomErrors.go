package customErrors

import "errors"

var (
    ErrFailedSave         = errors.New("could not save refresh record")
    ErrDB                 = errors.New("SQL querry could not go through")
    ErrWebhook            = errors.New("could not send webhook notification")
	ErrRecordNotFound     = errors.New("could not find such refresh record")
    ErrUserAgentMismatch  = errors.New("UserAgent mismatch")
    ErrUserIDMismatch     = errors.New("GUID from access token does not match DB id")
    ErrRefreshMismatch    = errors.New("refresh token mismatch")
    ErrMissingUserID      = errors.New("token missing user id")
    ErrTokenPairing       = errors.New("tokens are not paired")
    ErrUUIDParsing        = errors.New("could not parse UUID")
	ErrUnexpectedState    = errors.New("unexpected internal state")
)

func IsCustomError(err error) bool {
	return errors.Is(err, ErrFailedSave) ||
		errors.Is(err, ErrDB) ||
		errors.Is(err, ErrWebhook) ||
        errors.Is(err, ErrRecordNotFound) ||
        errors.Is(err, ErrUserAgentMismatch) ||
        errors.Is(err, ErrUserIDMismatch) ||
        errors.Is(err, ErrRefreshMismatch) ||
        errors.Is(err, ErrTokenPairing) ||
        errors.Is(err, ErrUUIDParsing) ||
        errors.Is(err, ErrUnexpectedState)
}