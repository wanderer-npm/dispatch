package events

import (
	"encoding/json"

	"github.com/wanderer-npm/dispatch/internal/discord"
)

type starPayload struct {
	Action     string `json:"action"`
	Repository repo   `json:"repository"`
	Sender     user   `json:"sender"`
}

func Star(body []byte) (*discord.Embed, error) {
	var p starPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	if p.Action != "created" {
		return nil, nil
	}

	return &discord.Embed{
		Title:       "Repository starred",
		Description: repoLink(p.Repository),
		Color:       colorStar,
		Author:      authorOf(p.Sender),
	}, nil
}
