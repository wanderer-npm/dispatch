package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Send(webhookURL string, embed Embed) error {
	body, err := json.Marshal(Payload{Embeds: []Embed{embed}})
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord returned %d", resp.StatusCode)
	}
	return nil
}
