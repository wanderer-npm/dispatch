package events

import (
	"encoding/json"

	"github.com/wanderer-npm/dispatch/internal/discord"
)

type forkPayload struct {
	Forkee     repo `json:"forkee"`
	Repository repo `json:"repository"`
	Sender     user `json:"sender"`
}

func Fork(body []byte) (*discord.Embed, error) {
	var p forkPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	desc := repoLink(p.Repository) + " was forked"
	if p.Forkee.FullName != "" {
		desc += "\nFork: " + repoLink(p.Forkee)
	}

	return &discord.Embed{
		Title:       "Repository forked",
		Description: desc,
		Color:       colorFork,
		Author:      authorOf(p.Sender),
	}, nil
}
