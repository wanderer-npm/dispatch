package events

import (
	"encoding/json"
	"strings"

	"github.com/wanderer-npm/dispatch/internal/discord"
)

type repositoryPayload struct {
	Action     string          `json:"action"`
	Repository repo            `json:"repository"`
	Sender     user            `json:"sender"`
	Changes    json.RawMessage `json:"changes"`
}

var repoActionLabel = map[string]string{
	"created":     "Repository created",
	"deleted":     "Repository deleted",
	"renamed":     "Repository renamed",
	"archived":    "Repository archived",
	"unarchived":  "Repository unarchived",
	"transferred": "Repository transferred",
	"publicized":  "Repository made public",
	"privatized":  "Repository made private",
}

var repoActionColor = map[string]int{
	"created":    colorRepoCreate,
	"deleted":    colorRepoDelete,
	"publicized": colorRepoCreate,
	"privatized": colorRepoDelete,
}

func Repository(body []byte) (*discord.Embed, error) {
	var p repositoryPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	label, ok := repoActionLabel[p.Action]
	if !ok {
		return nil, nil
	}

	color, ok := repoActionColor[p.Action]
	if !ok {
		color = colorRepoUpdate
	}

	var lines []string
	lines = append(lines, repoLink(p.Repository))

	if p.Action == "renamed" && p.Changes != nil {
		var changes struct {
			Repository struct {
				Name struct {
					From string `json:"from"`
				} `json:"name"`
			} `json:"repository"`
		}
		if err := json.Unmarshal(p.Changes, &changes); err == nil {
			if from := changes.Repository.Name.From; from != "" {
				lines = append(lines, "Renamed from **"+from+"**")
			}
		}
	}

	if p.Repository.Description != "" {
		lines = append(lines, "*"+p.Repository.Description+"*")
	}

	return &discord.Embed{
		Title:       label,
		Description: strings.Join(lines, "\n"),
		Color:       color,
		Author:      authorOf(p.Sender),
	}, nil
}
