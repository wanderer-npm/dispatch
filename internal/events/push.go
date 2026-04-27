package events

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wanderer-npm/dispatch/internal/discord"
)

type pushPayload struct {
	Ref        string   `json:"ref"`
	Before     string   `json:"before"`
	After      string   `json:"after"`
	Repository repo     `json:"repository"`
	Sender     user     `json:"sender"`
	Commits    []commit `json:"commits"`
}

func Push(body []byte) (*discord.Embed, error) {
	var p pushPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	if len(p.Commits) == 0 {
		return nil, nil
	}

	branch := refName(p.Ref)
	if strings.HasPrefix(branch, "gh-readonly-queue/") || strings.HasPrefix(branch, "pr-") {
		return nil, nil
	}

	count := len(p.Commits)
	noun := "commit"
	if count != 1 {
		noun = "commits"
	}

	compareURL := ""
	if p.Repository.HTMLURL != "" && len(p.Before) >= 12 && len(p.After) >= 12 {
		compareURL = fmt.Sprintf("%s/compare/%s...%s", p.Repository.HTMLURL, p.Before[:12], p.After[:12])
	}

	var lines []string
	for _, c := range p.Commits {
		if c.ID == "" || c.Message == "" {
			continue
		}
		shortSHA := c.ID[:7]
		firstLine := strings.SplitN(c.Message, "\n", 2)[0]
		msg := truncate(firstLine, 52)

		name := c.Author.Username
		if name == "" {
			name = c.Author.Name
		}
		if name == "" {
			name = "unknown"
		}

		commitURL := fmt.Sprintf("%s/commit/%s", p.Repository.HTMLURL, c.ID)
		lines = append(lines, fmt.Sprintf("[`%s`](%s) %s — %s", shortSHA, commitURL, msg, name))
	}

	desc := strings.Join(lines, "\n")
	if desc == "" {
		desc = "No commit details."
	}

	return &discord.Embed{
		Title:       fmt.Sprintf("[%s:%s] %d new %s", p.Repository.FullName, branch, count, noun),
		URL:         compareURL,
		Description: truncate(desc, 4000),
		Color:       colorPush,
		Author:      authorOf(p.Sender),
	}, nil
}

func refName(ref string) string {
	parts := strings.SplitN(ref, "/", 3)
	if len(parts) == 3 {
		return parts[2]
	}
	return ref
}
