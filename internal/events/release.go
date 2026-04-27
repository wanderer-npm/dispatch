package events

import (
	"encoding/json"

	"github.com/wanderer-npm/dispatch/internal/discord"
)

type releasePayload struct {
	Action  string `json:"action"`
	Release struct {
		TagName    string `json:"tag_name"`
		Name       string `json:"name"`
		HTMLURL    string `json:"html_url"`
		Prerelease bool   `json:"prerelease"`
		Body       string `json:"body"`
	} `json:"release"`
	Repository repo `json:"repository"`
	Sender     user `json:"sender"`
}

func Release(body []byte) (*discord.Embed, error) {
	var p releasePayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	if p.Action != "published" {
		return nil, nil
	}

	rel := p.Release
	name := rel.Name
	if name == "" {
		name = rel.TagName
	}

	relLink := name
	if rel.HTMLURL != "" {
		relLink = "[" + name + "](" + rel.HTMLURL + ")"
	}

	desc := relLink + " in " + repoLink(p.Repository)
	if rel.Prerelease {
		desc += "\n*Pre-release*"
	}
	if rel.Body != "" {
		desc += "\n\n" + truncate(rel.Body, 300)
	}

	return &discord.Embed{
		Title:       "Release published — " + rel.TagName,
		Description: desc,
		Color:       colorRelease,
		Author:      authorOf(p.Sender),
	}, nil
}
