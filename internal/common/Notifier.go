package common

type Notifier interface {
    Notify(storedIp string, attemptIp string, userAgent string, userID string) error
}
