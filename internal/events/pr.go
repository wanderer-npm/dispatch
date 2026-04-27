package events

import (
	"encoding/json"

	"github.com/wanderer-npm/dispatch/internal/discord"
)

type prPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		Merged  bool   `json:"merged"`
		Body    string `json:"body"`
		User    user   `json:"user"`
	} `json:"pull_request"`
	Repository repo `json:"repository"`
	Sender     user `json:"sender"`
}

func PullRequest(body []byte) (*discord.Embed, error) {
	var p prPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	var label string
	var color int

	switch p.Action {
	case "opened":
		label, color = "Pull request opened", colorPROpen
	case "reopened":
		label, color = "Pull request reopened", colorPROpen
	case "closed":
		if p.PullRequest.Merged {
			label, color = "Pull request merged", colorPRMerged
		} else {
			label, color = "Pull request closed", colorPRClose
		}
	default:
		return nil, nil
	}

	prLink := p.PullRequest.Title
	if p.PullRequest.HTMLURL != "" {
		prLink = "[" + p.PullRequest.Title + "](" + p.PullRequest.HTMLURL + ")"
	}

	desc := prLink + " in " + repoLink(p.Repository)
	if p.PullRequest.Body != "" {
		desc += "\n\n" + truncate(p.PullRequest.Body, 200)
	}

	return &discord.Embed{
		Title:       label,
		Description: desc,
		Color:       color,
		Author:      authorOf(p.Sender),
	}, nil
}
