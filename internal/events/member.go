package events

import (
	"encoding/json"

	"github.com/wanderer-npm/dispatch/internal/discord"
)

type memberPayload struct {
	Action     string `json:"action"`
	Member     user   `json:"member"`
	Repository repo   `json:"repository"`
	Sender     user   `json:"sender"`
}

func Member(body []byte) (*discord.Embed, error) {
	var p memberPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	if p.Member.Login == "" {
		return nil, nil
	}

	var label string
	var color int
	var verb string

	switch p.Action {
	case "added":
		label, color, verb = "Collaborator added", colorMemberAdd, "added to"
	case "removed":
		label, color, verb = "Collaborator removed", colorMemberRemove, "removed from"
	default:
		return nil, nil
	}

	memberLink := p.Member.Login
	if p.Member.HTMLURL != "" {
		memberLink = "[" + p.Member.Login + "](" + p.Member.HTMLURL + ")"
	}

	return &discord.Embed{
		Title:       label,
		Description: memberLink + " was " + verb + " " + repoLink(p.Repository),
		Color:       color,
		Author:      authorOf(p.Sender),
	}, nil
}
