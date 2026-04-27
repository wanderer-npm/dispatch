package events

import (
	"encoding/json"
	"fmt"

	"github.com/wanderer-npm/dispatch/internal/discord"
)

type refPayload struct {
	Ref        string `json:"ref"`
	RefType    string `json:"ref_type"`
	Repository repo   `json:"repository"`
	Sender     user   `json:"sender"`
}

func Create(body []byte) (*discord.Embed, error) {
	var p refPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	if p.Ref == "" || p.RefType == "" {
		return nil, nil
	}

	refURL := ""
	if p.Repository.HTMLURL != "" {
		if p.RefType == "tag" {
			refURL = fmt.Sprintf("%s/releases/tag/%s", p.Repository.HTMLURL, p.Ref)
		} else {
			refURL = fmt.Sprintf("%s/tree/%s", p.Repository.HTMLURL, p.Ref)
		}
	}

	refLink := fmt.Sprintf("[`%s`](%s)", p.Ref, refURL)
	if refURL == "" {
		refLink = "`" + p.Ref + "`"
	}

	return &discord.Embed{
		Title:       capitalize(p.RefType) + " created",
		Description: refLink + " in " + repoLink(p.Repository),
		Color:       colorCreate,
		Author:      authorOf(p.Sender),
	}, nil
}

func Delete(body []byte) (*discord.Embed, error) {
	var p refPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	if p.Ref == "" || p.RefType == "" {
		return nil, nil
	}

	return &discord.Embed{
		Title:       capitalize(p.RefType) + " deleted",
		Description: "`" + p.Ref + "` in " + repoLink(p.Repository),
		Color:       colorDelete,
		Author:      authorOf(p.Sender),
	}, nil
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return string(s[0]-32) + s[1:]
}
