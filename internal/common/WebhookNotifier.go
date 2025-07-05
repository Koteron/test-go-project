package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WebhookNotifier struct {
    URL string
    Client *http.Client
}

func (w *WebhookNotifier) Notify(storedIp string, attemptIp string, userAgent string, userID string) error {
    payload := gin.H{
		"event": "refresh_attempt_from_unrecognized_ip",
		"user_id": userID,
		"attempt_ip": attemptIp,
		"stored_ip": storedIp,
		"user_agent": userAgent,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	jsonData, _ := json.Marshal(payload)

    _, err := w.Client.Post(w.URL, "application/json", bytes.NewBuffer(jsonData))
    return err
}