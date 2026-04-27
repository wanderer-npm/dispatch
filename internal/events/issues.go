package events

import (
	"encoding/json"
	"fmt"

	"github.com/wanderer-npm/dispatch/internal/discord"
)

type issuesPayload struct {
	Action string `json:"action"`
	Issue  struct {
		Number  int    `json:"number"`
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		Body    string `json:"body"`
		User    user   `json:"user"`
	} `json:"issue"`
	Repository repo `json:"repository"`
	Sender     user `json:"sender"`
}

func Issues(body []byte) (*discord.Embed, error) {
	var p issuesPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	var label string
	var color int

	switch p.Action {
	case "opened":
		label, color = "Issue opened", colorIssueOpen
	case "closed":
		label, color = "Issue closed", colorIssueClosed
	case "reopened":
		label, color = "Issue reopened", colorIssueOpen
	default:
		return nil, nil
	}

	issue := p.Issue
	issueLink := fmt.Sprintf("#%d %s", issue.Number, issue.Title)
	if issue.HTMLURL != "" {
		issueLink = fmt.Sprintf("[#%d %s](%s)", issue.Number, issue.Title, issue.HTMLURL)
	}

	desc := issueLink + " in " + repoLink(p.Repository)
	if issue.Body != "" {
		desc += "\n\n" + truncate(issue.Body, 200)
	}

	return &discord.Embed{
		Title:       label,
		Description: desc,
		Color:       color,
		Author:      authorOf(p.Sender),
	}, nil
}
